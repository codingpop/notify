# Notify

A non-blocking notification Go client library.

## Running the executable

### Install the dependencies

Run:

```sh
cd ./notification
```

And then run:

```sh
  go get ./...
```

### Using the provided listner server:
1. Start the server.
```sh
go run cmd/listener/main.go
```
2. Start the cli.
```sh
  go run cmd/cli/main.go --url http://localhost:6340 --interval 1s
```  

### Send tons of requests at once:

```sh
go run cmd/cli/main.go --url http://localhost:6340 --interval 1s < messages.txt
```

*Note: You can observe the requests in the listener console*   



### Using an external server (e.g., https://jsonplaceholder.typicode.com/posts):
```sh
  go run cmd/cli/main.go --url https://jsonplaceholder.typicode.com/posts --interval 1s
```

## Testing

```sh
go test -v -cover ./... -count=1
```

## Thank you for coming! 🍻
