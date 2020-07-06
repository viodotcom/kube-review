FROM golang:1.14-alpine

RUN apk add --no-cache git

WORKDIR /app/

COPY . .

RUN go build -o ./cf-stop-stale-envs .
