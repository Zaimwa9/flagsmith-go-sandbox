## Flagsmith Go Sandbox

This is a sandbox for the Flagsmith Go client. It is not meant to be used as a library.

## Setup

- Install dependencies: `go mod tidy`
- Upgrade dependencies: `go get -u`
- Create a `.env` file based on the `.env.example` file

## Running

```bash
go run main.go

curl http://localhost:8082
```

To run a local version of the flagsmith go client, uncomment the following line in `main.go`:

```
// replace github.com/Flagsmith/flagsmith-go-client/v2 => {path-to-flagsmith-go-client}/flagsmith-go-client
```
