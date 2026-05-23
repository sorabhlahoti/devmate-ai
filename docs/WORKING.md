
# Working Guide

## Requirements

- Go installed
- Git installed
- Docker Desktop installed(for windows users) 
- swag CLI installed for Swagger generation

Install swag:

```
go install github.com/swaggo/swag/cmd/swag@latest
```

Add to PATH (Linux/Mac):

```
export PATH=$PATH:$(go env GOPATH)/bin
```

Add to PATH (Windows PowerShell):

```
$env:Path += ";$(go env GOPATH)\bin"
```

Check swag:

```
swag --version
```

## Run CLI

```
go run ./cmd/devmate
```

## Show Help

```
go run ./cmd/devmate --help
```

## Show Version

```
go run ./cmd/devmate version
```

## Ask a Question

```
go run ./cmd/devmate ask "explain goroutine in simple words"
```

## Explain Error from File

### Create a sample error file

**Linux/Mac:**
```
echo "open /data/events.log: permission denied" > error.log
```

**Windows PowerShell:**
```
echo "open /data/events.log: permission denied" > error.log
```

### Run

```
go run ./cmd/devmate explain-error --file error.log
```

or

```
go run ./cmd/devmate explain-error -f error.log
```

## Notes Management

### Save a Note

```
go run ./cmd/devmate save --title "Docker permission fix" --body "Check mounted volume permission"
```

### List Notes

```
go run ./cmd/devmate list
```

### Search Notes

```
go run ./cmd/devmate search "docker"
```

## Config Path

``` 
go run ./cmd/devmate config path
```

## Set API URL

```
go run ./cmd/devmate config set api_url http://localhost:8080
```

## Get API URL

``` 
go run ./cmd/devmate config get api_url
```

## Generate Swagger Docs

**Direct command:**

```
swag init -g cmd/api/main.go -o internal/swaggerdocs --outputTypes go,json --parseInternal
```

**Swagger UI:**

http://localhost:8080/swagger/index.html

**Swagger JSON:**

```
curl http://localhost:8080/swagger/doc.json
```

## Run API

```
go run ./cmd/api
```

### API Health Check

```
curl http://localhost:8080/health
```

### Create Note Through API

**Linux/Mac:**
```
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"Docker permission fix","body":"Check mounted volume permission"}'
```

**Windows PowerShell:**
```
curl.exe -X POST http://localhost:8080/api/v1/notes -H "Content-Type: application/json" -d "{\"title\":\"Docker permission fix\",\"body\":\"Check mounted volume permission\"}"
```

### List Notes Through API

```
curl http://localhost:8080/api/v1/notes
```

### Search Notes Through API

```
curl "http://localhost:8080/api/v1/notes/search?q=docker"
```

## Sync Local Notes With API

### Start API Server

**Linux/Mac:**
```
export DEVMATE_DB_PATH="$PWD/devmate-api.db"
go run ./cmd/api
```

**Windows PowerShell:**
```
$env:DEVMATE_DB_PATH="$PWD\devmate-api.db"
go run ./cmd/api
```

### Configure and Sync Notes

```
go run ./cmd/devmate config set api_url http://localhost:8080
go run ./cmd/devmate save --title "Sync test note" --body "This note should go to backend"
go run ./cmd/devmate sync
```

#### Run Sync Again

```
go run ./cmd/devmate sync
```

**Expected Output:**

```
No unsynced notes found ✅
```

## Run Tests

```
go test ./...
```

## Build CLI Binary

```
go build -o bin/devmate ./cmd/devmate
```

### Build API Binary

```
go build -o bin/devmate-api ./cmd/api
```

# Docker API Commands

## Build Docker Image

```
docker build -t devmate-api:local .
```

## Run Docker Container

```
docker run --rm -p 8080:8080 devmate-api:local
```

## Run With Docker Compose

```
docker compose up --build
```

## Health Check

```
curl http://localhost:8080/health
```

**Windows PowerShell:**
```
curl.exe http://localhost:8080/health
```

## Create Note Through Dockerized API

**Linux/Mac:**
```
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"Dockerized API","body":"DevMate API is running inside Docker"}'
```

**Windows PowerShell:**
```
curl.exe -X POST http://localhost:8080/api/v1/notes -H "Content-Type: application/json" -d "{\"title\":\"Dockerized API\",\"body\":\"DevMate API is running inside Docker\"}"
```

## List Notes

```
curl http://localhost:8080/api/v1/notes
```

## Sync Local CLI Notes To Dockerized API

**Terminal 1 (API Server):**

```
docker compose up --build
```

**Terminal 2 (CLI Commands):**

```
go run ./cmd/devmate config set api_url http://localhost:8080
go run ./cmd/devmate save --title "Docker sync test" --body "Local note synced to Docker API"
go run ./cmd/devmate sync
```

## Stop Docker Compose

```
docker compose down
```

## Stop And Delete Docker Volume

```
docker compose down -v
```
  
> **Note:** Use `-v` only if you want to delete saved API notes. This will permanently remove all note data.