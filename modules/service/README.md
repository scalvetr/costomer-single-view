# Service

## :gear: Build
Build go program
```shell
go build
```

Build docker image
```shell
docker build . -t service:latest

# Test

docker run --rm --name=datagen datagen:latest
```

## :rocket: Run
```shell
export PORT="8080"
export CUSTOMER_DB_URI="mongodb://localhost:27018"
export CUSTOMER_DB_NAME="customer-single-view"
export CUSTOMER_DB_USER="user"
export CUSTOMER_DB_PASSWORD="password"

go run -a *.go
```

## :white_check_mark: Test
```shell
curl -vvv http://localhost:8080/api/customers | jq

```
