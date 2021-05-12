FROM golang:1.16.4-alpine3.13 as base

LABEL maintainer="EEQ Team"
LABEL service="Codefresh Deploy"

RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet git
WORKDIR /

# Prune
COPY prune/* ./
RUN go build -o ./prune .
RUN chmod +x prune

# Helm #
FROM codefresh/cfstep-helm:3.2.4

ARG DEFAULT_HELM_REPO_URL

ENV CODEFRESH_VERSION=v0.75.18
ENV KUBECTL_VERSION=v1.21.0
ENV KR_HELM_REPO_URL $DEFAULT_HELM_REPO_URL

# Default packages #
RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet rhash gettext libstdc++

# Codefresh #
RUN curl -L --silent https://github.com/codefresh-io/cli/releases/download/${CODEFRESH_VERSION}/codefresh-${CODEFRESH_VERSION}-alpine-x64.tar.gz -o codefresh.tar.gz \
    && tar -zxf codefresh.tar.gz \
    && mv ./codefresh /usr/local/bin/codefresh \
    && chmod +x /usr/local/bin/codefresh

# Kubectl #
RUN curl -LO --silent https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && mv ./kubectl /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl

WORKDIR /
COPY deploy/* ./
COPY --from=base /prune ./
RUN chmod +x deploy
RUN chmod +x prune

# We need to do this because cfstep-helm has an entrypoint
ENTRYPOINT []
