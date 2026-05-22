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