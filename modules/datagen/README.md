# MS Customer

## Build
Build go program
```shell
go build
```

Build docker image
```shell
docker build . -t datagen:latest

# Test

docker run --rm --name=datagen datagen:latest
```

## Libs
```

module datagen

go 1.17

require (
github.com/confluentinc/confluent-kafka-go v1.7.0
github.com/jaswdr/faker v1.8.0
github.com/mycujoo/go-kafka-avro/v2 v2.0.0
)
```

```shell

go mod tidy
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

go run *.go \
--key-schema-file ${WORK_DIR}schemas/customer-key.avsc \
--value-schema-file ${WORK_DIR}schemas/customer-value.avsc
```

See:
* https://github.com/confluentinc/confluent-kafka-go
* https://docs.confluent.io/platform/current/clients/confluent-kafka-go/index.html