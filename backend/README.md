# Mastercard Backend - NLP-to-SQL Platform

Golang backend service for the Mastercard NLP-to-SQL Analytics Chatbot Platform.

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 14+ (running via Docker Compose)
- Google Gemini API key

### Setup

1. **Reset and start the database:**
   ```powershell
   # From project root
   .\reset-database.ps1
   ```
   
   Or manually:
   ```powershell
   docker-compose down -v
   docker-compose up -d postgres
   ```

2. **Configure environment:**
   ```bash
   cd backend
   cp env.template .env
   # Edit .env and set your GEMINI_API_KEY
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080`

## 🐳 Docker Setup

### Start everything with Docker:

```bash
# From project root
docker-compose up --build
```

This will:
- Build the backend container
- Start PostgreSQL with migrations
- Start the backend service

### Start only database:

```bash
docker-compose up -d postgres
```

## 🔧 Database Connection Issues

If you encounter authentication errors:

1. **Reset the database:**
   ```powershell
   .\reset-database.ps1
   ```

2. **Check database is running:**
   ```powershell
   docker ps
   ```

3. **Test connection manually:**
   ```powershell
   docker exec -it mastercard_postgres psql -U mastercard_user -d mastercard_db
   ```

4. **Verify credentials in .env match docker-compose.yml:**
   - DB_USER=mastercard_user
   - DB_PASSWORD=mastercard_pass
   - DB_NAME=mastercard_db

## 📁 Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and setup
│   ├── handlers/            # HTTP request handlers
│   ├── middleware/         # HTTP middleware (auth, CORS, etc.)
│   ├── models/              # GORM database models
│   ├── services/            # Business logic services
│   └── utils/               # Utility functions (JWT, password hashing)
├── pkg/
│   └── gemini/              # Google Gemini API integration
├── migrations/              # SQL migration files
└── go.mod                   # Go module definition
```

## 🔧 Configuration

All configuration is done via environment variables in `.env`:

- **Database**: Connection settings for PostgreSQL
  - When running locally: `DB_HOST=localhost`
  - When running in Docker: `DB_HOST=postgres` (automatically set)
- **JWT**: Secret keys and token expiry times
- **Gemini**: API key and model configuration
- **CORS**: Allowed origins, methods, and headers
- **Query**: Timeout and result limits

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get tokens
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/profile` - Get current user profile (protected)

### Queries
- `POST /api/v1/query` - Execute natural language query (protected)

### Conversations
- `POST /api/v1/conversations` - Create new conversation (protected)
- `GET /api/v1/conversations` - List user's conversations (protected)
- `GET /api/v1/conversations/:id` - Get conversation with messages (protected)
- `PUT /api/v1/conversations/:id` - Update conversation (protected)
- `DELETE /api/v1/conversations/:id` - Delete conversation (protected)
- `POST /api/v1/conversations/:id/branch` - Create conversation branch (protected)
- `GET /api/v1/conversations/search?q=keyword` - Search conversations (protected)

### Admin
- `GET /api/v1/admin/users` - List all users (admin only)
- `GET /api/v1/admin/audit-logs` - View audit logs (admin only)
- `GET /api/v1/admin/metrics` - System metrics (admin only)

## 🔐 Authentication

The API uses JWT tokens for authentication:

1. Register or login to get `access_token` and `refresh_token`
2. Include the access token in requests: `Authorization: Bearer <token>`
3. Access tokens expire in 15 minutes (configurable)
4. Use refresh token to get a new access token

## 🧪 Testing

### Example: Register and Login

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Example: Execute Query

```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "query": "Show me total transactions for Q1 2024",
    "conversation_id": 1
  }'
```

## 🛠️ Development

### Build
```bash
go build -o bin/server cmd/server/main.go
```

### Run
```bash
./bin/server
```

### Run with hot reload (using air)
```bash
# Install air: go install github.com/cosmtrek/air@latest
air
```

## 📝 Notes

- The backend uses GORM for database operations
- Row Level Security (RLS) is enabled on sensitive tables
- All queries are logged to audit_logs table
- Only SELECT queries are allowed (no DROP, DELETE, UPDATE, INSERT)
- Google Gemini is used for NLP-to-SQL conversion

## 🔍 Troubleshooting

1. **Database connection error**: 
   - Run `.\reset-database.ps1` to reset the database
   - Make sure PostgreSQL is running (`docker ps`)
   - Check credentials in `.env` match `docker-compose.yml`

2. **Gemini API error**: 
   - Check that `GEMINI_API_KEY` is set in `.env`
   - Verify the API key is valid

3. **JWT errors**: 
   - Ensure `JWT_SECRET` is set and consistent
   - Check token expiry settings

4. **CORS errors**: 
   - Add your frontend URL to `CORS_ALLOWED_ORIGINS` in `.env`

5. **Migration errors**:
   - Check `docker logs mastercard_postgres` for migration errors
   - Ensure migration files are in `backend/migrations/`
