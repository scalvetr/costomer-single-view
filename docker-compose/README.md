# Install Confluent Platform

## Prerequisites

Docker Compose

## Run

Start environment

```shell
docker-compose -p tfm \
-f docker-compose.confluent.yml \
-f docker-compose.confluent.changes.yml \
-f docker-compose.yml \
up -d
```

Services exposed