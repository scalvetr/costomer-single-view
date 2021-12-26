# Install Confluent Platform

## Prerequisites

Docker Compose

Download confluent docker compose file
```shell
curl --silent --output docker-compose.confluent.yml \
  https://raw.githubusercontent.com/confluentinc/cp-all-in-one/7.0.1-post/cp-all-in-one/docker-compose.yml
```

## :gear: Build

Start environment

```shell
docker-compose -p tfm \
-f docker-compose.confluent.yml \
-f docker-compose.confluent.changes.yml \
-f docker-compose.yml \
build
```

## :running_man: Run

Start environment

```shell
docker-compose -p tfm \
-f docker-compose.confluent.yml \
-f docker-compose.confluent.changes.yml \
-f docker-compose.yml \
up -d
```

Control center accessible through the following URL: http://localhost:9021/

