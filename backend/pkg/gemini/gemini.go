package gemini

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"mastercard-backend/internal/config"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewClient creates a new Gemini client
func NewClient() (*Client, error) {
	if config.AppConfig.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.AppConfig.GeminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(config.AppConfig.GeminiModel)

	// Convert config types to the required types for the genai model
	temp := float32(config.AppConfig.GeminiTemperature)
	maxTokens := int32(config.AppConfig.GeminiMaxTokens)
	model.Temperature = &temp
	model.MaxOutputTokens = &maxTokens

	return &Client{
		client: client,
		model:  model,
	}, nil
}

// Close closes the Gemini client
func (c *Client) Close() error {
	return c.client.Close()
}

// GenerateSQL generates SQL query from natural language using Gemini
func (c *Client) GenerateSQL(naturalLanguageQuery string, schemaContext string, conversationHistory []string) (string, error) {
	ctx := context.Background()

	// Build the prompt with schema context and conversation history
	prompt := buildPrompt(naturalLanguageQuery, schemaContext, conversationHistory)

	// Generate response
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate SQL: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Extract SQL from response
	sqlQuery := extractSQLFromResponse(resp.Candidates[0].Content.Parts[0].(genai.Text))

	return sqlQuery, nil
}

// GenerateAnalysis generates conversational analysis and insights about query results
func (c *Client) GenerateAnalysis(userQuery string, sqlQuery string, queryResults string, resultFormat string, conversationHistory []string) (string, error) {
	ctx := context.Background()

	// Build the analysis prompt
	prompt := buildAnalysisPrompt(userQuery, sqlQuery, queryResults, resultFormat, conversationHistory)

	// Generate response
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate analysis: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Extract analysis text from response
	analysis := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	analysis = strings.TrimSpace(analysis)

	return analysis, nil
}

