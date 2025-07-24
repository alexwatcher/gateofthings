FROM golang:1.24.5 AS builder

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM scratch

COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /migrations

ENTRYPOINT ["goose"]
