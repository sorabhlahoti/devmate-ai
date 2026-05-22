# Working Guide

## Run CLI

```bash
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
## Ask Question
```
go run ./cmd/devmate ask "explain goroutine in simple words"
```

## Explain Error From File

### Create sample error file:

```
echo "open /data/events.log: permission denied" > error.log
```
### Run

```
go run ./cmd/devmate explain-error --file error.log
```

OR

```
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

## Run Tests

```
go test ./...
```

## Build Binary
```
go build -o bin/devmate ./cmd/devmate
```

## Run Built Binary

### Windows:

```bash
./bin/devmate.exe ask "explain goroutine"
./bin/devmate.exe save --title "Docker fix" --body "Check volume permission"
./bin/devmate.exe list
```

### Linux/Mac:

```bash
./bin/devmate ask "explain goroutine"
./bin/devmate save --title "Docker fix" --body "Check volume permission"
./bin/devmate list
```