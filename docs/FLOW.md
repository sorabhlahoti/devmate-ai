# Application Flow

## Current CLI Flow

```text
User runs command
      |
      v
cmd/devmate/main.go
      |
      v
cli.Execute()
      |
      v
rootCmd.Execute()
      |
      v
Cobra checks command
      |
      v
Runs matching command
```
Example
```
go run ./cmd/devmate version
```
Flow:

main.go starts
      |
      v
cli.Execute() runs
      |
      v
Cobra finds version command
      |
      v
versionCmd Run function executes
      |
      v
Prints version

