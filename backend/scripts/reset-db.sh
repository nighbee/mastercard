#!/bin/bash

# A shell script to reset the Dockerized PostgreSQL database.
# This is the equivalent of reset-db.ps1 for use in Bash-like shells (e.g., Git Bash, WSL).

set -e # Exit immediately if a command exits with a non-zero status.

echo "Stopping and removing all containers, networks, and volumes..."
docker compose down --volumes

echo "Starting a fresh database container..."
docker compose up -d postgres

echo "Waiting for the database to be ready..."

# Load environment variables from the root .env file to get DB credentials
if [ -f ../.env ]; then
  # Use `set -a` to automatically export all variables defined in the sourced file.
  # This is a robust way to load .env files, correctly handling spaces in values.
  set -a
  source ../.env
  set +a
fi

# Wait for PostgreSQL to become available
echo "Waiting for database initialization and data loading to complete..."

# The postgres container runs init scripts and then becomes ready.
# We will poll `pg_isready` in a loop. This is more reliable than watching logs,
# as it confirms the database is actively accepting connections.
retries=30
for i in $(seq 1 $retries); do
  # Use `docker exec` to run pg_isready inside the container
  if docker exec mastercard_postgres pg_isready -U "${DB_USER:-mastercard_user}" -d "${DB_NAME:-mastercard_db}" -q; then
    echo "Database is ready!"
    echo "----------------------------------------------------"
    echo "Database reset complete! You can now start the backend."
    echo "Run: go run cmd/server/main.go"
    exit 0
  fi
  echo "Waiting for database... ($i/$retries)"
  sleep 2
done

echo "Error: Database did not become ready in time."
exit 1