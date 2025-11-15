#!/bin/bash

# Script to run database migrations manually
# Usage: ./run-migrations.sh

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values if not in .env
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-mastercard_user}
DB_NAME=${DB_NAME:-mastercard_db}
PGPASSWORD=${DB_PASSWORD:-mastercard_pass}

echo "Running database migrations..."
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo ""

# Run migrations in order
for migration in migrations/*.sql; do
    if [ -f "$migration" ]; then
        echo "Running: $migration"
        PGPASSWORD=$PGPASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
        echo "✓ Completed: $migration"
        echo ""
    fi
done

echo "All migrations completed successfully!"

