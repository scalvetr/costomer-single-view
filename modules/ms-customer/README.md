# MS Customer

## Build

```shell
go build
```

## Run
```shell
export WORK_DIR="../../";

# k8s node-port
# export BOOTSTRAP_SERVERS="192.168.64.3:9092"
export BOOTSTRAP_SERVERS="localhost:9092"

# kubectl port-forward
export SCHEMA_REGISTRY_URL="http://localhost:8081"
export TOPIC_NAME="event.customer.entity"

go run main.go \
--key-schema-file ${WORK_DIR}schemas/customer-key.avro \
--value-schema-file ${WORK_DIR}schemas/customer-value.avro \
--data-file ${WORK_DIR}test-data/customer-data.json
```

See:
* https://github.com/confluentinc/confluent-kafka-go
* https://docs.confluent.io/platform/current/clients/confluent-kafka-go/index.html
* 