#!/bin/bash


export CP_SCHEMA_REGISTRY_URL="http://localhost:8081"

curl --request GET -sL \
     --url ${CP_SCHEMA_REGISTRY_URL}/subjects | jq

curl --request DELETE -sL \
     --url ${CP_SCHEMA_REGISTRY_URL}/subjects/event.customer.entity-key

curl --request DELETE -sL \
     --url ${CP_SCHEMA_REGISTRY_URL}/subjects/event.customer.entity-value

curl --request DELETE -sL \
     --url ${CP_SCHEMA_REGISTRY_URL}/subjects/topic-key

curl --request DELETE -sL \
     --url ${CP_SCHEMA_REGISTRY_URL}/subjects/topic-value