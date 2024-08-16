FROM golang:1.18.2-alpine3.15 as base

LABEL maintainer="EEQ Team"
LABEL service="Kube Review"

RUN apk --no-cache --quiet update

WORKDIR /

# Prune.
COPY src/prune/* ./
RUN go build -o ./prune .
RUN chmod +x prune

FROM alpine

ARG DEFAULT_HELM_REPO_URL

ENV KUBECTL_VERSION=v1.24.11
ENV KUSTOMIZE_VERSION=v5.0.1
ENV KR_BASE_OVERLAY_PATH=/usr/local/kube-review/deploy/resources/base

# Default packages #
RUN apk --no-cache --quiet update \
    && apk add --no-cache --quiet rhash gettext moreutils curl bash git jq python3 py3-pip

# AWS CLI
RUN pip install wheel \
  pip install "Cython<3.0" "pyyaml<6" --no-build-isolation \
  pip3 install awscli
  
# Kubectl #
RUN curl -LO --silent https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
  && mv ./kubectl /usr/local/bin/kubectl \
  && chmod +x /usr/local/bin/kubectl

# Kustomize
RUN curl --proto "=https" -L --silent https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_amd64.tar.gz -o kustomize.tar.gz \
  && tar -zxf kustomize.tar.gz \
  && mv ./kustomize /usr/local/bin/kustomize \
  && rm -f kustomize.tar.gz

# Cleaning
RUN rm -rf /var/cache/apk/*

WORKDIR /usr/local

RUN mkdir -p kube-review/deploy
COPY src/deploy kube-review/deploy/
RUN chmod +x kube-review/deploy/deploy
RUN ln -s /usr/local/kube-review/deploy/deploy /deploy

RUN mkdir -p kube-review/restart
COPY src/restart kube-review/restart/
RUN chmod +x kube-review/restart/restart
RUN ln -s /usr/local/kube-review/restart/restart /restart

RUN mkdir -p kube-review/prune
COPY --from=base /prune kube-review/prune/
RUN chmod +x kube-review/prune/prune
RUN ln -s /usr/local/kube-review/prune/prune /prune
