# Install Confluent Platform

## Prerequisites

* Install common tools: [Install common tools](INSTALL-COMMON-TOOLS.md)
* K8S: [Create K8S cluster](K8S-CREATE-CLUSTER.md)
* Terraform: [Install terraform](INSTALL-TERRAFORM.md)

## Setup initial cluster

Configure terraform
```shell
# See: https://learn.hashicorp.com/tutorials/terraform/kubernetes-provider
export k8s_name="kind-data-lake"
export k8s_host=`kubectl config view -o json | jq -r --arg clusterName "${k8s_name}" '.clusters[] | select(.name == $clusterName) | .cluster.server'`
export k8s_cluster_ca_certificate=`kubectl config view --flatten -o json | jq -r --arg clusterName "${k8s_name}" '.clusters[] | select(.name == $clusterName) | .cluster["certificate-authority-data"]'`
export k8s_client_certificate=`kubectl config view --flatten -o json | jq -r --arg userName "${k8s_name}" '.users[] | select(.name == $userName) | .user["client-certificate-data"]'`
export k8s_client_key=`kubectl config view --flatten -o json | jq -r --arg userName "${k8s_name}" '.users[] | select(.name == $userName) | .user["client-key-data"]'`
export k8s_project_label="customer-single-view"
export k8s_namespace="customer-single-view"

cat > terraform.tfvars << EOF
# variables.tf

k8s_host                   = "${k8s_host}"
k8s_client_certificate     = "${k8s_client_certificate}"
k8s_client_key             = "${k8s_client_key}"
k8s_cluster_ca_certificate = "${k8s_cluster_ca_certificate}"
k8s_project_label          = "${k8s_project_label}"
k8s_namespace              = "${k8s_namespace}"

EOF
```

Run terraform
```shell
terraform apply -var-file terraform.tfvars
```

Expose ports
```shell
export k8s_namespace="customer-single-view"
# control center
export CONTROL_CENTER_POD_NAME=$(kubectl -n $k8s_namespace get pods -l "app=cp-control-center" -o jsonpath="{.items[0].metadata.name}")
echo $CONTROL_CENTER_POD_NAME
kubectl -n $k8s_namespace port-forward $CONTROL_CENTER_POD_NAME 9021:cc-http
```

Access to control center:

![Control Center Homepage](img/control-center-homepage.png)


## See
* https://confluentinc.github.io/cp-helm-charts/