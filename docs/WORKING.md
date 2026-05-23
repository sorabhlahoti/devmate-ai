# Working Guide

## Run CLI

```bash
go run ./cmd/devmate
```

## Show Help

```bash
go run ./cmd/devmate --help
```

## Show Version

```bash
go run ./cmd/devmate version
```

## Ask a Question

```bash
go run ./cmd/devmate ask "explain goroutine in simple words"
```

## Explain Error from File

### Create a sample error file

```bash
echo "open /data/events.log: permission denied" > error.log
```

### Run

```bash
go run ./cmd/devmate explain-error --file error.log
```

or

```bash
go run ./cmd/devmate explain-error -f error.log
```

## Notes Management

### Save a Note

```bash
go run ./cmd/devmate save --title "Docker permission fix" --body "Check mounted volume permission"
```

### List Notes

```bash
go run ./cmd/devmate list
```

### Search Notes

```bash
go run ./cmd/devmate search "docker"
```

## Config Path
```bash 
go run ./cmd/devmate config path
```

## Set API URL
```bash
go run ./cmd/devmate config set api_url http://localhost:8080
```

## Get API URL
```bash 
go run ./cmd/devmate config get api_url
```

## Run API

```bash
go run ./cmd/api
```


### API Health Check

```bash
curl http://localhost:8080/health
```

### Create Note Through API

**PowerShell:**
```powershell
curl.exe -X POST http://localhost:8080/api/v1/notes `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Docker permission fix\",\"body\":\"Check mounted volume permission\"}"
```

**Linux/Mac:**
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"Docker permission fix","body":"Check mounted volume permission"}'
```

### List Notes Through API

```bash
curl http://localhost:8080/api/v1/notes
```

### Search Notes Through API

```bash
curl "http://localhost:8080/api/v1/notes/search?q=docker"
```

## Sync Local Notes With API

### Start API Server

**Terminal 1:**

```powershell
$env:DEVMATE_DB_PATH="$PWD\devmate-api.db"
go run ./cmd/api
```

### Configure and Sync Notes

**Terminal 2:**

```bash
go run ./cmd/devmate config set api_url http://localhost:8080
go run ./cmd/devmate save --title "Sync test note" --body "This note should go to backend"
go run ./cmd/devmate sync
```

#### Run Sync Again

```bash
go run ./cmd/devmate sync
```

**Expected Output:**

```
No unsynced notes found ✅
```

## Run Tests

```bash
go test ./...
```

## Build CLI Binary

```bash
go build -o bin/devmate ./cmd/devmate
```

### Build API Binary

```bash
go build -o bin/devmate-api ./cmd/api
```
# Docker API Commands

## Build Docker Image

```bash
docker build -t devmate-api:local .
```

## Run Docker Container

```bash
docker run --rm -p 8080:8080 devmate-api:local
```

## Run With Docker Compose

```bash
docker compose up --build
```

## Health Check

```bash
curl http://localhost:8080/health
```

**PowerShell:**
```powershell
curl.exe http://localhost:8080/health
```

## Create Note Through Dockerized API

**PowerShell:**
```powershell
curl.exe -X POST http://localhost:8080/api/v1/notes `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Dockerized API\",\"body\":\"DevMate API is running inside Docker\"}"
```

**Linux/Mac:**
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"Dockerized API","body":"DevMate API is running inside Docker"}'
```

## List Notes

```bash
curl http://localhost:8080/api/v1/notes
```

## Sync Local CLI Notes To Dockerized API

**Terminal 1:**
```bash
docker compose up --build
```

**Terminal 2:**
```bash
go run ./cmd/devmate config set api_url http://localhost:8080
go run ./cmd/devmate save --title "Docker sync test" --body "Local note synced to Docker API"
go run ./cmd/devmate sync
```

## Stop Docker Compose

```bash
docker compose down
```

## Stop And Delete Docker Volume

```bash
docker compose down -v
```

> **⚠️ Note:** Use `-v` only if you want to delete saved API notes. This will permanently remove all note data.