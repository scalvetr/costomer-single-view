# Account Bookings - Kafka Streams

Kafka streams applications that creates a new topic producing Accounts with bookings.

## Build native binary

```shell
./mvnw -Pnative -DskipTests package
```

## Build native container
Option 1:
```shell
./mvnw spring-boot:build-image
```
Option 2:
```shell
docker build . -f Dockerfile.native -t account-bookings-kafka-streams:latest
```
Run 

```shell
docker run --rm account-bookings-kafka-streams:latest
```

The command used to initialize the applications

```shell
curl https://start.spring.io/#!type=maven-project&language=java&platformVersion=2.6.1&packaging=jar&jvmVersion=17&groupId=edu.uoc.scalvetr&artifactId=account-bookings-kafka-streams&name=account-bookings-kafka-streams&description=Account%20Bookings%20KafkaStreams&packageName=edu.uoc.scalvetr.accountbookings&dependencies=native,kafka,kafka-streams,lombok,configuration-processor
```

