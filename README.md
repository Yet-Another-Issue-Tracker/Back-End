# Overview
TODO

## Running the server
To run the server, follow these simple steps:

```
go run main.go
```

## Run test

```
go test ./routes
```

## Test the API performance

### Projects 
```
npx autocannon -c 1 -m POST -b '{ "client": "pippo", "type": "scrum" }' http://localhost:8080/v1/projects
```