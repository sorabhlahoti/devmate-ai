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
