# Create a K8S Cluster

## Pre Prerequisites

For examples docker + KinD


The commands in MacOS would be:
```shell
brew cask install docker
brew install kind@0.11.1
```

## Create the cluster
```shell
# install ...
kind create cluster --name data-lake
# ... or just start
docker start data-lake-control-plane

# check connectivity
kubectl cluster-info --context kind-data-lake

docker port data-lake-control-plane
# 6443/tcp -> 127.0.0.1:49507

kind get clusters
kubectl config use-context kind-data-lake
kubectl cluster-info

```