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

### Backend API

- Health check endpoint
- Create note endpoint
- List notes endpoint
- Search notes endpoint

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
```

## API Endpoints

```http
GET  /health
POST /api/v1/notes
GET  /api/v1/notes
GET  /api/v1/notes/search?q=docker
```

## Planned Features

- CLI sync with backend
- Real AI provider integration
- Generate public runbooks
- Git commit message helper
- Docker support
- Swagger API docs
- Kubernetes manifests
- Terraform starter
- CI/CD pipeline templates