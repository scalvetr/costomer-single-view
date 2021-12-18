FROM golang:1.17.4-alpine3.15 AS build-env

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

# All these steps will be cached
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/go.mod config/
COPY config/go.sum config/
COPY controllers/go.mod controllers/
COPY controllers/go.sum controllers/
COPY models/go.mod models/
COPY models/go.sum models/
COPY routes/go.mod routes/
COPY routes/go.sum routes/

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# https://github.com/confluentinc/confluent-kafka-go/blob/master/kafka/README.md#build-tags
# Build the binary  CGO_ENABLED=0
RUN go build -a -tags musl -o /go/bin/app

FROM alpine:3.15

RUN apk add --no-cache libc6-compat bash
COPY --from=build-env /go/bin/app /go/bin/app

# Run
ENV PORT="8080"
ENV CUSTOMER_DB_URI="mongodb://localhost:27018"
ENV CUSTOMER_DB_NAME="customer-single-view"
ENV CUSTOMER_DB_USER="user"
ENV CUSTOMER_DB_PASSWORD="password"

ENTRYPOINT ["/go/bin/app"]