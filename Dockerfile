FROM golang:1.14-alpine as base

RUN apk add --no-cache git
WORKDIR /

COPY prune/* ./
RUN go build -o ./prune .
RUN chmod +x prune

FROM codefresh/cfstep-helm:3.0.3
ARG DEFAULT_HELM_REPO_URL
ENV HELM_REPO_URL $DEFAULT_HELM_REPO_URL

RUN apk update && apk add --no-cache rhash gettext libstdc++
RUN curl -L "https://github.com/codefresh-io/cli/releases/download/v0.72.1/codefresh-v0.72.1-alpine-x64.tar.gz" -o codefresh.tar.gz \
    && tar -zxvf codefresh.tar.gz \
    && mv ./codefresh /usr/local/bin/codefresh
RUN chmod +x /usr/local/bin/codefresh

RUN curl -LO --silent "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kubectl"
RUN mv kubectl /usr/local/bin/kubectl && chmod +x /usr/local/bin/kubectl

WORKDIR /
COPY deploy/* ./
COPY --from=base /prune ./
RUN chmod +x deploy
RUN chmod +x prune

# We need to do this because cfstep-helm has an entrypoint
ENTRYPOINT []
