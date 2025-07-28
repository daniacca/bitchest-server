FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bitchest ./cmd/server/main.go

# Stage finale minimale
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bitchest .

EXPOSE 7463

CMD ["./bitchest"]
