FROM golang:1.14 as builder

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN apt update

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/chesscli ./cmd/chesscli
RUN go build -o bin/game-server ./cmd/game-server
RUN go build -o bin/user-stats ./cmd/user-stats

FROM alpine:3.12

COPY --from=builder /app/bin/chesscli /app/
COPY --from=builder /app/bin/game-server /app/
COPY --from=builder /app/bin/user-stats /app/

WORKDIR /app
