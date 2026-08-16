# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/server

# Final Runtime Stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/server /app/server

EXPOSE 8080

ENV PORT=8080
ENV GIN_MODE=release

ENTRYPOINT ["/app/server"]
