## Live chat system

### Chat client
A console application that connects to the server via WebSocket. It shows the most recent messages and allows to chat with other clients.

### Chat service
Manages broadcast of messages between clients, produces messages to Kafka, retrieves message history from Storage service with http.
### Storage microservice
A separate app that receives user messages from the main service via Kafka and stores them in the PostgreSQL database. It uses Redis to perform message history caching.


## Running locally

### Client:
```cd chat-service && go run ./cmd/client/main.go```

### Server:
Start the Docker services:

```docker compose up```

To apply database migration:

```cd storage-service && make migrate-up```

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

