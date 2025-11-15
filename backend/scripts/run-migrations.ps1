# PowerShell script to run database migrations manually
# Usage: .\run-migrations.ps1

$ErrorActionPreference = "Stop"

# Load environment variables from .env file
if (Test-Path ".env") {
    Get-Content ".env" | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            Set-Variable -Name $name -Value $value -Scope Script
        }
    }
}

# Default values if not in .env
$DB_HOST = if ($DB_HOST) { $DB_HOST } else { "localhost" }
$DB_PORT = if ($DB_PORT) { $DB_PORT } else { "5432" }
$DB_USER = if ($DB_USER) { $DB_USER } else { "mastercard_user" }
$DB_NAME = if ($DB_NAME) { $DB_NAME } else { "mastercard_db" }
$DB_PASSWORD = if ($DB_PASSWORD) { $DB_PASSWORD } else { "mastercard_pass" }

Write-Host "Running database migrations..."
Write-Host "Database: $DB_NAME"
Write-Host "Host: ${DB_HOST}:${DB_PORT}"
Write-Host "User: $DB_USER"
Write-Host ""

# Set PGPASSWORD environment variable
$env:PGPASSWORD = $DB_PASSWORD

# Run migrations in order
Get-ChildItem -Path "migrations" -Filter "*.sql" | Sort-Object Name | ForEach-Object {
    Write-Host "Running: $($_.Name)"
    & psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $_.FullName
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Completed: $($_.Name)" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed: $($_.Name)" -ForegroundColor Red
        exit 1
    }
    Write-Host ""
}

Write-Host "All migrations completed successfully!" -ForegroundColor Green

