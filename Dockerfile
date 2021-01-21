FROM golang:1.14-alpine as base

LABEL maintainer="EEQ Team"
LABEL service="Codefresh Deploy"

RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet git
WORKDIR /

# Prune
COPY prune/* ./
RUN go build -o ./prune .
RUN chmod +x prune

FROM alpine:3.13.0

ARG DEFAULT_HELM_REPO_URL
ARG KUBE_CONTEXT
ARG APP_DOMAIN

ENV CODEFRESH_VERSION=v0.74.4
ENV KUBECTL_VERSION=v1.19.3
ENV HELM_VERSION=v3.4.2
ENV HELM_REPO_URL $DEFAULT_HELM_REPO_URL
ENV KUBE_CONTEXT $KUBE_CONTEXT
ENV APP_DOMAIN $APP_DOMAIN

# Default packages #
RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet rhash gettext libstdc++ curl jq

# Codefresh #
RUN curl -L --silent https://github.com/codefresh-io/cli/releases/download/${CODEFRESH_VERSION}/codefresh-${CODEFRESH_VERSION}-alpine-x64.tar.gz -o codefresh.tar.gz \
    && tar -zxf codefresh.tar.gz \
    && mv ./codefresh /usr/local/bin/codefresh \
    && chmod +x /usr/local/bin/codefresh \
    && rm -f codefresh.tar.gz

# Kubectl #
RUN curl -LO --silent https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && mv ./kubectl /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl

# Helm
RUN curl -LO --silent https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && tar -zxf helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && mv linux-amd64/helm /usr/local/bin/helm \
    && rm -f helm-${HELM_VERSION}-linux-amd64.tar.gz

WORKDIR /
COPY deploy/* ./
COPY --from=base /prune ./
RUN chmod +x deploy
RUN chmod +x prune

# We need to do this because cfstep-helm has an entrypoint
ENTRYPOINT []
