# Features

## Current Features

### CLI

- Root command
- Help command
- Version command
- Ask command with mock AI provider
- Explain error command from file
- Save developer notes locally
- List saved notes
- Search saved notes
- CLI config management
- Sync local notes with backend API

### Backend API

- Health check endpoint
- Create note endpoint
- List notes endpoint
- Search notes endpoint
- Swagger UI page

### DevOps

- Dockerfile for API
- Docker Compose setup
- SQLite Docker volume
- Container health check
- Makefile commands


## CLI Commands

```bash
# CLI Commands
devmate
devmate --help
devmate version
devmate ask "explain goroutine in simple words"
devmate explain-error --file examples/error.log
devmate save --title "Docker permission fix" --body "Check mounted volume permission"
devmate list
devmate search "docker"
devmate config set api_url http://localhost:8080
devmate config get api_url
devmate config path
devmate sync
```

## API Endpoints

```http
GET  /health
POST /api/v1/notes
GET  /api/v1/notes
GET  /api/v1/notes/search?q=docker
```

## Swagger

**Swagger UI:**
```http
http://localhost:8080/swagger/index.html
```

**Swagger JSON:**
```http
http://localhost:8080/swagger/doc.json
```

**Generate Swagger docs:**
```bash
swag init -g cmd/api/main.go -o internal/swaggerdocs --outputTypes go,json --parseInternal  
```


## Docker Commands

```bash
docker build -t devmate-api:local .
docker compose up --build
docker compose down
```

## Planned Features

- Real AI provider integration
- Generate public runbooks
- Git commit message helper
- Kubernetes manifests
- Terraform starter
- CI/CD pipeline templates