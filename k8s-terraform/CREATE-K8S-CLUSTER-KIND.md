# Create a K8S Cluster With KindD

## Pre Prerequisites

For examples docker + KinD

[KinD Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager)

The commands in MacOS would be:
```shell
brew cask install docker
brew install kind@0.11.1
```

## Create the cluster
```shell
export CLUSTER_NAME="tfm";
cat <<EOF | kind create cluster --name ${CLUSTER_NAME} --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  # the following ports can be used as NodePort -> see https://github.com/confluentinc/cp-helm-charts/blob/master/charts/cp-kafka/values.yaml#L137
  - containerPort: 31090
    hostPort: 31090
    listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
    protocol: TCP
  - containerPort: 31091
    hostPort: 31091
    listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
    protocol: TCP
  - containerPort: 31092
    hostPort: 31092
    listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
    protocol: TCP
    
EOF


# ... or just start
docker start tfm-control-plane

# check connectivity
kubectl cluster-info --context kind-tfm

docker port tfm-control-plane
# 6443/tcp -> 127.0.0.1:55313
# 31090/tcp -> 0.0.0.0:31090
# 31091/tcp -> 0.0.0.0:31091
# 31092/tcp -> 0.0.0.0:31092

kind get clusters
kubectl config use-context kind-tfm
kubectl cluster-info
```

Cluster deletion

```shell
kind delete cluster --name=tfm
# or
docker stop tfm-control-plane
docker rm tfm-control-plane

# check there's no "tfm-control-plane" container
docker container ls
```