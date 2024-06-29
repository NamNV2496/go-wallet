# build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY /internal/databaseaccess/migration/sql ./internal/databaseaccess/migration/sql

EXPOSE 8080
CMD [ "/app/main" ]
