FROM golang:1.14-alpine

RUN apk add --no-cache git

WORKDIR /

COPY . .

RUN go build -o ./cf-stop-stale-envs .
RUN chmod +x cf-stop-stale-envs
