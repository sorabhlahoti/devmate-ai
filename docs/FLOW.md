# Application Flow

## Ask Command Flow

```text
User runs:
devmate ask "explain goroutine"

        |
        v

cmd/devmate/main.go

        |
        v

cli.Execute()

        |
        v

Cobra finds ask command

        |
        v

ask.go joins user arguments into one question

        |
        v

AI provider is created

        |
        v

MockProvider.Ask() is called

        |
        v

Mock answer is returned

        |
        v

CLI prints answer in terminal
```
## Explain Error Flow   

```text      
User runs:
devmate explain-error --file error.log

        |
        v

Cobra finds explain-error command

        |
        v

CLI reads --file flag value

        |
        v

os.ReadFile() reads error.log

        |
        v

BuildErrorExplanationPrompt() creates AI prompt

        |
        v

MockProvider.Ask() receives prompt

        |
        v

Mock explanation is returned

        |
        v

CLI prints explanation 
```
## Save Note Flow

```text
User runs:
devmate save --title "Docker fix" --body "Check permission"

        |
        v

Cobra finds save command

        |
        v

CLI reads title and body flags

        |
        v

storage.NewSQLiteStore("") opens local SQLite DB

        |
        v

Store.Init() creates notes table if missing

        |
        v

Store.SaveNote() validates and inserts note

        |
        v

CLI prints saved note ID
```

## List Notes Flow
```text
User runs:
devmate list

        |
        v

Cobra finds list command

        |
        v

SQLite database opens

        |
        v

Store.ListNotes() fetches latest notes

        |
        v

CLI prints notes in terminal
```

## Search Notes Flow
```text
User runs:
devmate search "docker"

        |
        v

Cobra finds search command

        |
        v

CLI joins search arguments

        |
        v

Store.SearchNotes() searches title and body

        |
        v

CLI prints matching notes
```

## Why Mock AI First?

We are using mock AI first because we want to learn clean Go architecture before connecting real AI APIs.

Later we can replace:

```
provider := ai.NewMockProvider()
```

with:

```
provider := ai.NewCloudflareProvider(...)
```

or:

```
provider := ai.NewOpenAIProvider(...)
```
without changing the whole CLI.
