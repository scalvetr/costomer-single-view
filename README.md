# customer-single-view

## Prerequisites

* Install common tools: [Install common tools](INSTALL-COMMON-TOOLS.md)
* K8S: [Create K8S cluster](K8S-CREATE-CLUSTER.md)
* Terraform: [Install terraform](INSTALL-TERRAFORM.md)

## Setup initial cluster

Configure terraform
```shell
# See: https://learn.hashicorp.com/tutorials/terraform/kubernetes-provider
k8s_name="kind-data-lake"
k8s_host=`kubectl config view -o json | jq -r --arg clusterName "${k8s_name}" '.clusters[] | select(.name == $clusterName) | .cluster.server'`
k8s_cluster_ca_certificate=`kubectl config view --flatten -o json | jq -r --arg clusterName "${k8s_name}" '.clusters[] | select(.name == $clusterName) | .cluster["certificate-authority-data"]'`
k8s_client_certificate=`kubectl config view --flatten -o json | jq -r --arg userName "${k8s_name}" '.users[] | select(.name == $userName) | .user["client-certificate-data"]'`
k8s_client_key=`kubectl config view --flatten -o json | jq -r --arg userName "${k8s_name}" '.users[] | select(.name == $userName) | .user["client-key-data"]'`

cat > terraform.tfvars << EOF
# variables.tf

k8s_host                   = "${k8s_host}"
k8s_client_certificate     = "${k8s_client_certificate}"
k8s_client_key             = "${k8s_client_key}"
k8s_cluster_ca_certificate = "${k8s_cluster_ca_certificate}"
k8s_project_label          = "customer-single-view"
k8s_namespace              = "customer-single-view"

EOF

```

terraform apply -var-file terraform.tfvars