# Build stage
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/global-bundle.pem /app/global-bundle.pem

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/api"]
