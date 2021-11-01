terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

variable "k8s_host" {
  type = string
}

variable "k8s_client_certificate" {
  type = string
}

variable "k8s_client_key" {
  type = string
}

variable "k8s_cluster_ca_certificate" {
  type = string
}

provider "kubernetes" {
  host = var.k8s_host

  client_certificate     = base64decode(var.k8s_client_certificate)
  client_key             = base64decode(var.k8s_client_key)
  cluster_ca_certificate = base64decode(var.k8s_cluster_ca_certificate)
}