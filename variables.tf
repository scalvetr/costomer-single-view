# Variable declaration

variable "k8s_host" {
  type = string
}

variable "k8s_client_certificate" {
  type = string
  sensitive = true
}

variable "k8s_client_key" {
  type = string
  sensitive = true
}

variable "k8s_cluster_ca_certificate" {
  type = string
  sensitive = true
}


variable "k8s_project_label" {
  type = string
}

variable "k8s_namespace" {
  type = string
}