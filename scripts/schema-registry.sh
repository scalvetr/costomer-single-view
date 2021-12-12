#!/bin/bash

export CP_SCHEMA_REGISTRY_URL="http://localhost:8081"

function delete_schemas() {
    # shellcheck disable=SC2006
    subjects=`curl --request GET --url ${CP_SCHEMA_REGISTRY_URL}/subjects`
    for subject in $(echo "${subjects}" | jq -r '.[]'); do
        curl --request DELETE -sL \
             --url "${CP_SCHEMA_REGISTRY_URL}/subjects/${subject}"
    done
}


delete_schemas;
#curl --request DELETE -sL \
#     --url ${CP_SCHEMA_REGISTRY_URL}/subjects/event_customer_entity-key

#curl --request DELETE -sL \
#     --url ${CP_SCHEMA_REGISTRY_URL}/subjects/event_customer_entity-value
