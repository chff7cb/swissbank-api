# SwissBank API

[![codecov](https://codecov.io/gh/chff7cb/swissbank-api/graph/badge.svg?token=843EQJ5CJ7)](https://codecov.io/gh/chff7cb/swissbank-api)

Simple application written in Golang that provides an API for managing financial accounts and transactions.

## About
SwissBank is built mainly using:

- [Fx](https://github.com/uber-go/fx) for handling applicaton lifecycle.
- [Gin](https://github.com/gin-gonic/gin) for HTTP request handling.

Unit tests leverage:

- [testify](https://github.com/stretchr/testify) for test suite and mocking capabilities.
- [mockery](https://github.com/vektra/mockery) for mock code generation and integration.

The application also uses [DynamoDB](https://aws.amazon.com/pm/dynamodb/) as primary database.

## Getting started

### Running locally

The application itself may be run locally from the `src/` directory as follows:
```bash
go run cmd/http/main.go
```

or built and ran like so:
```bash
go build cmd/http/main.go
```
```bash
./main
```

A `Dockerfile` is also available for building and running the application as a Docker container.
```bash
docker build -t swissbank-api .
```
```bash
docker run -p 8182:8182 -e AWS_ACCESS_KEY_ID=DUMMYIDEXAMPLE -e AWS_SECRET_ACCESS_KEY=DUMMYEXAMPLEKEY -e SWISSBANK_AWS_REGION=us-east-1 -e SWISSBANK_DYNAMODB_ENDPOINT_URL=http://172.17.0.1:8000 swissbank-api
```

Now you should be able to navigate to:
```
http://localhost:8182/swagger/index.html
```

and see the generated swagger documentation.

**NOTE: In order for the application to work properly we will have to provide it with a functional DynamoDB instance.**
This can be achieved by either connecting with an existing DynamoDB instance on AWS cloud or using [DynamoDB local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.DownloadingAndRunning.html).

### Docker Compose

Finally there's also a `docker-compose` file that wraps up the construction of a container image using the available Dockerfile and running a DynamoDB local container as well as the required configuration parameters.

## Database setup

### Using DynamoDB local

[DynamoDB local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.DownloadingAndRunning.html) is available as a Docker image that we can use for running it locally:
```base
docker run -p 8000:8000 amazon/dynamodb-local:latest
```
For the application to access  DynamoDB local running on Docker we have to configure the environment variables `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `SWISSBANK_AWS_REGION` **using dummy/passthrough values** which are required by the AWS SDK.

**See previous example of running with Docker.**

The source project also contains a [`local.env`](local.env) file outlying these required parameters.

### DynamoDB on AWS cloud

Firstly, we will have to manually create the required tables. Checkout [src/providers/migrations.go](src/providers/migrations.go) to see what the expected table schema should look like.

To use an AWS DynamoDB instance we may run the application in any of the ways mentioned before 
but the following environment parameters will have to be provided:
```base
SWISSBANK_AWS_PROFILE=<ACTUAL_AWS_PROFILE>
SWISSBANK_AWS_REGION=<REAL_AWS_REGION>
```

these are available in any environment after running the CLI with `aws configure` which will make the SDK read the credentials from `~/.aws/config` and/or `~/.aws/credentials`.

## Code/Contributing

### Tests

All the application code resides in the `src/` folder.

You can run and report tests coverage with the following commands:
```bash
go test ./... -coverprofile /tmp/cover.out
```
```bash
go tool cover -func /tmp/cover.out
```

**Right now only code under the `src/cmd`, `src/mocks` and `src/docs` directories are untested.**

### Checking

Code is being checked with the `golangci-lint` tool:

```bash
golangci-lint run --enable-all -D depguard -D paralleltest -D tagliatelle -D godot -D gofumpt -D exhaustivestruct -D lll -D exhaustruct -D nonamedreturns -D ireturn -D wrapcheck -D nlreturn
```

The last version of the tool used was `v1.54.2` including all available linters except the ones excluded above.

## License

Unlincesed