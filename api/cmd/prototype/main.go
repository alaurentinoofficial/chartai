package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// Global Variables
var (
	DatabaseConnection string
	LLM                *openai.LLM
)

func init() {
	DatabaseConnection = os.Getenv("DATABASE_CONNECTION")

	llm, err := openai.New(openai.WithModel("gpt-4o-mini"))
	if err != nil {
		log.Fatal(err.Error())
	}
	LLM = llm
}

func main() {
	app := fiber.New()
	ctx := context.Background()

	app.Post("/chart", func(c *fiber.Ctx) error {
		var request RequestChartPrompt
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		_, sql, chart, data := HandleUserPrompt(ctx, request)

		result := ChartResponse{
			Chart: chart.Name,
			Query: sql,
			Data:  data,
		}

		return c.JSON(result)
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

type RequestChartPrompt struct {
	Prompt string `json:"prompt"`
}

func HandleUserPrompt(ctx context.Context, request RequestChartPrompt) ([]Table, string, *Chart, ChartData) {
	prompt := request.Prompt

	// Get Database connection
	conn, err := getDatabaseConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get list of tables from Database
	tables, err := getTablesFromDatabase(ctx, conn)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Query LLM to get best Chart for this prompt
	chart, err := QueryLLMChart(ctx, prompt)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Query LLM to generate a SQL Statement compatible with the Database
	sql := QueryLLMSqlStatement(ctx, prompt, tables, chart)

	// Run the SQL Statement on the target Database
	output := RunSqlQuery(ctx, conn, sql, chart)

	// Return the data for this Chart
	return tables, sql, chart, output
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	IsPrimaryKey bool        `json:"is_primary_key"`
	ForeignKey   *ForeignKey `json:"foreign_key,omitempty"`
}

type ForeignKey struct {
	TargetTable  string `json:"target_table"`
	TargetColumn string `json:"target_column"`
}

type OutputTable struct {
	Columns []string
	Values  []any
}

type ChartResponse struct {
	Chart string
	Query string
	Data  ChartData
}

func getDatabaseConnection(ctx context.Context) (*sql.DB, error) {
	conn, err := sql.Open("postgres", DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getTablesFromDatabase(ctx context.Context, conn *sql.DB) ([]Table, error) {
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

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tables []Table

	for rows.Next() {
		var tableName string
		var columnsJSON string

		if err := rows.Scan(&tableName, &columnsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var columns []Column
		if err := json.Unmarshal([]byte(columnsJSON), &columns); err != nil {
			return nil, fmt.Errorf("failed to unmarshal columns JSON: %w", err)
		}

		tables = append(tables, Table{
			Name:    tableName,
			Columns: columns,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tables, nil
}

func QueryLLMChart(ctx context.Context, userPrompt string) (*Chart, error) {
	availableTypes := make([]string, len(RegistradedCharts))
	i := 0
	for _, chart := range RegistradedCharts {
		availableTypes[i] = `"` + chart.Name + `"`
		i += 1
	}
	availableTypesString := strings.Join(availableTypes, " or ")

	prompt := fmt.Sprintf("Given the following user request: '%s', suggest the most appropriate chart type from the following options: %s. Only reply with the chart type name and nothing else with no markdown formatting or characters.", userPrompt, availableTypesString)

	answer, err := llms.GenerateFromSinglePrompt(ctx, LLM, prompt)
	if err != nil {
		log.Fatal(err)
	}

	return StringToChart(answer)
}

func QueryLLMSqlStatement(ctx context.Context, userPrompt string, tables []Table, chart *Chart) string {
	// Convert tables schema to a readable format for the LLM
	var schemaDesc strings.Builder
	schemaDesc.WriteString("Database Schema:\n")
	for _, table := range tables {
		schemaDesc.WriteString(fmt.Sprintf("Table '%s':\n", table.Name))
		for _, col := range table.Columns {
			fkInfo := ""
			if col.ForeignKey != nil {
				fkInfo = fmt.Sprintf(" (Foreign Key to %s.%s)",
					col.ForeignKey.TargetTable,
					col.ForeignKey.TargetColumn)
			}
			pkInfo := ""
			if col.IsPrimaryKey {
				pkInfo = " (Primary Key)"
			}
			schemaDesc.WriteString(fmt.Sprintf("  - %s (%s)%s%s\n",
				col.Name, col.Type, pkInfo, fkInfo))
		}
	}

	// Construct the full prompt
	prompt := fmt.Sprintf(`Given the following database schema and chart requirements, generate a PostgreSQL query for this request: "%s"

Expected result Struct formated on Go Lang: %s

Current Database modeling:
%s

Rules:
1. The output must follow the expected struct format
2. The query must return exactly the columns needed for the chart type
3. Use appropriate joins if needed
4. Only return the SQL query, no explanations and no markdown formatting or comments
5. Respect the words of the tables and columns names are case sensitive and all the tables and columns names should be in double quotes
6. Generate a query compatible with the postgres database`,
		userPrompt,
		chart.Schema,
		schemaDesc,
	)

	// Query the LLM
	answer, err := llms.GenerateFromSinglePrompt(ctx, LLM, prompt)
	if err != nil {
		log.Printf("Error generating SQL query: %v", err)
		return ""
	}

	// Clean up the response and remove any markdown code blocks if present
	answer = strings.TrimSpace(answer)
	answer = strings.TrimPrefix(answer, "SQL Query: ")
	answer = strings.TrimPrefix(answer, "```sql")
	answer = strings.TrimPrefix(answer, "```")
	answer = strings.TrimSuffix(answer, "```")
	answer = strings.TrimSpace(answer)

	return answer
}

type ChartEntry map[string]any
type ChartData []ChartEntry

func RunSqlQuery(ctx context.Context, conn *sql.DB, sql string, chart *Chart) ChartData {
	// Execute the query
	rows, err := conn.QueryContext(ctx, sql)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil
	}
	defer rows.Close()

	// Get column names and types
	columns, err := rows.ColumnTypes()
	if err != nil {
		log.Printf("Error getting column info: %v", err)
		return nil
	}

	// Create result slice
	var result ChartData

	// Iterate through rows
	for rows.Next() {
		// Create a slice of interface{} to hold the row values
		scanArgs := make([]interface{}, len(columns))
		for i := range scanArgs {
			scanArgs[i] = new(interface{})
		}

		// Scan the row into the interfaces
		if err := rows.Scan(scanArgs...); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Create a new entry for this row
		entry := make(ChartEntry)

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
		log.Printf("Error iterating rows: %v", err)
	}

	return result
}

var RegistradedCharts []Chart

type ChartType int

const (
	BarChartType ChartType = iota
	LineChartType
)

type Chart struct {
	Name   string
	Type   ChartType
	Schema string
}

func StringToChart(value string) (*Chart, error) {
	for _, chart := range RegistradedCharts {
		if strings.ToLower(chart.Name) == strings.ToLower(value) {
			return &chart, nil
		}
	}

	return nil, fmt.Errorf("unsupported chart type: %s", value)
}

func StringToChartType(value string) (ChartType, error) {
	for _, chart := range RegistradedCharts {
		if strings.ToLower(chart.Name) == strings.ToLower(value) {
			return chart.Type, nil
		}
	}

	return -1, fmt.Errorf("unsupported chart type: %s", value)
}

func ChartTypeToString(value ChartType) string {
	for _, chart := range RegistradedCharts {
		if chart.Type == value {
			return chart.Name
		}
	}

	return "Line"
}

func init() {
	RegistradedCharts = []Chart{
		{
			Name: "Bar",
			Type: BarChartType,
			Schema: `type BarChart struct {
	Category string
	Value    float32
}`,
		},
		{
			Name: "Line",
			Type: LineChartType,
			Schema: `type LineChart struct {
	Category string
	Value    float32
}`,
		},
	}
}
