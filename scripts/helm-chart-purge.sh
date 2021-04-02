#!/bin/bash

#####################################################################################
# Script to remove versions of the Helm Chart in the Helm Repositories on Codefresh #
#####################################################################################

# Variables #
TOKEN=""
CODEFRESH_HELM_REPO="CF_HELM_LAB"
REPO_NAME="lab"
CHART_NAME=""
CHART_VERSION=""
ALL=false

usage_help() {
  echo "Helm Clean Repository"
  echo ""
  echo "Removing all Chart Versions: ./helm-chart-purge.sh --token ABC --codefresh-helm-repo CF_HELM_LAB --repo-name lab --chart-name kube-review --all true"
  echo ""
  echo "Removing a specific Chart Version: ./helm-chart-purge.sh --token ABC --repo-name lab --chart-name kube-review --chart-version 0.0.1+39a1e9b --all false"
  echo ""
  echo "OPTIONS:"
  echo "--token               - Your personal token - Ref.: https://codefresh.io/docs/docs/integrations/codefresh-api/"
  echo "--codefresh-helm-repo - Name of the Helm Chart Repository, CF_HELM_LAB or CF_HELM_DEFAULT"
  echo "--repo-name           - Name of the Helm Repository, lab or default"
  echo "--chart-name          - Name of the Helm Chart"
  echo "--chart-version       - Version of the Helm Chart"
  echo "--all                 - Remove all versions, default value is false"
  echo "--help"
}

while [ $# -gt 0 ]; do
  case "$1" in
    --help)
        usage_help
        exit 1
        ;;
    --token)
        TOKEN=$2
        ;;
    --codefresh-helm-repo)
        CODEFRESH_HELM_REPO=$2
        ;;
    --repo-name)
        REPO_NAME=$2
        ;;
    --chart-name)
        CHART_NAME=$2
        ;;
    --chart-version)
        CHART_VERSION=$2
        ;;
    --all)
        ALL=$2
        ;;
  esac
  shift
done

# Checking the inputs #
if [ "${TOKEN}" = "" ]; then
  echo "Oops - You should set the --token parameter"
  echo "Ref.: https://codefresh.io/docs/docs/integrations/codefresh-api/"
  exit 1
fi

if [ "${CHART_NAME}" = "" ]; then
  echo "Oops - You should set the --chart-name parameter"
  exit 1
fi

if [ "${ALL}" != true ]; then

  if [ "${CHART_VERSION}" = "" ]; then
    echo "Oops - Remove a specific chart version, you should set the --chart-version parameter"
    exit 1
  fi

  echo "--> Executing Helm Purge Command on Codefresh - Repo Name: ${REPO_NAME} - Chart Name: ${CHART_NAME} - Chart Version: ${CHART_VERSION}"
  curl -X DELETE -H "Authorization: Bearer ${TOKEN}" "https://h.cfcr.io/api/findhotel/${REPO_NAME}/charts/${CHART_NAME}/${CHART_VERSION}"

else

  if [ "${CHART_VERSION}" != "" ]; then
    echo "Oops - Remove all chart versions, you should remove the --chart-version parameter"
    exit 1
  fi

  for CHART_VERSIONS in $(curl --silent -X GET -H "Authorization: Bearer "${TOKEN} -H "Content-Type: application/json; charset=utf-8" "https://g.codefresh.io/api/charts/${CODEFRESH_HELM_REPO}/${CHART_NAME}?full=true" | jq '.versions[].version' | tr -d '"'); \
    do \
      echo "--> Executing Helm Purge Command on Codefresh - Repo Name: ${REPO_NAME} - Chart Name: ${CHART_NAME} - Chart Versions: ${CHART_VERSIONS}" ; \
      curl -X DELETE -H "Authorization: Bearer ${TOKEN}" "https://h.cfcr.io/api/findhotel/${REPO_NAME}/charts/${CHART_NAME}/${CHART_VERSION}" ; \
    done

fi
