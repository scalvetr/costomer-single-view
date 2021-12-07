FROM golang:1.17.4-alpine3.15 AS build-env

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

# All these steps will be cached
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# https://github.com/confluentinc/confluent-kafka-go/blob/master/kafka/README.md#build-tags
# Build the binary  CGO_ENABLED=0
RUN go build -a -tags musl -o /go/bin/app

FROM alpine:3.15

RUN apk add --no-cache libc6-compat
COPY --from=build-env /go/bin/app /go/bin/app

# Run
ENV TOPIC_NAME="topic"
ENV BOOTSTRAP_SERVERS="localhost:9092"
ENV SCHEMA_REGISTRY_URL="localhost:8081"

ENTRYPOINT ["/go/bin/app"]