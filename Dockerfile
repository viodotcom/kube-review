FROM golang:1.14-alpine

RUN apk add --no-cache git

WORKDIR /

COPY prune/* ./

RUN go build -o ./prune .
RUN chmod +x prune
