FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o service1-go .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/service1-go .

# Create a non-root user with UID in range 10000-20000 (required by Choreo)
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10014 \
    "choreo"
USER 10014

ENTRYPOINT ["./service1-go"]
