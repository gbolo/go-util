#!/bin/bash
set -e

API_URL="http://127.0.0.1:60081"
SERVICE_ID=71a8b045bfac

TITLE() {
  echo "------------------------------------------------------"
  echo " > ${1}"
  echo "------------------------------------------------------"
}

TITLE "CREATE A SERVICE"
http --check-status -v POST ${API_URL}/api/v1/service description="new service"

TITLE "CREATE ANOTHER SERVICE WITH A SPECIFIED ID"
http --check-status -v POST ${API_URL}/api/v1/service description="new service with ID ${SERVICE_ID}" id=${SERVICE_ID}

TITLE "LIST ALL SERVICES"
http --check-status ${API_URL}/api/v1/service

TITLE "GENERATE AN API KEY"
http --check-status -v POST ${API_URL}/api/v1/key name="key-1" service_id=${SERVICE_ID}

TITLE "GENERATE ANOTHER API KEY"
API_KEY=$(http POST ${API_URL}/api/v1/key name="key-2" service_id=${SERVICE_ID} | jq . -r)

TITLE "LIST ALL SERVICES"
http --check-status ${API_URL}/api/v1/service

TITLE "VALIDATE API KEY"
http --check-status -v ${API_URL}/api/v1/validate/${SERVICE_ID} "X-API-KEY: ${API_KEY}"

TITLE "REVOKE SECOND API KEY"
http --check-status -v DELETE ${API_URL}/api/v1/key prefix="${API_KEY:0:8}" service_id=${SERVICE_ID}

TITLE "CHECK THAT KEY IS NOT VALID"
http -v ${API_URL}/api/v1/validate/${SERVICE_ID} "X-API-KEY: ${API_KEY}"

TITLE "LIST ALL SERVICES"
http --check-status ${API_URL}/api/v1/service
