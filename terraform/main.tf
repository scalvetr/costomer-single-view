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
  debug = true
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
resource "kubernetes_config_map" "cp-kafka-connect-init-script" {
  metadata {
    name = "cp-kafka-connect-init-script"
    namespace = var.k8s_namespace
  }
  data = {
    "init-script.sh" = <<EOF
#!/bin/bash
echo "[INIT] Installing Additional Connectors"
echo "[INIT] confluent-hub install --no-prompt debezium/debezium-connector-postgresql:1.7.1"
confluent-hub install --no-prompt debezium/debezium-connector-postgresql:1.7.1
echo "[INIT] confluent-hub install --no-prompt debezium/debezium-connector-mongodb:1.7.1"
confluent-hub install --no-prompt debezium/debezium-connector-mongodb:1.7.1
EOF
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
  # install connect connectors
  # https://hub.docker.com/r/confluentinc/cp-kafka-connect
  # https://github.com/confluentinc/cp-helm-charts/tree/master/charts/cp-kafka-connect
  set {
    name = "cp-kafka-connect.image"
    value = "confluentinc/cp-kafka-connect"
  }
  set {
    name  = "cp-kafka-connect.customEnv.CUSTOM_SCRIPT_PATH"
    value = "/etc/config/init-script.sh"
  }
  set {
    name  = "cp-kafka-connect.volumes[0].name"
    value = "config-volume"
  }
  set {
    name  = "cp-kafka-connect.volumes[0].configMap.name"
    value = "cp-kafka-connect-init-script"
  }
  set {
    name  = "cp-kafka-connect.volumes[0].configMap.defaultMode"
    # value 0777 in decimal notation. Check here: https://kubernetes.io/docs/concepts/configuration/secret/#secret-files-permissions
    value = 511
  }
  set {
    name  = "cp-kafka-connect.volumeMounts[0].name"
    value = "config-volume"
  }
  set {
    name  = "cp-kafka-connect.volumeMounts[0].mountPath"
    value = "/etc/config"
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