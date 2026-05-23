# Application Flow

## Current Architecture

```text
Developer
   |
   | runs terminal command
   v
DevMate CLI
   |
   | saves/searches
   v
Local SQLite Database
```

```text
API Client / curl / future CLI sync
   |
   | HTTP JSON request
   v
DevMate API Server
   |
   | calls storage layer
   v
SQLite Database
```

## Config flow

```textUser runs:
devmate config set api_url http://localhost:8080

        |
        v

Cobra finds config set command

        |
        v

internal/config loads existing config

        |
        v

updates api_url value

        |
        v

writes config.json into user config directory
```     

## Config Get Flow

```text
User runs:
devmate config get api_url

        |
        v

Cobra finds config get command

        |
        v

internal/config loads config.json

        |
        v

returns api_url

        |
        v

CLI prints value
```     


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

## API Create Note Flow

```text
Client sends:
POST /api/v1/notes

        |
        v

cmd/api/main.go starts server

        |
        v

internal/api/server.go receives request

        |
        v

handleNotes() checks method

        |
        v

handleCreateNote() decodes JSON body

        |
        v

storage.SaveNote() validates and inserts data

        |
        v

API returns JSON response
```

## API List Notes Flow

```text
Client sends:
GET /api/v1/notes

        |
        v

handleNotes() checks method

        |
        v

handleListNotes() reads limit

        |
        v

storage.ListNotes() fetches notes

        |
        v

API returns notes JSON

## API Search Notes Flow
```text
Client sends:
GET /api/v1/notes/search?q=docker

        |
        v

handleSearchNotes() reads q parameter

        |
        v

storage.SearchNotes() searches title/body

        |
        v

API returns matching notes      
```

## Sync Flow

```text 
User runs:
devmate sync

        |
        v

CLI loads config.json

        |
        v

CLI reads api_url

        |
        v

CLI opens local SQLite database

        |
        v

ListUnsyncedNotes() finds notes where remote_id is NULL

        |
        v

For each unsynced note:
POST /api/v1/notes

        |
        v

Backend API saves note and returns remote ID

        |
        v

CLI marks local note as synced

        |
        v

Local note now has remote_id and synced_at
```

## Swagger Flow

Flow:

1. swag reads annotations from cmd/api/main.go.
2. swag reads handler annotations from internal/api/server.go.
3. swag generates files in internal/swaggerdocs.
4. API imports generated Swagger docs.
5. http-swagger serves Swagger UI at /swagger/index.html.
6. Browser loads Swagger JSON from /swagger/doc.json.

## Docker Flow

```text
User runs:
docker compose up --build

        |
        v

Docker reads Dockerfile

        |
        v

Builder stage compiles Go API binary

        |
        v

Runtime stage copies only compiled binary

        |
        v

Container starts /app/devmate-api

        |
        v

Go API listens on port 8080

        |
        v

SQLite database is stored at /data/devmate.db

        |
        v

Docker volume keeps /data persistent
```

## Sync With Dockerized API Flow

```text
User runs:
devmate sync

        |
        v

CLI reads api_url from config

        |
        v

api_url points to http://localhost:8080

        |
        v

CLI sends POST /api/v1/notes

        |
        v

Dockerized API receives request

        |
        v

API stores note in SQLite volume

        |
        v

CLI marks local note as synced
```

## Why Docker Volume Exists

Container filesystem is temporary.

If SQLite was stored only inside container:
- container removed
- notes gone

With Docker volume:
- container removed
- volume remains
- notes survive

