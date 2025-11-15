# PowerShell script to check database connection
param(
    [string]$Host = "localhost",
    [string]$Port = "5432",
    [string]$User = "mastercard_user",
    [string]$Password = "mastercard_pass",
    [string]$Database = "mastercard_db"
)

Write-Host "Testing database connection..." -ForegroundColor Yellow
Write-Host "Host: $Host" -ForegroundColor Cyan
Write-Host "Port: $Port" -ForegroundColor Cyan
Write-Host "User: $User" -ForegroundColor Cyan
Write-Host "Database: $Database" -ForegroundColor Cyan
Write-Host ""

$env:PGPASSWORD = $Password

try {
    $result = docker exec mastercard_postgres psql -h $Host -p $Port -U $User -d $Database -c "SELECT version();"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Database connection successful!" -ForegroundColor Green
        Write-Host $result
    } else {
        Write-Host "✗ Database connection failed!" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Make sure the database container is running:" -ForegroundColor Yellow
    Write-Host "  docker-compose up -d postgres" -ForegroundColor Cyan
}

