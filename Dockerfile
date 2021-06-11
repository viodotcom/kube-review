FROM golang:1.16.4-alpine as base

LABEL maintainer="EEQ Team"
LABEL service="Codefresh Deploy"

RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet git
WORKDIR /

# Prune
COPY prune/* ./
RUN go build -o ./prune .
RUN chmod +x prune

FROM alpine:3.13.5

ARG DEFAULT_HELM_REPO_URL

ENV CODEFRESH_VERSION=v0.75.18
ENV KUBECTL_VERSION=v1.20.5
ENV HELM_VERSION=v3.2.4
ENV KR_HELM_REPO_URL $DEFAULT_HELM_REPO_URL

# Default packages #
RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet rhash gettext libstdc++ curl bash git jq

# Codefresh #
RUN curl -L --silent https://github.com/codefresh-io/cli/releases/download/${CODEFRESH_VERSION}/codefresh-${CODEFRESH_VERSION}-alpine-x64.tar.gz -o codefresh.tar.gz \
    && tar -zxf codefresh.tar.gz \
    && mv ./codefresh /usr/local/bin/codefresh \
    && chmod +x /usr/local/bin/codefresh

# Kubectl #
RUN curl -LO --silent https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && mv ./kubectl /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl

# Helm #
RUN curl -L --silent https://get.helm.sh/helm-${HELM_VERSION}-linux-386.tar.gz -o helm.tar.gz \
    && tar -zxf helm.tar.gz \
    && mv linux-386/helm /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm

WORKDIR /
COPY deploy/* ./
COPY --from=base /prune ./
RUN chmod +x deploy
RUN chmod +x prune
