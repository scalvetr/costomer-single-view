# Create a K8S Cluster With Minikube

## Pre Prerequisites


The commands in MacOS would be:
```shell
brew cask install docker
brew install hyperkit@0.20200908
brew install minikube@1.24.0
```

## Create the cluster
```shell
export CLUSTER_NAME="data-lake";
minikube start --driver=hyperkit --profile=${CLUSTER_NAME} --kubernetes-version=v1.22.3
# --driver=docker
minikube profile list


# check connectivity
kubectl cluster-info --context data-lake
minikube ip --profile=${CLUSTER_NAME}

docker port data-lake
# 22/tcp -> 127.0.0.1:64553
# 2376/tcp -> 127.0.0.1:64549
# 32443/tcp -> 127.0.0.1:64550
# 5000/tcp -> 127.0.0.1:64551
# 8443/tcp -> 127.0.0.1:64552


kubectl config use-context data-lake
kubectl cluster-info
```

Cluster deletion

```shell
minikube delete --profile=${CLUSTER_NAME}
# or
docker stop data-lake
docker rm data-lake

# check there's no "data-lake-control-plane" container
docker container ls
```