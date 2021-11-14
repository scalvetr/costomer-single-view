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
  version    = "0.6.1"
  namespace = var.k8s_namespace

  set {
    name  = "cp-zookeeper.enabled"
    value = "true"
  }
  set {
    name  = "cp-zookeeper.servers"
    value = "3"
  }
  set {
    name  = "cp-kafka.enabled"
    value = "true"
  }
  set {
    name  = "cp-kafka.brokers"
    value = "3"
  }
  set {
    name  = "cp-kafka-connect.enabled"
    value = "true"
  }
  set {
    name  = "cp-ksql-server.enabled"
    value = "true"
  }
  set {
    name  = "cp-control-center.enabled"
    value = "true"
  }
  set {
    name  = "cp-schema-registry.enabled"
    value = "true"
  }
  set {
    name  = "cp-kafka-rest.enabled"
    value = "false"
  }
}

resource "helm_release" "postgresql" {
  name       = "postgresql"

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "postgresql"
  version    = "10.13.8"
  namespace = var.k8s_namespace

  # https://github.com/bitnami/charts/tree/master/bitnami/postgresql/#parameters
  set {
    name  = "postgresqlUsername"
    value = var.postgresql_username
  }
  set {
    name  = "postgresqlPassword"
    value = var.postgresql_password
  }
}

resource "helm_release" "mongodb" {
  name       = "mongodb"

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "mongodb"
  version    = "10.29.2"
  namespace  = var.k8s_namespace

  values = [
  ]
  # https://github.com/bitnami/charts/tree/master/bitnami/mongodb/#parameters
  set {
    name  = "auth.database"
    value = var.mongodb_database
  }
  set {
    name  = "auth.username"
    value = var.mongodb_username
  }
  set {
    name  = "auth.password"
    value = var.mongodb_password
  }
}