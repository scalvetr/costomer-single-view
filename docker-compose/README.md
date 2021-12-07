# Install Confluent Platform

## Prerequisites

Docker Compose

## Build

Start environment

```shell
docker-compose -p tfm \
-f docker-compose.confluent.yml \
-f docker-compose.confluent.changes.yml \
-f docker-compose.yml \
build
```
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