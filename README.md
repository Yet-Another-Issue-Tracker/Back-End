# Overview
TODO

## Running the server
To run the server, follow these simple steps:

```
go run main.go
```

## Run test

Until a fix will be released tests must be run withouth parallelism otherwise database error will be thrown caused be concurrency of tests
```
go test -p 1 ./... 
```

## Test the API performance

### Projects 
```
npx autocannon -c 1 -m POST -b '{ "client": "pippo", "type": "scrum" }' http://localhost:8080/v1/projects
```