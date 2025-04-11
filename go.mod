module flagsmith-test

go 1.22

toolchain go1.24.2

require (
	github.com/Flagsmith/flagsmith-go-client/v2 v2.3.1
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/Flagsmith/flagsmith-go-client/v4 v4.2.0 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/go-resty/resty/v2 v2.14.0 // indirect
	github.com/itlightning/dateparse v0.2.0 // indirect
	golang.org/x/exp v0.0.0-20230713183714-613f0c0eb8a1 // indirect
	golang.org/x/net v0.27.0 // indirect
)

replace github.com/Flagsmith/flagsmith-go-client/v2 => ../flagsmith-go-client
