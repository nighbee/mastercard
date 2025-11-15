# Database Migrations

This directory contains SQL migration files that create and configure the database schema.

## Migration Files

1. **001_create_transactions_table.sql**
   - Creates the main `transactions` table with all 28 fields
   - Creates indexes for common query patterns
   - Optimized for date range queries and filtering

2. **002_create_users_and_roles_tables.sql**
   - Creates `roles`, `permissions`, `role_permissions`, and `users` tables
   - Sets up RBAC (Role-Based Access Control) structure
   - Inserts default roles: admin, manager, analyst, viewer
   - Inserts default permissions and assigns them to roles

3. **003_create_conversations_and_messages_tables.sql**
   - Creates `conversations` table for chat sessions
   - Creates `messages` table for individual chat messages
   - Supports conversation branching via `parent_branch_id`
   - Includes full-text search indexes for searching conversations

4. **004_create_audit_logs_table.sql**
   - Creates `audit_logs` table for comprehensive audit trail
   - Tracks all user actions, queries, and system events
   - Indexed for efficient querying by user, timestamp, and action

5. **005_create_row_level_security.sql**
   - Enables Row Level Security (RLS) on sensitive tables
   - Creates RLS policies for transactions, conversations, messages, and audit_logs
   - Ensures users can only access data they're authorized to see

## Running Migrations

### Using Docker Compose (Automatic)
When you start the PostgreSQL container with `docker-compose up`, migrations in this directory are automatically executed in alphabetical order.

### Manual Execution
To run migrations manually:

```bash
# Connect to PostgreSQL
psql -h localhost -U mastercard_user -d mastercard_db

# Run migrations in order
\i 001_create_transactions_table.sql
\i 002_create_users_and_roles_tables.sql
\i 003_create_conversations_and_messages_tables.sql
\i 004_create_audit_logs_table.sql
\i 005_create_row_level_security.sql
```

## Database Schema Overview

### Core Tables
- **transactions**: Main transaction data (28 fields)
- **users**: User accounts and authentication
- **roles**: User roles (admin, manager, analyst, viewer)
- **permissions**: Resource-action permissions
- **role_permissions**: Many-to-many role-permission mapping
- **conversations**: Chat conversation sessions
- **messages**: Individual messages with queries and results
- **audit_logs**: Comprehensive audit trail

### Security Features
- Row Level Security (RLS) enabled on sensitive tables
- Role-based access control (RBAC)
- Encrypted password storage (bcrypt)
- Audit logging for all actions

## Notes

- All tables include `created_at` and `updated_at` timestamps
- Foreign keys are properly set up with CASCADE/SET NULL as appropriate
- Indexes are created for common query patterns
- Full-text search is enabled for conversations and messages
- RLS policies ensure data isolation between users

