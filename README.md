# ЗАПУСК

клиент:
```
cd chat-service && go run ./cmd/client/main.go
```

сервер:

1. ```docker compose up```
 
2. ```cd storage-service && make migrate-up``` чтобы накатить миграцию 

