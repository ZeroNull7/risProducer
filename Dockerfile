FROM golang:1.17.4-alpine3.15 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ris_producer

# final stage
FROM alpine:3.15.0
COPY --from=builder /app/ris_producer /app/
EXPOSE 8123
CMD ["run"]
ENTRYPOINT ["/app/ris_producer"]