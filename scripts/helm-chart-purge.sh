#!/bin/bash

#####################################################################################
# Script to remove versions of the Helm Chart in the Helm Repositories on Codefresh #
#####################################################################################

# Variables #
TOKEN=""
REPO_NAME="lab"
CHART_NAME=""
CHART_VERSION=""

usage_help() {
  echo "Helm Clean Repository"
  echo ""
  echo "Usage: ./helm-chart-purge.sh --token ABC --repo-name lab --chart-name cf-review-env --chart-version 0.0.1+39a1e9b"
  echo ""
  echo "OPTIONS:"
  echo "--token         - Your personal token - Ref.: https://codefresh.io/docs/docs/integrations/codefresh-api/"
  echo "--repo-name     - Name of the Helm Repository. Default value is Lab"
  echo "--chart-name    - Name of the Helm Chart"
  echo "--chart-version - Version of the Helm Chart"
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
    --repo-name)
        REPO_NAME=$2
        ;;
    --chart-name)
        CHART_NAME=$2
        ;;
    --chart-version)
        CHART_VERSION=$2
        ;;
  esac
  shift
done

# Checking the inputs #
if [ "${TOKEN}" = "" ]; then
  echo "Oops - You should set the --token"
  echo "Ref.: https://codefresh.io/docs/docs/integrations/codefresh-api/"
  exit 1
fi

if [ "${CHART_NAME}" = "" ]; then
  echo "Oops - You should set the --chart-name"
  exit 1
fi

if [ "${CHART_VERSION}" = "" ]; then
  echo "Oops - You should set the --chart-version"
  exit 1
fi

# Running #
echo "--> Executing Helm Purge Command on Codefresh - Repo Name: ${REPO_NAME} - Chart Name: ${CHART_NAME} - Chart Version: ${CHART_VERSION}"
curl -X DELETE -H "Authorization: Bearer ${TOKEN}" "https://h.cfcr.io/api/findhotel/${REPO_NAME}/charts/${CHART_NAME}/${CHART_VERSION}"
