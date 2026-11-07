FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/coderunner ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache python3 nodejs bash

WORKDIR /app
COPY --from=builder /app/coderunner /app/coderunner

EXPOSE 8080
ENTRYPOINT ["/app/coderunner"]
