# Running the BookSwap application locally

## Dependencies
You will need the following tools to run your the code in this repository: 
- [Go 1.19 or later](https://go.dev/doc/install)
- [PostgresSQL 15 or later](https://www.postgresql.org/download/)
- [Docker 4.17 or later](https://www.docker.com/products/docker-desktop/)

This project depends on a variety of testing tools, which are used for instruction and demonstration purposes:
- [testify](https://github.com/stretchr/testify)
- [mockery](https://github.com/vektra/mockery)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [ginkgo](https://github.com/onsi/ginkgo)
- [godog](https://github.com/cucumber/godog)
- [pact](https://github.com/pact-foundation)

## Run locally
From `Chapter06` onwards, you will require a database to run the `BookSwap` application:
1. Install [PostgresSQL](https://www.postgresql.org/download/) according to the installation steps for your  operating system.
1. Export the following user variables: 
```
$ export BOOKSWAP_DB_URL=XXX
$ export BOOKSWAP_PORT=XXX
$ export BOOKSWAP_BASE_URL=XXX
```
3. Run the `BookSwap` executable using the `go run chapterXX/cmd/main.go` command. The application will then listen on the configured port.

## Run in Docker 
From `chapter06` onwards, you can run the `BookSwap` application with Docker: 
1. Install [Docker](https://docs.docker.com/get-docker/) according to the installation steps for your operating system. Separate Docker configuration files have been provided for each chapter. For example, `docker-compose.book-swap.chapter06.yml` runs the version of the application corresponding to the `chapter06` directory.
1. Create a `docker.env` file will all the variables required for running the application. Example values can be seen below: 
```
POSTGRES_USER=root
POSTGRES_PASSWORD=root
POSTGRES_DB=books
BOOKSWAP_DB_URL=postgres://root:root@db:5432/books?sslmode=disable
BOOKSWAP_PORT=3000
```
3. Run the application using `docker compose`: 
```
$  docker compose -f docker-compose.book-swap.chapterXX.yml up --build
$  docker compose -f docker-compose.book-swap.chapterXX.yml down 
```
4. The application will then listen on the configured port. 

```
curl --location --request GET 'http://localhost:3000/
```

## Generate mocks
```
$ mockery --dir "chapterXX" --output "chapterXX/mocks" --all
```

## Run unit tests
1. Install all [Dependencies](#dependencies) and follow the [Run locally](#run-locally) steps. 
1. Run all tests:
```
$ go test ./... -v
```

## Run all tests
**WARNING:** as tests often require a specific version of the test endpoint to be running, I do not recommend running all tests from the root directory.

Instead, I recommend running unit & integration tests per chapter:
1. Run the Docker version of the `chapterXX` according to [Run in Docker](#run-in-docker).
1. In a separate window, change to the desired directory `cd chapterXX`.
1. Run all tests in the `chapterXX` directory by setting to `LONG` commandline argument:
```
$ LONG=true go test ./... -v
```
 
## Postman collection
For your convenience, a [Postman](https://www.postman.com/downloads/) collection with requests for the BookSwap application has been provided. See `BookSwap.postman_collection.json`. This file can then be used to [import the collection into Postman](https://learning.postman.com/docs/getting-started/importing-and-exporting-data/#importing-data-into-postman).