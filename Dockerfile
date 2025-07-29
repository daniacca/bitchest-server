FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o bitchest ./cmd/server/main.go

# Final minimal stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bitchest .

EXPOSE 7463

CMD ["./bitchest", "-host", "0.0.0.0"]
