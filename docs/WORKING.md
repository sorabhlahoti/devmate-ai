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
```

### Linux/Mac:

```bash
./bin/devmate ask "explain goroutine"
```   