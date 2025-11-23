# Mastercard NLP-to-SQL Analytics Chatbot Platform

An enterprise-grade AI-powered conversational analytics system that enables users to query transactional databases using natural language. Built for the Mastercard case study, this platform translates natural language queries into SQL statements using Google Gemini AI, executes them securely, and presents results in multiple formats (tables, charts, text).

## 📋 Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [API Reference](#api-reference)
- [Security Features](#security-features)
- [Development Guide](#development-guide)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## 🎯 Overview

The Mastercard NLP-to-SQL Analytics Chatbot is a full-stack application designed to democratize data access for business users. Instead of writing complex SQL queries, users can ask questions in natural language like:

- "Show me total transactions for Q1 2024"
- "What are the top 5 merchants by revenue in Kazakhstan?"
- "List all cities in the transactions table"

The system uses Google Gemini AI to intelligently convert these questions into valid PostgreSQL queries, execute them safely, and provide conversational analysis of the results.

### Use Case

Built as part of the Mastercard case study competition to demonstrate:
- **Technical Excellence**: Modern architecture with Go backend and React frontend
- **AI Integration**: Sophisticated NLP-to-SQL conversion with context awareness
- **Security**: Enterprise-grade authentication, authorization, and audit logging
- **User Experience**: Intuitive chat interface with voice input support
- **Business Value**: Reduces analyst workload and democratizes data access

## ✨ Key Features

### Natural Language Processing
- **AI-Powered Query Translation**: Google Gemini converts natural language to SQL
- **Context-Aware**: Maintains conversation history for follow-up questions
- **Multi-Format Results**: Automatically formats results as tables, charts, or text
- **Conversational Analysis**: AI-generated insights and explanations of query results

### Security & Compliance
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Role-Based Access Control (RBAC)**: Four role levels (Admin, Manager, Analyst, Viewer)
- **Row-Level Security (RLS)**: Database-level data isolation policies
- **Audit Logging**: Complete trail of all queries and user actions
- **SQL Injection Prevention**: Query validation and parameterization
- **Read-Only Queries**: Only SELECT statements allowed, no data modification

### User Interface
- **Modern Chat Interface**: Clean, responsive design with real-time updates
- **Voice Input**: Speech-to-text query input support
- **Conversation Management**: Create, branch, search, and organize conversations
- **Interactive Results**: Sortable tables, interactive charts, exportable data
- **Multi-Language Support**: Interface supports English, Russian, and Kazakh

### Data Analysis
- **28-Field Transaction Schema**: Comprehensive payment transaction data model
- **Advanced Filtering**: Date ranges, merchants, locations, card types, etc.
- **Aggregations**: SUM, COUNT, AVG, MIN, MAX with grouping
- **Time Series Analysis**: Quarter, month, year-based analytics
- **Merchant Analytics**: Revenue analysis, transaction patterns

## 🏗️ Architecture

### System Architecture (3-Tier)

```
┌─────────────────────────────────────────────────────────┐
│                 FRONTEND (React + TypeScript)           │
│  - Vite build system with hot reload                    │
│  - shadcn/ui components with Tailwind CSS               │
│  - React Query for state management & caching           │
│  - React Router for navigation                          │
│  - Recharts for data visualization                      │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTP/REST API (JSON)
┌──────────────────────▼──────────────────────────────────┐
│              BACKEND (Go + Fiber Framework)             │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │  Handlers    │  │  Middleware  │  │   Services   │ │
│  │  (HTTP API)  │  │  (Auth/RBAC) │  │  (Business)  │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │
│         │                 │                  │         │
│  ┌──────▼─────────────────▼──────────────────▼───────┐ │
│  │              Models (GORM ORM)                    │ │
│  └───────────────────────────────────────────────────┘ │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │      Gemini Client (NLP-to-SQL Engine)          │  │
│  │  - SQL Generation with schema context            │  │
│  │  - Query analysis and explanation                 │  │
│  │  - Error handling and fallback mechanisms         │  │
│  └──────────────────────────────────────────────────┘  │
└──────────────────────┬──────────────────────────────────┘
                       │ PostgreSQL Protocol
┌──────────────────────▼──────────────────────────────────┐
│            DATABASE (PostgreSQL 14)                      │
│  - Transactions table (28 fields, indexed)              │
│  - Users, Roles, Permissions (RBAC system)              │
│  - Conversations, Messages (chat history)               │
│  - Audit Logs (compliance tracking)                     │
│  - Row-Level Security policies                          │
└──────────────────────────────────────────────────────────┘
```

### Request Flow

1. **User Input**: User types or speaks a query in the frontend
2. **Authentication**: JWT token validated by auth middleware
3. **Authorization**: RBAC middleware checks user permissions
4. **Handler**: Query handler receives and validates request
5. **Service Layer**: Query service orchestrates the process:
   - Loads conversation history for context
   - Calls Gemini API to generate SQL
   - Validates SQL (read-only, safe)
   - Executes query on database
   - Calls Gemini API to analyze results
   - Saves message to database
6. **Audit**: All actions logged to audit_logs table
7. **Response**: Results returned to frontend with SQL, data, and analysis

## 🛠️ Tech Stack

### Backend (Go)
- **Framework**: [Fiber v2](https://gofiber.io/) - Fast Express-inspired web framework
- **ORM**: [GORM](https://gorm.io/) - Go object-relational mapper
- **Database Driver**: PostgreSQL driver for GORM
- **Authentication**: [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) - JWT tokens
- **Password Hashing**: golang.org/x/crypto/bcrypt
- **AI/NLP**: [Google Generative AI Go SDK](https://github.com/google/generative-ai-go) - Gemini API
- **Configuration**: [godotenv](https://github.com/joho/godotenv) - Environment variables
- **Lines of Code**: ~2,735 lines

### Frontend (TypeScript/React)
- **Build Tool**: [Vite 5](https://vitejs.dev/) - Lightning-fast dev server and build
- **Framework**: [React 18](https://react.dev/) - UI library with hooks
- **Language**: [TypeScript 5](https://www.typescriptlang.org/) - Type-safe JavaScript
- **Routing**: [React Router 6](https://reactrouter.com/) - Client-side routing
- **State Management**: [@tanstack/react-query 5](https://tanstack.com/query) - Server state management
- **UI Components**: [shadcn/ui](https://ui.shadcn.com/) - High-quality component library
- **Styling**: [Tailwind CSS 3](https://tailwindcss.com/) - Utility-first CSS
- **Charts**: [Recharts 2](https://recharts.org/) - Composable charting library
- **Form Handling**: [React Hook Form 7](https://react-hook-form.com/) + [Zod](https://zod.dev/) - Type-safe forms
- **Icons**: [Lucide React](https://lucide.dev/) - Beautiful icon set
- **Lines of Code**: ~6,098 lines

### Database
- **RDBMS**: [PostgreSQL 14](https://www.postgresql.org/) - Advanced open-source database
- **Deployment**: Docker with Alpine Linux
- **Migrations**: SQL scripts (9 migration files)
- **Features**: JSONB columns, full-text search, row-level security, indexes

### Infrastructure
- **Containerization**: [Docker](https://www.docker.com/) & Docker Compose
- **Database Volume**: Persistent storage with named volumes
- **Networking**: Bridge network for service communication
- **Health Checks**: Automated container health monitoring

### Development Tools
- **Version Control**: Git + GitHub
- **Code Quality**: ESLint (frontend), Go fmt (backend)
- **API Testing**: cURL, Postman (examples provided)

## 🚀 Quick Start

### Prerequisites
- **Docker** and **Docker Compose** (for containerized setup)
- **Go 1.21+** (for backend development)
- **Node.js 18+** and **npm** (for frontend development)
- **Google Gemini API Key** (get from [Google AI Studio](https://makersuite.google.com/app/apikey))

### Option 1: Full Docker Setup (Recommended)

1. **Clone the repository**:
```bash
git clone https://github.com/nighbee/mastercard.git
cd mastercard
```

2. **Configure environment**:
```bash
# Create .env file in project root
cat > .env << EOF
# Database
DB_USER=mastercard_user
DB_PASSWORD=mastercard_pass
DB_NAME=mastercard_db
DB_PORT=5432

# Application
APP_PORT=8080
EOF

# Configure backend
cd backend
cp env.template .env
# Edit .env and set your GEMINI_API_KEY
```

3. **Start all services**:
```bash
cd ..
docker-compose up --build
```

This starts:
- PostgreSQL on port 5432 (with automatic migrations)
- Backend API on port 8080
- Frontend dev server on port 5173

4. **Access the application**:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

### Option 2: Local Development Setup

#### 1. Database Setup

```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Wait for database to initialize (10-15 seconds)
# Or use the reset script
./reset-database.ps1  # Windows
# or
bash backend/scripts/reset-db.sh  # Linux/Mac
```

#### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Configure environment
cp env.template .env
# Edit .env and set required variables

# Run development server
go run cmd/server/main.go

# Or build and run
go build -o bin/server cmd/server/main.go
./bin/server
```

Backend will start on http://localhost:8080

#### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will start on http://localhost:5173

### First Steps

1. **Register an account**: Navigate to http://localhost:5173/register
2. **Login**: Use your credentials to access the dashboard
3. **Start chatting**: Try example queries like:
   - "Show all transactions"
   - "What cities are in the database?"
   - "Total transactions by merchant"

## 📁 Project Structure

```
mastercard/
├── backend/                      # Go backend service
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Application entry point
│   ├── internal/                # Internal packages (not importable)
│   │   ├── config/              # Configuration management
│   │   │   └── config.go        # Environment variable loading
│   │   ├── database/            # Database connection setup
│   │   │   └── database.go      # GORM initialization
│   │   ├── handlers/            # HTTP request handlers
│   │   │   ├── auth_handler.go  # Login, register, refresh, profile
│   │   │   ├── query_handler.go # NLP query execution
│   │   │   ├── conversation_handler.go  # Conversation CRUD
│   │   │   └── admin_handler.go # User management, audit logs
│   │   ├── middleware/          # HTTP middleware
│   │   │   ├── auth.go          # JWT validation
│   │   │   ├── rbac.go          # Role-based access control
│   │   │   ├── cors.go          # CORS configuration
│   │   │   └── logger.go        # Request logging
│   │   ├── models/              # Database models
│   │   │   └── models.go        # GORM structs (User, Transaction, etc.)
│   │   ├── services/            # Business logic
│   │   │   ├── auth_service.go  # Authentication logic
│   │   │   ├── query_service.go # NLP-to-SQL orchestration
│   │   │   ├── conversation_service.go  # Conversation management
│   │   │   └── audit_service.go # Audit logging
│   │   └── utils/               # Utility functions
│   │       ├── jwt.go           # JWT token generation/validation
│   │       └── password.go      # Password hashing
│   ├── pkg/                     # Public packages (importable)
│   │   └── gemini/              # Google Gemini API client
│   │       └── gemini.go        # NLP-to-SQL and analysis
│   ├── migrations/              # Database migrations (9 files)
│   │   ├── 001_create_transactions_table.sql
│   │   ├── 002_create_users_and_roles_tables.sql
│   │   ├── 003_create_conversations_and_messages_tables.sql
│   │   ├── 004_create_audit_logs_table.sql
│   │   ├── 005_create_row_level_security.sql
│   │   ├── 006_load_transaction_data.sql
│   │   ├── 007_sampeldata.sql
│   │   ├── 008_add_analysis_to_messages.sql
│   │   └── 009_update_rbac_roles.sql
│   ├── scripts/                 # Helper scripts
│   │   ├── reset-db.sh          # Reset database script
│   │   └── run-migrations.sh    # Manual migration runner
│   ├── Dockerfile              # Backend container image
│   ├── go.mod                  # Go module definition
│   ├── go.sum                  # Go dependency checksums
│   └── env.template            # Environment variable template
│
├── frontend/                    # React TypeScript frontend
│   ├── src/
│   │   ├── components/          # React components
│   │   │   ├── ui/             # shadcn/ui base components
│   │   │   ├── ConversationList.tsx
│   │   │   ├── MessageBubble.tsx
│   │   │   ├── ResultsViewer.tsx  # Table/chart display
│   │   │   └── VoiceInputModal.tsx
│   │   ├── contexts/            # React contexts
│   │   │   └── AuthContext.tsx  # Authentication state
│   │   ├── hooks/               # Custom hooks
│   │   │   └── use-toast.ts
│   │   ├── lib/                 # Utilities
│   │   │   ├── api.ts          # API client with axios
│   │   │   └── utils.ts        # Helper functions
│   │   ├── pages/               # Page components
│   │   │   ├── Landing.tsx     # Landing page
│   │   │   ├── Login.tsx       # Login page
│   │   │   ├── Register.tsx    # Registration page
│   │   │   ├── Dashboard.tsx   # Main chat interface
│   │   │   ├── Profile.tsx     # User profile
│   │   │   ├── Admin.tsx       # Admin panel
│   │   │   └── NotFound.tsx    # 404 page
│   │   ├── App.tsx             # Main app component
│   │   ├── main.tsx            # Entry point
│   │   └── index.css           # Global styles
│   ├── public/                 # Static assets
│   ├── package.json            # npm dependencies
│   ├── tsconfig.json           # TypeScript configuration
│   ├── vite.config.ts          # Vite build configuration
│   └── tailwind.config.ts      # Tailwind CSS configuration
│
├── doc/                         # Documentation
│   ├── overview.md             # Architecture overview (Russian)
│   ├── caseoverview.md         # Case study requirements
│   ├── userflow.md             # User flow documentation
│   ├── fieldsofdatabase.md     # Database field descriptions
│   ├── RBAC.md                 # Role-based access control
│   └── Mastercard-NLP-SQL-PRD(1).md  # Product requirements
│
├── docker-compose.yml           # Docker services definition
├── .gitignore                  # Git ignore patterns
├── reset-database.ps1          # PowerShell database reset script
└── README.md                   # This file
```

## 🗄️ Database Schema

### Core Tables

#### 1. transactions (28 fields)
The main table storing all payment transaction data.

**Key Fields**:
- `id` - Primary key
- `card_no` - Masked card number
- `date` - Transaction date (indexed)
- `process_date` - Processing date (indexed)
- `trx_amount_usd`, `trx_amount_eur`, `trx_amount_local` - Amounts in different currencies
- `merch_name` - Merchant name (indexed)
- `location_city` - Transaction city (indexed)
- `issuer_country` - Card issuing country (indexed)
- `mcc` - Merchant Category Code (indexed)
- `authorization_status` - Approved/declined status (indexed)
- And 16 more fields for comprehensive transaction analysis

**Indexes**: date, process_date, merch_name, location_city, issuer_country, mcc, authorization_status

#### 2. users
User accounts with authentication credentials.

**Fields**:
- `id` - Primary key
- `email` - Unique user email
- `password_hash` - bcrypt hashed password
- `full_name` - User's full name
- `role_id` - Foreign key to roles
- `is_active` - Account status
- `last_login` - Last login timestamp

#### 3. roles
User role definitions for RBAC.

**Default Roles**:
- **admin** - Full system access, user management
- **manager** - View all data, manage users (except delete)
- **analyst** - Execute queries, view own conversations
- **viewer** - Read-only access to data

#### 4. permissions
Resource-action permission mappings.

**Structure**:
- `resource` - What is being accessed (users, queries, audit_logs)
- `action` - What action (read, write, delete)
- `conditions` - JSON conditions for fine-grained control

#### 5. conversations
Chat conversation sessions.

**Fields**:
- `id` - Primary key
- `user_id` - Owner of the conversation
- `title` - Conversation title
- `parent_id` - For conversation branching
- `is_active` - Active/archived status
- `created_at`, `updated_at` - Timestamps

#### 6. messages
Individual messages in conversations.

**Fields**:
- `id` - Primary key
- `conversation_id` - Parent conversation
- `user_message` - Original natural language query
- `sql_query` - Generated SQL statement
- `result_data` - Query results (JSONB)
- `result_format` - Format type (table/chart/text)
- `analysis` - AI-generated analysis text
- `execution_time_ms` - Query execution time
- `error_message` - Error details if query failed

#### 7. audit_logs
Comprehensive audit trail for compliance.

**Fields**:
- `id` - Primary key
- `user_id` - User who performed action
- `action` - Action type (login, query, admin_action)
- `resource` - Resource affected
- `details` - Additional context (JSONB)
- `ip_address` - User's IP address
- `user_agent` - Browser/client information
- `created_at` - Timestamp

### Database Features

- **Migrations**: 9 SQL migration files executed in order
- **Indexes**: Strategic indexes on commonly queried fields
- **JSONB**: Schema-less JSON storage for flexible data
- **Row-Level Security**: PostgreSQL RLS policies for data isolation
- **Foreign Keys**: Referential integrity between tables
- **Timestamps**: Automatic created_at/updated_at tracking

For detailed field descriptions, see [Database Fields Documentation](./doc/fieldsofdatabase.md).

## 📡 API Reference

Full API documentation with all endpoints, request/response formats, and examples.

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

All auth endpoints return JWT tokens for secure API access.

#### POST /auth/register
Register a new user account.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
```

**Response** (201):
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "role_id": 3
  }
}
```

#### POST /auth/login
Authenticate and receive access + refresh tokens.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response** (200):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": {"id": 3, "name": "analyst"}
  }
}
```

#### POST /auth/refresh
Refresh expired access token.

#### GET /auth/profile
Get current user profile (requires authentication).

**Headers**: `Authorization: Bearer <token>`

### Query Endpoints

#### POST /query
Execute natural language query (requires authentication).

**Request**:
```json
{
  "query": "Show me total transactions for Q1 2024",
  "conversation_id": 1
}
```

**Response** (200):
```json
{
  "message": {
    "id": 42,
    "user_message": "Show me total transactions for Q1 2024",
    "sql_query": "SELECT COUNT(*) FROM transactions WHERE date >= '2024-01-01' AND date < '2024-04-01'",
    "result_data": "[{\"total\": 15234}]",
    "result_format": "text",
    "analysis": "There were 15,234 transactions in Q1 2024...",
    "execution_time_ms": 145
  }
}
```

### Conversation Endpoints

- **POST /conversations** - Create new conversation
- **GET /conversations** - List user's conversations
- **GET /conversations/:id** - Get conversation with messages
- **PUT /conversations/:id** - Update conversation
- **DELETE /conversations/:id** - Delete conversation
- **POST /conversations/:id/branch** - Create conversation branch
- **GET /conversations/search?q=keyword** - Search conversations

### Admin Endpoints (Manager/Admin only)

- **GET /admin/users** - List all users
- **POST /admin/users** - Create new user
- **PUT /admin/users/:id** - Update user
- **DELETE /admin/users/:id** - Delete user (admin only)
- **GET /admin/audit-logs** - View audit logs
- **GET /admin/metrics** - System metrics

**Error Responses**: Standard HTTP codes (400, 401, 403, 404, 500)

For complete API documentation with all parameters and examples, see the [Backend README](./backend/README.md).

## 🔐 Security Features

Comprehensive security implementation following industry best practices.

### Authentication & Authorization

**JWT Token System**:
- **Access Tokens**: Short-lived (15 min), used for API requests
- **Refresh Tokens**: Long-lived (7 days), obtain new access tokens
- **Token Rotation**: Refresh tokens rotate on use
- **Bcrypt Hashing**: Password storage with cost factor 10

**Role-Based Access Control (RBAC)**:
- **Admin**: Full system access, user management, configuration
- **Manager**: View all data, manage users (create/update), audit logs
- **Analyst**: Execute queries, manage own conversations
- **Viewer**: Read-only data access

### Database Security

- **Row-Level Security (RLS)**: PostgreSQL policies restrict data by user role
- **Prepared Statements**: GORM prevents SQL injection with parameterized queries
- **Query Validation**: Only SELECT statements allowed
- **Connection Pooling**: Prevents resource exhaustion
- **SSL/TLS**: Encrypted database connections in production

### Query Safety

**SQL Injection Prevention**:
- Regex validation blocks dangerous keywords (DROP, DELETE, UPDATE, INSERT)
- Query parser ensures read-only operations
- Gemini-generated SQL validated before execution

**Resource Protection**:
- 30-second timeout on all queries
- Maximum 10,000 rows per result (configurable)
- Schema isolation (public schema only)

### Audit & Compliance

**Comprehensive Logging**:
- Every action logged with user ID, action type, resource, IP, user agent, timestamp
- Query history stored for review
- Immutable audit logs
- GDPR-compliant data export/deletion

### Network Security

- **CORS**: Configurable allowed origins, methods, headers
- **HTTPS**: TLS encryption required in production
- **Environment Variables**: Secrets never in code

## 🔧 Development Guide

Complete guide for local development and contribution.

### Backend Development

**Setup**:
```bash
cd backend
go mod download
cp env.template .env
# Edit .env with your configuration
go run cmd/server/main.go
```

**Hot Reload** (recommended):
```bash
go install github.com/cosmtrek/air@latest
air  # Auto-rebuilds on file changes
```

**Code Structure**:
- **Handlers**: Thin HTTP layer, delegates to services
- **Services**: Business logic and orchestration
- **Models**: GORM database structs
- **Middleware**: Auth, RBAC, CORS, logging
- **Utils**: Reusable helpers

**Testing API**:
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!","full_name":"Test User"}'

# Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!"}' | jq -r '.access_token')

# Execute query
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query":"Show all transactions","conversation_id":1}'
```

### Frontend Development

**Setup**:
```bash
cd frontend
npm install
echo "VITE_API_URL=http://localhost:8080/api/v1" > .env
npm run dev
```

**Commands**:
- `npm run dev` - Development server with hot reload
- `npm run build` - Production build
- `npm run preview` - Preview production build
- `npm run lint` - Lint code

**Adding UI Components**:
```bash
npx shadcn-ui@latest add button
npx shadcn-ui@latest add dialog
```

**Code Guidelines**:
- Use TypeScript for type safety
- Tailwind utility classes for styling
- React Query for server state
- Functional components with hooks

### Database Management

**Migrations**:
```bash
# Automatic on container start
docker-compose up -d postgres

# Manual execution
docker exec -i mastercard_postgres psql -U mastercard_user -d mastercard_db < backend/migrations/001_create_transactions_table.sql
```

**Database Console**:
```bash
docker exec -it mastercard_postgres psql -U mastercard_user -d mastercard_db

# Inside psql:
\dt              # List tables
\d transactions  # Describe table
\q               # Quit
```

**Reset Database**:
```bash
./reset-database.ps1  # Windows
bash backend/scripts/reset-db.sh  # Linux/Mac
```

## 🚢 Deployment

Production deployment guide for Docker-based infrastructure.

### Docker Production Setup

1. **Configure Environment**:
```bash
# Project root .env
cat > .env << EOF
DB_USER=prod_user
DB_PASSWORD=strong_prod_password_change_this
DB_NAME=mastercard_db
DB_PORT=5432
APP_PORT=8080
EOF

# Backend .env
cd backend
cp env.template .env
# Set production values:
# - Strong JWT_SECRET (32+ chars)
# - Valid GEMINI_API_KEY
# - DB_SSLMODE=require
# - Specific CORS_ALLOWED_ORIGINS
```

2. **Build and Deploy**:
```bash
docker-compose up -d --build
```

3. **Verify**:
```bash
docker-compose ps
docker-compose logs -f backend
curl http://localhost:8080/health
```

### Production Security Checklist

- [ ] All default passwords changed
- [ ] JWT_SECRET is cryptographically random (32+ characters)
- [ ] Database SSL enabled (`DB_SSLMODE=require`)
- [ ] CORS restricted to production frontend URL only
- [ ] HTTPS enabled (use reverse proxy with SSL certificate)
- [ ] Environment files not in version control
- [ ] Gemini API key has appropriate quota limits
- [ ] Database backups scheduled
- [ ] Monitoring and logging configured

### Scaling Recommendations

**Backend Scaling**:
- Horizontal: Add backend containers behind load balancer (nginx, HAProxy)
- Stateless design allows unlimited scaling
- Use Redis for session/query caching

**Database Scaling**:
- PostgreSQL read replicas for read-heavy workloads
- PgBouncer for connection pooling
- Table partitioning for large transaction tables
- Regular VACUUM and ANALYZE

**Frontend Scaling**:
- Build static assets: `npm run build`
- Serve via CDN (CloudFlare, AWS CloudFront)
- Enable gzip/brotli compression
- Cache static assets aggressively

### Reverse Proxy (nginx example)

```nginx
server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 🐛 Troubleshooting

Common issues and solutions for development and deployment.

### Database Connection Errors

**Symptom**: `failed to connect to database`, authentication errors

**Solutions**:
```bash
# 1. Reset database completely
./reset-database.ps1  # Windows
bash backend/scripts/reset-db.sh  # Linux/Mac

# 2. Check PostgreSQL is running
docker ps | grep postgres

# 3. View logs for errors
docker logs mastercard_postgres

# 4. Verify credentials match
cat backend/.env | grep DB_
docker-compose config | grep -A5 postgres

# 5. Test connection manually
docker exec -it mastercard_postgres psql -U mastercard_user -d mastercard_db
```

### Backend Won't Start

**Symptom**: `Failed to load config`, Gemini API errors

**Solutions**:
```bash
# 1. Verify .env exists
ls -la backend/.env

# 2. Check Go version
go version  # Must be 1.21+

# 3. Verify Gemini API key
# Test at: https://makersuite.google.com/app/apikey

# 4. Check environment loading
cd backend && go run cmd/server/main.go
# Look for specific error messages

# 5. Reinstall dependencies
rm go.sum
go mod download
```

### Frontend API Errors

**Symptom**: CORS errors, 401 Unauthorized, network failures

**Solutions**:
```bash
# 1. Verify backend is running
curl http://localhost:8080/health

# 2. Check CORS configuration
# In backend/.env:
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# 3. Clear browser data
# In browser console:
localStorage.clear()
sessionStorage.clear()

# 4. Verify API URL
cat frontend/.env
# Should be: VITE_API_URL=http://localhost:8080/api/v1

# 5. Check browser console for specific errors
# Press F12 → Console tab
```

### Query Execution Fails

**Symptom**: Queries timeout, return errors, or produce invalid SQL

**Solutions**:

1. **Gemini API Issues**:
   - Verify API key is valid and not rate-limited
   - Check quota at Google AI Studio
   - Review Gemini model name in .env

2. **SQL Validation Errors**:
   - Check backend logs for generated SQL
   - Verify schema context is correct
   - May need to adjust AI prompt configuration in Gemini client files

3. **Database Timeouts**:
   - Increase `QUERY_TIMEOUT` in backend/.env
   - Add indexes for slow queries
   - Check `pg_stat_activity` for long-running queries

4. **Permission Errors**:
   - Verify user has correct role
   - Check RLS policies in migration 005
   - Review RBAC permissions

### JWT Token Issues

**Symptom**: Constant re-authentication, token invalid errors

**Solutions**:
```bash
# 1. Check JWT configuration
cat backend/.env | grep JWT
# JWT_SECRET must be consistent
# JWT_ACCESS_TOKEN_EXPIRY=15m
# JWT_REFRESH_TOKEN_EXPIRY=168h

# 2. Clear tokens and re-login
# In browser console:
localStorage.clear()

# 3. Verify token format
# Should be: "Bearer <token>" in Authorization header

# 4. Check server time is correct
date  # On server
# Time skew can cause token validation failures
```

### Migration Errors

**Symptom**: Tables missing, schema incorrect

**Solutions**:
```bash
# 1. Check which migrations ran
docker exec mastercard_postgres psql -U mastercard_user -d mastercard_db -c "\dt"

# 2. View migration logs
docker logs mastercard_postgres | grep -i migration

# 3. Run specific migration manually
docker exec -i mastercard_postgres psql -U mastercard_user -d mastercard_db \
  < backend/migrations/001_create_transactions_table.sql

# 4. Nuclear option: Full reset
./reset-database.ps1
```

### Performance Issues

**Symptom**: Slow queries, high CPU/memory usage

**Solutions**:
```bash
# 1. Check query execution times
# Look at execution_time_ms in API responses

# 2. Monitor database
docker exec mastercard_postgres psql -U mastercard_user -d mastercard_db \
  -c "SELECT * FROM pg_stat_activity;"

# 3. Review slow queries
docker exec mastercard_postgres psql -U mastercard_user -d mastercard_db \
  -c "SELECT query, mean_exec_time FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;"

# 4. Add indexes for common queries
# Edit migrations to add indexes on frequently queried columns

# 5. Monitor container resources
docker stats
```

### Debug Mode

**Enable detailed logging**:
```bash
# Backend logs
docker-compose logs -f backend

# Database logs
docker logs -f mastercard_postgres

# Frontend console
# Open browser DevTools (F12) → Console tab

# Health check
curl -v http://localhost:8080/health
```

## 📚 Documentation

- **[Architecture Overview](./doc/overview.md)**: Detailed system architecture (Russian)
- **[Case Study Requirements](./doc/caseoverview.md)**: Original case study brief
- **[User Flow](./doc/userflow.md)**: Complete user journey documentation
- **[Database Fields](./doc/fieldsofdatabase.md)**: Transaction table field descriptions
- **[RBAC Model](./doc/RBAC.md)**: Role-based access control matrix
- **[Product Requirements](./doc/Mastercard-NLP-SQL-PRD(1).md)**: Full PRD document
- **[Backend README](./backend/README.md)**: Backend-specific documentation
- **[Frontend README](./frontend/README.md)**: Frontend-specific documentation

## 🤝 Contributing

This is an internal project for the Mastercard case study. For team members:

1. **Create a feature branch**: `git checkout -b feature/your-feature`
2. **Make changes**: Follow code structure guidelines
3. **Test locally**: Ensure all features work
4. **Commit**: Use clear commit messages
5. **Push**: `git push origin feature/your-feature`
6. **Pull Request**: Create PR for review

### Code Style

**Backend (Go)**:
- Use `go fmt` before committing
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Comment exported functions and types

**Frontend (TypeScript)**:
- Run `npm run lint` before committing
- Use TypeScript types (avoid `any`)
- Follow React best practices (hooks, functional components)

## 📄 License

Internal project for Mastercard case study competition.

---

## 📊 Project Statistics

- **Backend**: ~2,735 lines of Go code
- **Frontend**: ~6,098 lines of TypeScript/React code
- **Database**: 9 migration files, 7 tables, 28-field transaction schema
- **API Endpoints**: 20+ REST endpoints
- **Tech Stack**: 15+ major technologies
- **Documentation**: 6 comprehensive documentation files

---

**Built with ❤️ for the Mastercard Case Study**

For questions or issues, contact the development team.