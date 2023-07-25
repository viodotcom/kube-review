#!/bin/bash
set -e

if [[ -z "${CHECKOUT_STRIPE_RESTRICTED_API_KEY}" || -z "${CHECKOUT_URL}" ]]; then
  echo "One or more required variables are missing: CHECKOUT_STRIPE_RESTRICTED_API_KEY or CHECKOUT_URL"
  exit 1
fi

STRIPE_ENDPOINT="https://api.stripe.com/v1/apple_pay/domains"
# extract domain name from url
CHECKOUT_DOMAIN_NAME=$(echo $CHECKOUT_URL | awk -F[/:] '{print $4}')

# register domain name with apple pay
echo "Registering ${CHECKOUT_DOMAIN_NAME} with Apple Pay"

curl "${STRIPE_ENDPOINT}" \
 -u "${CHECKOUT_STRIPE_RESTRICTED_API_KEY}": \
 -d domain_name="${CHECKOUT_DOMAIN_NAME}"
