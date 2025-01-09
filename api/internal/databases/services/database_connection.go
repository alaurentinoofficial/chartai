package databases_services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	databases_models "github.com/alaurentinoofficial/chartai/internal/databases/models"
)

func NewDatabaseConnectionFactory(dtype databases_models.DatabaseType, connectionString string) (DatabaseConnection, error) {
	if dtype == databases_models.PostgresDatabaseType {
		conn, err := sql.Open("postgres", connectionString)
		if err != nil {
			return nil, err
		}

		return &postgreDatabaseConnection{
			conn: conn,
		}, nil
	}

	return nil, fmt.Errorf("Invalid DatabaseType")
}

type DatabaseConnection interface {
	GetSchema(ctx context.Context) (*[]databases_models.Table, error)
	RunQuery(ctx context.Context, query string) (*[]map[string]any, error)
	Close(ctx context.Context)
}

type postgreDatabaseConnection struct {
	conn *sql.DB
}

func (c *postgreDatabaseConnection) GetSchema(ctx context.Context) (*[]databases_models.Table, error) {
	query := `
WITH table_constraints AS (
    SELECT
        tc.table_name,
        kcu.column_name,
        tc.constraint_type,
        cc.table_name AS foreign_target_table,
        cc.column_name AS foreign_target_column
    FROM
        information_schema.table_constraints tc
    LEFT JOIN information_schema.key_column_usage kcu
        ON tc.constraint_name = kcu.constraint_name
        AND tc.table_schema = kcu.table_schema
    LEFT JOIN information_schema.constraint_column_usage cc
        ON tc.constraint_name = cc.constraint_name
        AND tc.table_schema = cc.table_schema
    WHERE
        tc.table_schema = 'public' -- Adjust schema if needed
),
table_columns AS (
    SELECT
        t.table_name,
        c.column_name,
        c.data_type,
        CASE
            WHEN tc.constraint_type = 'PRIMARY KEY' THEN TRUE
            ELSE FALSE
        END AS is_primary_key,
        CASE
            WHEN tc.constraint_type = 'FOREIGN KEY' THEN
                json_build_object(
                    'target_table', tc.foreign_target_table,
                    'target_column', tc.foreign_target_column
                )
            ELSE NULL
        END AS foreign_key_details
    FROM
        information_schema.tables t
    JOIN information_schema.columns c
        ON t.table_name = c.table_name
        AND t.table_schema = c.table_schema
    LEFT JOIN table_constraints tc
        ON t.table_name = tc.table_name
        AND c.column_name = tc.column_name
        AND t.table_schema = 'public' -- Adjust schema if needed
    WHERE
        t.table_schema = 'public' -- Adjust schema if needed
        AND t.table_type = 'BASE TABLE'
)
SELECT
    table_name,
    json_agg(
        json_build_object(
            'name', column_name,
            'type', data_type,
            'is_primary_key', is_primary_key,
            'foreign_key', foreign_key_details
        )
    ) AS columns
FROM
    table_columns
GROUP BY
    table_name;
	`

	rows, err := c.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tables []databases_models.Table

	for rows.Next() {
		var tableName string
		var columnsJSON string

		if err := rows.Scan(&tableName, &columnsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var columns []databases_models.Column
		if err := json.Unmarshal([]byte(columnsJSON), &columns); err != nil {
			return nil, fmt.Errorf("failed to unmarshal columns JSON: %w", err)
		}

		tables = append(tables, databases_models.Table{
			Name:    tableName,
			Columns: columns,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &tables, nil
}

func (c *postgreDatabaseConnection) RunQuery(ctx context.Context, query string) (*[]map[string]any, error) {
	// Execute the query
	rows, err := c.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names and types
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// Create result slice
	var result [](map[string]any)

	// Iterate through rows
	for rows.Next() {
		// Create a slice of interface{} to hold the row values
		scanArgs := make([]any, len(columns))
		for i := range scanArgs {
			scanArgs[i] = new(interface{})
		}

		// Scan the row into the interfaces
		if err := rows.Scan(scanArgs...); err != nil {
			continue
		}

		// Create a new entry for this row
		entry := make(map[string]any)

		// Convert each column value to the appropriate type
		for i, col := range columns {
			value := *(scanArgs[i].(*interface{}))

			switch v := value.(type) {
			case []byte:
				// Handle numeric types that might come as []byte
				switch col.DatabaseTypeName() {
				case "NUMERIC", "DECIMAL", "FLOAT", "FLOAT4", "FLOAT8", "DOUBLE":
					if f, err := strconv.ParseFloat(string(v), 64); err == nil {
						entry[col.Name()] = f
					} else {
						entry[col.Name()] = string(v)
					}
				default:
					entry[col.Name()] = string(v)
				}
			case time.Time:
				entry[col.Name()] = v.Format(time.RFC3339)
			case nil:
				entry[col.Name()] = nil
			default:
				entry[col.Name()] = v
			}
		}

		result = append(result, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *postgreDatabaseConnection) Close(ctx context.Context) {
	c.conn.Close()
}
