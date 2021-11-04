terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.6.1"
    }
    helm = {
      source = "hashicorp/helm"
      version = "2.3.0"
    }
  }
}

provider "kubernetes" {
  host = var.k8s_host

  client_certificate     = base64decode(var.k8s_client_certificate)
  client_key             = base64decode(var.k8s_client_key)
  cluster_ca_certificate = base64decode(var.k8s_cluster_ca_certificate)
}

provider "helm" {
  kubernetes {
    host = var.k8s_host

    client_certificate     = base64decode(var.k8s_client_certificate)
    client_key             = base64decode(var.k8s_client_key)
    cluster_ca_certificate = base64decode(var.k8s_cluster_ca_certificate)
  }
}

resource "kubernetes_namespace" "project-namespace" {
  metadata {
    annotations = {
      name = var.k8s_namespace
    }

    labels = {
      project = var.k8s_project_label
    }

    name = var.k8s_namespace
  }
}

resource "helm_release" "confluent" {
  name       = "confluent"

  repository = "https://confluentinc.github.io/cp-helm-charts"
  chart      = "cp-helm-charts"
  namespace = var.k8s_namespace

  set {
    name  = "cp-zookeeper.enabled"
    value = "true"
  }
  set {
    name  = "cp-zookeeper.servers"
    value = "1"
  }
  set {
    name  = "cp-kafka.enabled"
    value = "true"
  }
  set {
    name  = "cp-kafka.brokers"
    value = "1"
  }
  set {
    name  = "cp-schema-registry.enabled"
    value = "false"
  }
  set {
    name  = "cp-kafka-rest.enabled"
    value = "false"
  }
  set {
    name  = "cp-kafka-connect.enabled"
    value = "false"
  }
  set {
    name  = "cp-ksql-server.enabled"
    value = "false"
  }
  set {
    name  = "cp-control-center.enabled"
    value = "false"
  }
}