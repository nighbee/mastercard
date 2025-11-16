Архитектура системы Mastercard NLP-to-SQL Chatbot
Общая архитектура (3-слойная)
┌─────────────────────────────────────────────────────────┐│                    FRONTEND (React)                      ││  - React + TypeScript + Tailwind CSS                    ││  - React Query для управления состоянием                ││  - API клиент для взаимодействия с backend              │└──────────────────────┬──────────────────────────────────┘                       │ HTTP/REST API┌──────────────────────▼──────────────────────────────────┐│                 BACKEND (Golang + Fiber)                ││                                                          ││  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ ││  │  Handlers    │  │  Middleware  │  │   Services   │ ││  │  (HTTP)      │  │  (Auth/RBAC)  │  │  (Business)  │ ││  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ ││         │                 │                  │         ││  ┌──────▼─────────────────▼──────────────────▼───────┐ ││  │              Models (GORM)                        │ ││  └───────────────────────────────────────────────────┘ ││                                                          ││  ┌──────────────────────────────────────────────────┐  ││  │         Gemini Client (NLP Engine)              │  ││  │  - SQL Generation                                │  ││  │  - Analysis Generation                           │  ││  └──────────────────────────────────────────────────┘  │└──────────────────────┬──────────────────────────────────┘                       │┌──────────────────────▼──────────────────────────────────┐│              DATABASE (PostgreSQL)                       ││  - Transactions table (транзакционные данные)          ││  - Users, Roles, Permissions (RBAC)                     ││  - Conversations, Messages (история чатов)              ││  - Audit Logs (аудит действий)                          ││  - Row-Level Security (RLS) policies                   │└──────────────────────────────────────────────────────────┘
Детальная архитектура компонентов
1. Слой Handlers (HTTP обработчики)
Расположение: backend/internal/handlers/
query_handler.go — обработка NLP-запросов
conversation_handler.go — управление разговорами
auth_handler.go — аутентификация (login, register, refresh)
admin_handler.go — администрирование (пользователи, метрики, аудит)
Роль: принимает HTTP-запросы, валидирует входные данные, вызывает сервисы, возвращает ответы.
2. Слой Middleware (промежуточное ПО)
Расположение: backend/internal/middleware/
auth.go — JWT-аутентификация
rbac.go — проверка ролей и прав
cors.go — CORS
logger.go — логирование запросов
Роль: обработка запросов до handlers (аутентификация, авторизация, логирование).
3. Слой Services (бизнес-логика)
Расположение: backend/internal/services/
query_service.go — основная логика NLP-to-SQL
conversation_service.go — управление разговорами
auth_service.go — логика аутентификации
audit_service.go — аудит действий
Роль: бизнес-логика, взаимодействие с БД и внешними сервисами.
4. Слой Models (модели данных)
Расположение: backend/internal/models/models.go
Модели GORM:
User, Role, Permission — RBAC
Transaction — транзакционные данные
Conversation, Message — история чатов
AuditLog — аудит
Роль: структуры данных и маппинг на таблицы БД.
5. Gemini Client (NLP Engine)
Расположение: backend/pkg/gemini/gemini.go
Роль: взаимодействие с Google Gemini API для генерации SQL и анализа.
Процесс обработки NLP-запроса (детально)
Шаг 1: HTTP-запрос от фронтенда
POST /api/v1/query{  "query": "какие города есть в таблице транзакций?",  "conversation_id": 123}
Шаг 2: Middleware (аутентификация и авторизация)
1. AuthMiddleware (auth.go)   ├─ Извлекает JWT токен из заголовка Authorization   ├─ Валидирует токен   ├─ Загружает пользователя из БД с ролью   └─ Сохраняет user в c.Locals("user")2. RBAC проверка (если требуется)   └─ RequireRole/RequirePermission проверяет права доступа
Шаг 3: QueryHandler.ExecuteQuery
// query_handler.gofunc (h *QueryHandler) ExecuteQuery(c *fiber.Ctx) error {    userID := c.Locals("userID").(uint)  // Извлечен из JWT    var req QueryRequest    c.BodyParser(&req)        // Вызывает сервис    message, err := h.queryService.ExecuteQuery(userID, req.Query, req.ConversationID)    return c.JSON(fiber.Map{"message": message})}
Шаг 4: QueryService.ExecuteQuery (основная логика)
// query_service.gofunc (s *QueryService) ExecuteQuery(...) {    // 4.1 Загрузка истории разговора    history := загрузить_последние_10_сообщений(conversationID)        // 4.2 Генерация SQL через Gemini    sqlQuery := s.geminiClient.GenerateSQL(query, schemaContext, history)        // 4.3 Валидация SQL (только SELECT)    if !isValidReadOnlyQuery(sqlQuery) {        return error    }        // 4.4 Выполнение SQL в БД    result, format := s.executeSQL(sqlQuery)        // 4.5 Генерация анализа через Gemini    analysis := s.geminiClient.GenerateAnalysis(query, sqlQuery, result, format, history)        // 4.6 Сохранение в БД    message := создать_Message_запись(...)    database.DB.Create(&message)        return message}
NLP с генерацией SQL (Gemini)
Процесс генерации SQL
1. Подготовка контекста схемы БД
// gemini.gofunc GetSchemaContext() string {    return `    CREATE TABLE transactions (        id SERIAL PRIMARY KEY,        card_no VARCHAR(255) NOT NULL,        date DATE NOT NULL,        location_city VARCHAR(100),        trx_amount_usd DECIMAL(15, 2),        ...    );        Common query patterns:    - Date filtering: WHERE date >= '2024-01-01'    - Location filtering: WHERE location_city = 'Almaty'    - IMPORTANT: Always use SINGLE QUOTES (') for strings    `}
2. Построение промпта для Gemini
func buildPrompt(query, schemaContext, history []string) string {    prompt := `    You are a SQL expert assistant.        Database Schema:    [полная схема таблицы transactions]        Previous conversation context:    - "покажи транзакции за январь"    - "сколько их всего?"        Rules:    1. Generate ONLY valid PostgreSQL SQL    2. Use SINGLE QUOTES (') for string literals    3. Return ONLY SQL, no explanations        User Query: "какие города есть в таблице транзакций?"        Generate the SQL query:    `    return prompt}
3. Вызов Gemini API
func (c *Client) GenerateSQL(...) (string, error) {    prompt := buildPrompt(...)        // Вызов Google Gemini API    resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))        // Извлечение SQL из ответа    sqlQuery := extractSQLFromResponse(resp)        // Пост-обработка: исправление двойных кавычек    sqlQuery = fixDoubleQuotesInSQL(sqlQuery)        return sqlQuery}
4. Пост-обработка SQL
// Исправление двойных кавычек на одинарныеfunc fixDoubleQuotesInSQL(sql string) string {    // WHERE location_city = "Almaty"     // → WHERE location_city = 'Almaty'    re := regexp.MustCompile(`(=\s*|IN\s*\(|LIKE\s+)"([^"]+)"`)    sql = re.ReplaceAllString(sql, `${1}'${2}'`)    return sql}
Результат: SELECT DISTINCT location_city FROM transactions WHERE location_city IS NOT NULL
Генерация анализа (Gemini)
Процесс генерации анализа
После выполнения SQL:
// query_service.goif result != "" && resultFormat != "error" {    // Подготовка контекста для анализа    analysisHistory := загрузить_последние_5_сообщений(conversationID)        // Вызов Gemini для генерации анализа    analysis := geminiClient.GenerateAnalysis(        userQuery:    "какие города есть в таблице транзакций?",        sqlQuery:      "SELECT DISTINCT location_city FROM transactions...",        queryResults: `[{"location_city":"Almaty"},{"location_city":"Astana"}]`,        resultFormat: "table",        history:      analysisHistory    )}
Промпт для анализа:
func buildAnalysisPrompt(...) string {    return `    You are a helpful data analyst assistant.        User's Question: "какие города есть в таблице транзакций?"        SQL Query Executed: SELECT DISTINCT location_city FROM transactions...        Query Results:    [{"location_city":"Almaty"},{"location_city":"Astana"}]        Provide a conversational analysis and insights about these results.    Be natural, engaging, and helpful.    `}
Результат анализа: "В таблице транзакций представлены два города: Алматы и Астана. Это основные финансовые центры Казахстана..."
Взаимодействие с базой данных
1. Подключение к БД
// database.gofunc Connect() error {    dsn := fmt.Sprintf(        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",        config.DBHost, config.DBUser, config.DBPassword,         config.DBName, config.DBPort, config.DBSSLMode,    )        DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{...})    return err}
2. Выполнение SQL-запроса
// query_service.go - executeSQL()func (s *QueryService) executeSQL(sqlQuery string) (string, string, error) {    // Получение raw SQL connection (не GORM)    sqlDB, err := database.DB.DB()        // Создание контекста с таймаутом    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)    defer cancel()        // Выполнение запроса    rows, err := sqlDB.QueryContext(ctx, sqlQuery)    defer rows.Close()        // Получение колонок    columns, _ := rows.Columns()        // Сканирование результатов    var results []map[string]interface{}    for rows.Next() {        values := make([]interface{}, len(columns))        valuePtrs := make([]interface{}, len(columns))        for i := range values {            valuePtrs[i] = &values[i]        }        rows.Scan(valuePtrs...)                // Преобразование в map        row := make(map[string]interface{})        for i, col := range columns {            row[col] = values[i]        }        results = append(results, row)    }        // Определение формата результата    format := "table"    if len(results) == 1 && len(columns) == 1 {        format = "text"  // Одно значение    }        // Конвертация в JSON    resultJSON, _ := json.Marshal(results)    return string(resultJSON), format, nil}
3. Сохранение сообщения в БД
// query_service.gomessage := models.Message{    UserMessage:     query,    SQLQuery:        &sqlQuery,    ResultData:      &resultJSON,  // JSONB поле    ResultFormat:    &format,    Analysis:        &analysis,    ExecutionTimeMs: &executionTime,    ConversationID:  conversationID,}// GORM автоматически создает INSERT запросdatabase.DB.Create(&message)
SQL запрос (выполняется GORM):
INSERT INTO messages (    conversation_id, user_message, sql_query,     result_data, result_format, analysis, execution_time_ms) VALUES (    123, 'какие города...', 'SELECT DISTINCT...',    '{"location_city":"Almaty"}', 'table', 'В таблице...', 150);
4. Row-Level Security (RLS)
-- migrations/005_create_row_level_security.sqlCREATE POLICY transactions_read_policy ON transactionsFOR SELECTUSING (    -- Admin и Manager видят все    EXISTS (SELECT 1 FROM users u JOIN roles r ON u.role_id = r.id            WHERE u.id = current_setting('app.current_user_id')::integer            AND r.name IN ('admin', 'manager'))    OR    -- Analyst и Viewer видят только свои данные (если нужно)    ...);
Установка user_id перед запросом:
// database.gofunc SetCurrentUserID(userID uint) error {    return DB.Exec("SET app.current_user_id = ?", userID).Error}
Полный поток данных (пример)
1. Пользователь: "какие города есть в таблице транзакций?"   ↓2. Frontend → POST /api/v1/query   ↓3. AuthMiddleware → проверка JWT → userID = 5   ↓4. QueryHandler.ExecuteQuery   ↓5. QueryService.ExecuteQuery   ├─ Загрузка истории: [] (новый разговор)   ├─ Gemini.GenerateSQL()   │  ├─ buildPrompt() → промпт с схемой БД   │  ├─ API вызов к Google Gemini   │  └─ extractSQL() → "SELECT DISTINCT location_city FROM transactions..."   ├─ Валидация: isValidReadOnlyQuery() → ✅   ├─ executeSQL()   │  ├─ database.DB.DB() → raw connection   │  ├─ QueryContext() → выполнение SQL   │  ├─ rows.Scan() → сканирование результатов   │  └─ JSON marshal → '[{"location_city":"Almaty"},{"location_city":"Astana"}]'   ├─ Gemini.GenerateAnalysis()   │  ├─ buildAnalysisPrompt() → промпт с результатами   │  ├─ API вызов к Google Gemini   │  └─ "В таблице транзакций представлены два города..."   └─ database.DB.Create(&message) → сохранение в БД   ↓6. Response → Frontend   {     "message": {       "id": 456,       "user_message": "какие города...",       "sql_query": "SELECT DISTINCT...",       "result_data": "[{\"location_city\":\"Almaty\"}...]",       "result_format": "table",       "analysis": "В таблице транзакций...",       "execution_time_ms": 150     }   }   ↓7. Frontend отображает:   - Анализ: "В таблице транзакций..."   - Таблица: Almaty, Astana
Особенности архитектуры
Разделение ответственности: Handlers → Services → Models → Database
Безопасность:
JWT-аутентификация
RBAC на уровне middleware и handlers
RLS в PostgreSQL
Валидация SQL (только SELECT)
Масштабируемость:
Stateless backend
Connection pooling через GORM
Таймауты для SQL-запросов
Надежность:
Graceful shutdown
Обработка ошибок на всех уровнях
Аудит действий
Логирование
Система спроектирована как модульная, безопасная и масштабируемая платформа для NLP-to-SQL запросов.