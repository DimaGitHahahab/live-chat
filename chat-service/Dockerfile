FROM golang:alpine AS builder

WORKDIR /chat

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o chatserver ./cmd/server/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /chat/.env .

COPY --from=builder /chat/chatserver .

CMD ["./chatserver"]