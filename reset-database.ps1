# PowerShell script to reset the database
# Run this from the project root directory

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Resetting Mastercard Database" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Stopping and removing containers and volumes..." -ForegroundColor Yellow
docker-compose down -v

Write-Host ""
Write-Host "Starting fresh database container..." -ForegroundColor Green
docker-compose up -d postgres

Write-Host ""
Write-Host "Waiting for database to be ready (this may take 10-15 seconds)..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0
$ready = $false

while ($attempt -lt $maxAttempts -and -not $ready) {
    Start-Sleep -Seconds 1
    $attempt++
    $result = docker exec mastercard_postgres pg_isready -U mastercard_user -d mastercard_db 2>&1
    if ($LASTEXITCODE -eq 0) {
        $ready = $true
        Write-Host "✓ Database is ready!" -ForegroundColor Green
    } else {
        Write-Host "." -NoNewline -ForegroundColor Gray
    }
}

if (-not $ready) {
    Write-Host ""
    Write-Host "✗ Database failed to start. Check logs with: docker logs mastercard_postgres" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Checking database connection..." -ForegroundColor Yellow
$env:PGPASSWORD = "mastercard_pass"
$testResult = docker exec mastercard_postgres psql -U mastercard_user -d mastercard_db -c "SELECT version();" 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Database connection successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Checking if migrations ran..." -ForegroundColor Yellow
    $tables = docker exec mastercard_postgres psql -U mastercard_user -d mastercard_db -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>&1
    if ($tables -match '\d+') {
        Write-Host "✓ Migrations completed. Found tables in database." -ForegroundColor Green
    }
} else {
    Write-Host "✗ Database connection failed!" -ForegroundColor Red
    Write-Host $testResult
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Database reset complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "You can now:" -ForegroundColor Yellow
Write-Host "  1. Start the backend: cd backend && go run cmd/server/main.go" -ForegroundColor Cyan
Write-Host "  2. Or start everything with Docker: docker-compose up" -ForegroundColor Cyan
Write-Host ""

