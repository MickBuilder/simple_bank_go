# Build stage
FROM golang:1.22-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM golang:1.22-alpine3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY db/migration db/migration

EXPOSE 8080
CMD [ "/app/main" ]