// buildAnalysisPrompt constructs the prompt for generating conversational analysis
func buildAnalysisPrompt(userQuery string, sqlQuery string, queryResults string, resultFormat string, history []string) string {
	var prompt strings.Builder

	prompt.WriteString("You are a helpful data analyst assistant. Your task is to provide conversational analysis and insights about query results.\n\n")
	prompt.WriteString("You should:\n")
	prompt.WriteString("1. Analyze the query results and provide meaningful insights\n")
	prompt.WriteString("2. Explain what the data shows in a conversational, natural way\n")
	prompt.WriteString("3. Identify patterns, trends, or interesting findings\n")
	prompt.WriteString("4. Answer follow-up questions about the data\n")
	prompt.WriteString("5. Be conversational and friendly, like ChatGPT\n")
	prompt.WriteString("6. If asked about seasonality, trends, or 'why', provide analytical explanations\n")
	prompt.WriteString("7. Write in a natural, engaging style\n\n")

	if len(history) > 0 {
		prompt.WriteString("Previous conversation context:\n")
		for i, h := range history {
			if i < 5 { // Limit to last 5 messages for context
				prompt.WriteString(fmt.Sprintf("- %s\n", h))
			}
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("User's Question: ")
	prompt.WriteString(userQuery)
	prompt.WriteString("\n\n")

	prompt.WriteString("SQL Query Executed: ")
	prompt.WriteString(sqlQuery)
	prompt.WriteString("\n\n")

	prompt.WriteString("Query Results (Format: ")
	prompt.WriteString(resultFormat)
	prompt.WriteString("):\n")
	prompt.WriteString(queryResults)
	prompt.WriteString("\n\n")

	prompt.WriteString("Provide a conversational analysis and insights about these results. ")
	prompt.WriteString("Be natural, engaging, and helpful. If the user asked a specific question, answer it directly. ")
	prompt.WriteString("If they asked for analysis or insights, provide meaningful commentary about the data.\n\n")
	prompt.WriteString("Your analysis:")

	return prompt.String()
}

// buildPrompt constructs the prompt for Gemini
func buildPrompt(query string, schemaContext string, history []string) string {
	var prompt strings.Builder

	prompt.WriteString("You are a SQL expert assistant. Your task is to convert natural language queries into PostgreSQL SQL statements.\n\n")
	prompt.WriteString("Database Schema:\n")
	prompt.WriteString(schemaContext)
	prompt.WriteString("\n\n")

	if len(history) > 0 {
		prompt.WriteString("Previous conversation context:\n")
		for i, h := range history {
			if i < 10 { // Limit to last 10 messages
				prompt.WriteString(fmt.Sprintf("- %s\n", h))
			}
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("Rules:\n")
	prompt.WriteString("1. Generate ONLY valid PostgreSQL SQL queries\n")
	prompt.WriteString("2. Use proper table and column names from the schema\n")
	prompt.WriteString("3. Always use SINGLE QUOTES (') for string literals, NEVER double quotes (\") - double quotes are only for identifiers\n")
	prompt.WriteString("4. For date ranges, use proper date functions (e.g., DATE_TRUNC, INTERVAL)\n")
	prompt.WriteString("5. For aggregations, use appropriate GROUP BY clauses\n")
	prompt.WriteString("6. Return ONLY the SQL query, no explanations or markdown formatting\n")
	prompt.WriteString("7. If the query is ambiguous, generate a reasonable interpretation\n")
	prompt.WriteString("8. Use proper JOINs when needed\n")
	prompt.WriteString("9. Limit results to reasonable sizes (use LIMIT when appropriate)\n")
	prompt.WriteString("10. Handle NULL values appropriately\n")
	prompt.WriteString("11. Example: WHERE merchant_city = 'Almaty' (correct) NOT WHERE merchant_city = \"Almaty\" (wrong)\n\n")

	prompt.WriteString("User Query: ")
	prompt.WriteString(query)
	prompt.WriteString("\n\n")
	prompt.WriteString("Generate the SQL query:")

	return prompt.String()
}

// extractSQLFromResponse extracts SQL query from Gemini's response
func extractSQLFromResponse(response genai.Text) string {
	sql := string(response)

	// Remove markdown code blocks if present
	sql = strings.TrimSpace(sql)
	if strings.HasPrefix(sql, "```sql") {
		sql = strings.TrimPrefix(sql, "```sql")
		sql = strings.TrimSuffix(sql, "```")
	} else if strings.HasPrefix(sql, "```") {
		sql = strings.TrimPrefix(sql, "```")
		sql = strings.TrimSuffix(sql, "```")
	}

	sql = strings.TrimSpace(sql)

	// More robustly remove a single pair of leading/trailing quotes if they exist.
	// This prevents `strings.Trim` from removing the quote from a date or string literal at the end of the query.
	if (strings.HasPrefix(sql, "'") && strings.HasSuffix(sql, "'")) ||
		(strings.HasPrefix(sql, "\"") && strings.HasSuffix(sql, "\"")) {
		sql = sql[1 : len(sql)-1]
	}

	// Fix double quotes in string literals (replace with single quotes)
	// This regex finds double-quoted strings and replaces them with single quotes
	// Pattern: "([^"]*)" but only if it's not a table/column identifier
	// We'll do a simple replacement for common cases
	sql = fixDoubleQuotesInSQL(sql)

	return sql
}

// fixDoubleQuotesInSQL replaces double quotes with single quotes in string literals
// PostgreSQL uses single quotes for string literals, double quotes for identifiers
func fixDoubleQuotesInSQL(sql string) string {
	// Find all double-quoted strings that appear to be string literals
	// Pattern: looks for = "value", IN ("value"), LIKE "value", etc.
	// This regex matches double-quoted strings after operators
	re := regexp.MustCompile(`(=\s*|IN\s*\(|LIKE\s+|ILIKE\s+|,\s*)"([^"]+)"`)
	sql = re.ReplaceAllString(sql, `${1}'${2}'`)

	// Also handle standalone double-quoted strings in WHERE clauses
	// Pattern: WHERE column = "value"
	re2 := regexp.MustCompile(`(WHERE|AND|OR|HAVING)\s+(\w+)\s*=\s*"([^"]+)"`)
	sql = re2.ReplaceAllString(sql, `${1} ${2} = '${3}'`)

	// Simple fallback: replace remaining double quotes in value positions
	// This is a last resort for edge cases
	re3 := regexp.MustCompile(`"([^"]+)"`)
	sql = re3.ReplaceAllStringFunc(sql, func(match string) string {
		// Only replace if it's not likely an identifier (doesn't start with uppercase or contain schema.table)
		content := strings.Trim(match, `"`)
		if !strings.Contains(content, ".") && !regexp.MustCompile(`^[A-Z][a-zA-Z0-9_]*$`).MatchString(content) {
			return `'` + content + `'`
		}
		return match
	})

	return sql
}

// GetSchemaContext returns the database schema context as a string
func GetSchemaContext() string {
	return `CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    card_no VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    process_date DATE NOT NULL,
    trx_amount_usd DECIMAL(15, 2),
    trx_amount_eur DECIMAL(15, 2),
    trx_amount_local DECIMAL(15, 2),
    trx_cnt_usd INTEGER DEFAULT 0,
    trx_cnt_eur INTEGER DEFAULT 0,
    trx_cnt_local INTEGER DEFAULT 0,
    interchange_fee DECIMAL(15, 2),
    merch_name VARCHAR(255),
    agg_merch_name VARCHAR(255),
    issuer_code VARCHAR(50),
    issuer_country VARCHAR(100),
    bin6_code VARCHAR(6),
    acquirer_code VARCHAR(50),
    acquirer_country VARCHAR(100),
    trx_type VARCHAR(50),
    trx_direction VARCHAR(10) CHECK (trx_direction IN ('plus', 'minus')),
    mcc VARCHAR(10),
    mcc_group VARCHAR(100),
    input_mode VARCHAR(50),
    wallet_type VARCHAR(50),
    product_type VARCHAR(50),
    authorization_status VARCHAR(50),
    authorization_response_code VARCHAR(10),
    location_id VARCHAR(100),
    location_city VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

Common query patterns:
- Date filtering: WHERE date >= '2024-01-01' AND date <= '2024-03-31' (for Q1 2024)
- Merchant filtering: WHERE merch_name = 'Merchant Name' OR agg_merch_name = 'Merchant Name'
- Location filtering: WHERE location_city = 'Almaty' (use SINGLE quotes for strings)
- Aggregations: SUM(trx_amount_usd), SUM(trx_amount_eur), SUM(trx_amount_local), COUNT(*), AVG(trx_amount_usd)
- Top N queries: ORDER BY column DESC LIMIT N
- Grouping: GROUP BY location_city, merch_name, mcc_group, etc.
- Type filtering: WHERE trx_type = 'POS'
- Direction filtering: WHERE trx_direction = 'plus' (outgoing) OR trx_direction = 'minus' (incoming)
- Status filtering: WHERE authorization_status = 'approved' OR authorization_status = 'declined'
- IMPORTANT: Always use SINGLE QUOTES (') for string values, NEVER double quotes (")
- IMPORTANT: Use 'date' column for transaction date, 'process_date' for processing date
- IMPORTANT: Amount columns: trx_amount_usd, trx_amount_eur, trx_amount_local (not transaction_amount_kzt)`
}
