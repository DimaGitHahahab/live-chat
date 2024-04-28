FROM golang:alpine AS builder

WORKDIR /storage

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o storageserver ./cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /storage/.env .

COPY --from=builder /storage/storageserver .

CMD ["./storageserver"]