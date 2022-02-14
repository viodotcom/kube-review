FROM golang:1.16.4-alpine as base

LABEL maintainer="EEQ Team"
LABEL service="Codefresh Deploy"

RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet git
WORKDIR /

# Prune
COPY src/prune/* ./
RUN go build -o ./prune .
RUN chmod +x prune

FROM alpine:3.13.5

ARG DEFAULT_HELM_REPO_URL

ENV KUBECTL_VERSION=v1.20.5
ENV KUSTOMIZE_VERSION=v4.5.2
ENV KR_BASE_OVERLAY_PATH=/usr/local/kube-review/deploy/resources/base

# Default packages #
RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet rhash gettext moreutils libstdc++ curl bash git jq

# Kubectl #
RUN curl -LO --silent https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && mv ./kubectl /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl

# Kustomize
RUN curl -L --silent https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_amd64.tar.gz -o kustomize.tar.gz \
    && tar -zxf kustomize.tar.gz \
    && mv ./kustomize /usr/local/bin/kustomize

WORKDIR /usr/local

RUN mkdir -p kube-review/deploy
COPY src/deploy kube-review/deploy/
RUN chmod +x kube-review/deploy/deploy
RUN ln -s /usr/local/kube-review/deploy/deploy /deploy

RUN mkdir -p kube-review/prune
COPY --from=base /prune kube-review/prune/
RUN chmod +x kube-review/prune/prune
RUN ln -s /usr/local/kube-review/prune/prune /prune
