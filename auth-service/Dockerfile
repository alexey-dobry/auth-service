FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app/service/build

COPY ./server/go.mod .
COPY ./server/go.sum .
RUN go mod download

COPY ./server/ .
RUN go build -ldflags="-s -w" -o /app/server ./cmd

FROM scratch

WORKDIR /app/migrations

COPY ./migrations .

WORKDIR /app/server

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/server /app/server

CMD ["./server"]