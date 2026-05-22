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

## Commands

```bash
devmate
devmate --help
devmate version
devmate ask "explain goroutine in simple words"
devmate explain-error --file error.log
devmate explain-error -f error.log
devmate save --title "Docker permission fix" --body "Check mounted volume permission"
devmate list
devmate search "docker"
```

## Planned Features
- Real AI provider integration
- Sync with backend
- Generate public runbooks
- Git commit message helper
- Docker support
- Swagger API doc
