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
export CLUSTER_NAME="tfm";
minikube start --driver=hyperkit -p=${CLUSTER_NAME} --kubernetes-version=v1.22.3
# --driver=docker
# next time: minikube start --profile=tfm



minikube profile list
# delete
minikube profile list
minikube delete -p=${CLUSTER_NAME} 


# check connectivity
kubectl cluster-info --context tfm
minikube ip --profile=${CLUSTER_NAME}

docker port tfm
# 22/tcp -> 127.0.0.1:64553
# 2376/tcp -> 127.0.0.1:64549
# 32443/tcp -> 127.0.0.1:64550
# 5000/tcp -> 127.0.0.1:64551
# 8443/tcp -> 127.0.0.1:64552


kubectl config use-context tfm
kubectl cluster-info
```

Cluster deletion

```shell
minikube delete --profile=${CLUSTER_NAME}
# or
docker stop tfm
docker rm tfm

# check there's no "tfm-control-plane" container
docker container ls
```