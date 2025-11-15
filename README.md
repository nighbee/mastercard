# Mastercard NLP-to-SQL Analytics Chatbot Platform

An AI-powered conversational analytics system that enables users to query transactional databases using natural language. The platform translates user queries into SQL statements, executes them securely, and presents results in multiple formats.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for backend development)
- Node.js 18+ and npm (for frontend development)

### 1. Database Setup

Start the PostgreSQL container:

```bash
docker-compose up -d postgres
```

This will:
- Start PostgreSQL 14 in a Docker container
- Automatically run all database migrations
- Create the database schema with all required tables

### 2. Environment Configuration

Copy the environment template and configure it:

```bash
cd backend
cp env.template .env
```

Edit `.env` and set:
- Database credentials (if different from defaults)
- `GEMINI_API_KEY` - Your Google Gemini API key
- `JWT_SECRET` - A secure random string for JWT tokens

### 3. Verify Database

Check if the database is running:

```bash
docker-compose ps
```

Connect to the database:

```bash
docker exec -it mastercard_postgres psql -U mastercard_user -d mastercard_db
```

## 📁 Project Structure

```
mastercard/
├── backend/              # Golang backend (to be implemented)
│   ├── migrations/      # SQL migration files
│   └── .env            # Environment configuration
├── frontend/            # React frontend (existing)
├── docker-compose.yml   # Docker services configuration
└── doc/                 # Documentation
```

## 🗄️ Database Schema

### Core Tables

1. **transactions** - Main transaction data (28 fields)
   - All transaction details including amounts, merchants, locations, etc.
   - Indexed for efficient querying

2. **users** - User accounts and authentication
   - Email, password hash, role assignment

3. **roles** - User roles (admin, manager, analyst, viewer)
   - Pre-configured with default roles

4. **permissions** - Resource-action permissions
   - Defines what actions can be performed on resources

5. **conversations** - Chat conversation sessions
   - Supports conversation branching

6. **messages** - Individual chat messages
   - Stores user queries, generated SQL, and results

7. **audit_logs** - Comprehensive audit trail
   - Tracks all user actions and queries

### Security Features

- **Row Level Security (RLS)** - Enabled on sensitive tables
- **Role-Based Access Control (RBAC)** - Granular permission system
- **Audit Logging** - Complete audit trail for compliance

## 🔧 Development

### Database Migrations

Migrations are automatically run when the PostgreSQL container starts. They are located in `backend/migrations/` and executed in alphabetical order.

To manually run migrations:

```bash
docker exec -i mastercard_postgres psql -U mastercard_user -d mastercard_db < backend/migrations/001_create_transactions_table.sql
```

### Database Connection

Default connection details:
- Host: `localhost`
- Port: `5432`
- Database: `mastercard_db`
- User: `mastercard_user`
- Password: `mastercard_pass`

## 📝 Next Steps

1. ✅ Database schema created
2. ✅ Docker PostgreSQL container configured
3. ⏳ Backend implementation (Golang + Fiber)
4. ⏳ Google Gemini integration
5. ⏳ Frontend-backend integration

## 📚 Documentation

- [User Flow](./doc/userflow.md)
- [Case Overview](./doc/caseoverview.md)
- [Product Requirements](./doc/Mastercard-NLP-SQL-PRD(1).md)
- [Database Fields](./doc/fieldsofdatabase.md)

## 🔐 Security Notes

- Change default database passwords in production
- Use strong JWT secrets
- Enable SSL for database connections in production
- Review and adjust RLS policies based on requirements

## 📄 License

Internal project for Mastercard case study.





test query: 
show total sum of transctions for q1