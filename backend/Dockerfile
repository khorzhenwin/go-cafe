# Build stage
FROM golang:1.25.7 AS builder
WORKDIR /app
ARG GENERATE_SWAGGER=true

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN if [ "$GENERATE_SWAGGER" = "true" ]; then /go/bin/swag init -g main.go -d ./cmd/api,./internal -o ./docs --parseInternal; fi
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/global-bundle.pem /app/global-bundle.pem

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/api"]
