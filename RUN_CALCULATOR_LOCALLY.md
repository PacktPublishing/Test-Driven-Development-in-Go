# Running the Calculator CLI application locally
These instructions are required from `chapter01` to `chapter03`.

## Dependencies
You will need the following tools to run your the code in this repository: 
- [Go 1.19 or later](https://go.dev/doc/install)

This project depends on a variety of testing tools, which are used for instruction and demonstration purposes:
- [testify](https://github.com/stretchr/testify)
- [mockery](https://github.com/vektra/mockery)

## Run locally
The Calculator CLI application requires an expression to parse: 
```bash
$ go run chapterXX/main.go -expression "2 + 3"
```

## Run unit tests
1. Install all [Dependencies](#dependencies). 
1. Run all tests in a given chapter:
```
$ cd chapterXX 
$ go test ./... -v
```
