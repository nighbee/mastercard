# PowerShell script to reset the database
Write-Host "Stopping and removing containers and volumes..." -ForegroundColor Yellow
docker-compose -f ../docker-compose.yml down -v

Write-Host "Starting fresh database container..." -ForegroundColor Green
docker-compose -f ../docker-compose.yml up -d postgres

Write-Host "Waiting for database to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host "Database reset complete! You can now start the backend." -ForegroundColor Green
Write-Host "Run: go run cmd/server/main.go" -ForegroundColor Cyan

