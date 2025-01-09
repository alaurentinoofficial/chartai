package charts_handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	charts_models "github.com/alaurentinoofficial/chartai/internal/charts/models"
	charts_repositories "github.com/alaurentinoofficial/chartai/internal/charts/repositories"
	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	core_services "github.com/alaurentinoofficial/chartai/internal/core/services"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	databases_handlers "github.com/alaurentinoofficial/chartai/internal/databases/handlers"
	databases_models "github.com/alaurentinoofficial/chartai/internal/databases/models"
)

type CreateChartRequest struct {
	Prompt     string `json:"prompt"`
	DatabaseId string `json:"databaseId" validate:"uuid"`
}

func CreateChartHandler(
	unitOfWork core_services.UnitOfWork,
	chartRepository charts_repositories.ChartRepository,
	structValidator *core_validations.StructValidator,
	getDatabaseById core_handlers.HandlerFunc[databases_handlers.GetDatabaseByIdRequest, databases_handlers.DatabaseResponse],
	llm core_services.LlmService,
) core_handlers.HandlerFunc[CreateChartRequest, ChartResponse] {
	return func(ctx context.Context, request CreateChartRequest) (*ChartResponse, error) {
		// Validates
		if err := structValidator.Validate(request); err != nil {
			return nil, err
		}

		// Begin transactionCtx
		transactionCtx, err := unitOfWork.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer transactionCtx.Rollback()

		database, err := getDatabaseById(
			ctx,
			databases_handlers.GetDatabaseByIdRequest{Id: request.DatabaseId},
		)
		if err != nil {
			return nil, err
		}

		// Query LLM the Chart Title
		title, err := queryLLMTitle(ctx, llm, request.Prompt)
		if err != nil {
			return nil, err
		}

		// Query LLM to get best Chart for this prompt
		chartKind, err := queryLLMChart(ctx, llm, request.Prompt)
		if err != nil {
			return nil, err
		}

		// Query LLM to generate a SQL Statement compatible with the Database
		sql, err := queryLLMSqlStatement(ctx, llm, request.Prompt, *database.Schema, chartKind)
		if err != nil {
			return nil, err
		}

		var chart charts_models.ChartEntity
		switch chartKind.(type) {
		case charts_models.BarChartEntity:
			chart = charts_models.NewBarChart(title, sql, database.Id)
			break
		default:
			chart = charts_models.NewLineChart(title, sql, database.Id)
			break
		}

		// CreateNewBarChart
		err = chartRepository.Create(transactionCtx, chart)
		if err != nil {
			return nil, err
		}

		// Commit changes
		err = transactionCtx.Commit()
		if err != nil {
			return nil, err
		}

		return NewChartResponse(chart), nil
	}
}

func queryLLMTitle(ctx context.Context, llm core_services.LlmService, userPrompt string) (string, error) {
	prompt := fmt.Sprintf("Given the following user request: '%s', suggest a title for this chart consise into a single short phrase. The output must be just the chart title, no markdown, no code, just text", userPrompt)

	answer, err := llm.SinglePrompt(ctx, prompt)
	if err != nil {
		return "", errors.New("Failed to process the chart title")
	}

	return answer, nil
}

func queryLLMChart(ctx context.Context, llm core_services.LlmService, userPrompt string) (charts_models.ChartEntity, error) {
	availableTypes := make([]string, charts_models.GetRegistredChartsLen())
	i := 0
	for _, chart := range charts_models.GetRegistredCharts() {
		availableTypes[i] = `"` + chart.GetType() + `"`
		i += 1
	}
	availableTypesString := strings.Join(availableTypes, " or ")

	prompt := fmt.Sprintf("Given the following user request: '%s', suggest the most appropriate chart type from the following options: %s. Only reply with the chart type name and nothing else with no markdown formatting or characters.", userPrompt, availableTypesString)

	answer, err := llm.SinglePrompt(ctx, prompt)
	if err != nil {
		return nil, errors.New("Failed to process the chart type")
	}

	return charts_models.StringToChart(answer)
}

func queryLLMSqlStatement(
	ctx context.Context,
	llm core_services.LlmService,
	userPrompt string,
	tables []databases_models.Table,
	chart charts_models.ChartEntity,
) (string, error) {
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

	schema := ""
	for _, column := range chart.GetSchema() {
		isRequred := ""

		if !column.Optional {
			isRequred = " NOT NULL"
		}

		schema += fmt.Sprintf(`- "%s" %s%s '%s'`, column.Name, column.Type, isRequred, column.Description)
	}

	// Construct the full prompt
	prompt := fmt.Sprintf(`Given the following database schema and chart requirements, generate a PostgreSQL query for this request: "%s"

Expected Columns of the Query are:
%s

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
		schema,
		schemaDesc,
	)

	// Query the LLM
	answer, err := llm.SinglePrompt(ctx, prompt)
	if err != nil {
		return "", errors.New("Failed to process the chart type")
	}

	// Clean up the response and remove any markdown code blocks if present
	answer = strings.TrimSpace(answer)
	answer = strings.TrimPrefix(answer, "SQL Query: ")
	answer = strings.TrimPrefix(answer, "```sql")
	answer = strings.TrimPrefix(answer, "```")
	answer = strings.TrimSuffix(answer, "```")
	answer = strings.TrimSpace(answer)

	return answer, nil
}
