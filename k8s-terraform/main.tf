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
echo "[INIT] confluent-hub install --no-prompt mongodb/kafka-connect-mongodb:1.6.1"
confluent-hub install --no-prompt mongodb/kafka-connect-mongodb:1.6.1
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
    name  = "cp-zookeeper.imageTag"
    value = var.confluent_platform_version
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
    name  = "cp-kafka.imageTag"
    value = var.confluent_platform_version
  }
  set {
    name  = "cp-kafka.brokers"
    value = "3"
  }
  # see: https://medium.com/swlh/enable-external-access-to-confluent-kafka-on-kubernetes-step-by-step-e4647ca7a927
  # see: https://github.com/confluentinc/cp-helm-charts/blob/master/charts/cp-kafka/values.yaml#L137
  set {
    name  = "cp-kafka.nodeport.enabled"
    value = "true"
  }
  # in case you are not able to route to the k8s hostIp (kubectl get nodes -o wide), then you can manipulate the advertised listeners
  # to listen to localhost
  # https://tsuyoshiushio.medium.com/configuring-kafka-on-kubernetes-makes-available-from-an-external-client-with-helm-96e9308ee9f4
/*
  set {
    name = "cp-kafka.configurationOverrides.advertised.listeners"
    value = "EXTERNAL0://localhost:31090\\,EXTERNAL1://localhost:31091\\,EXTERNAL2://localhost:31092"
  }
  set {
    name = "cp-kafka.configurationOverrides.listener.security.protocol.map"
    value = "PLAINTEXT:PLAINTEXT\\,EXTERNAL0:PLAINTEXT\\,EXTERNAL1:PLAINTEXT\\,EXTERNAL2:PLAINTEXT"
  }
  set {
    name = "cp-kafka.configurationOverrides.listeners"
    value = "PLAINTEXT://:9092\\,EXTERNAL0://:31090\\,EXTERNAL1://:31091\\,EXTERNAL2://:31092"
  }
  set {
    name = "cp-kafka.configurationOverrides.inter.broker.listener.name"
    value = "PLAINTEXT"
  }
  */
  set {
    name  = "cp-kafka-connect.enabled"
    value = "true"
  }
  set {
    name  = "cp-kafka-connect.imageTag"
    value = var.confluent_platform_version
  }
  # install connect connectors
  # https://hub.docker.com/r/confluentinc/cp-kafka-connect
  # https://github.com/confluentinc/cp-helm-charts/tree/master/charts/cp-kafka-connect
  set {
    name = "cp-kafka-connect.image"
    value = "confluentinc/cp-kafka-connect" # default one. Add connectors via script
  }
  set {
    name = "cp-kafka-connect.replicaCount"
    value = "1"
  }
  # see: https://github.com/confluentinc/cp-helm-charts/blob/master/charts/cp-kafka-connect/templates/deployment.yaml#L111
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
    name  = "cp-ksql-server.imageTag"
    value = var.confluent_platform_version
  }
  set {
    name  = "cp-control-center.enabled"
    value = "true"
  }
  set {
    name  = "cp-control-center.imageTag"
    value = var.confluent_platform_version
  }
  set {
    name  = "cp-schema-registry.enabled"
    value = "true"
  }
  set {
    name  = "cp-schema-registry.imageTag"
    value = var.confluent_platform_version
  }
  set {
    name  = "cp-kafka-rest.enabled"
    value = "false"
  }
  set {
    name  = "cp-kafka-rest.imageTag"
    value = var.confluent_platform_version
  }
}

resource "helm_release" "postgresql-core-banking" {
  name       = "postgresql-core-banking"
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
  # See: https://github.com/bitnami/charts/blob/master/bitnami/postgresql/templates/_helpers.tpl#L105
  # See: https://github.com/bitnami/charts/blob/master/bitnami/postgresql/templates/statefulset.yaml#L205
  set {
    name  = "global.postgresql.replicationUser"
    value = var.postgresql_replication_username
  }
  # See: https://github.com/bitnami/charts/blob/master/bitnami/postgresql/values.yaml#L31
  set {
    name  = "global.postgresql.replicationPassword"
    value = var.postgresql_replication_password
  }

  # See: https://github.com/bitnami/charts/blob/master/bitnami/postgresql/values.yaml#L249
  set {
    name  = "extraEnv[0].name"
    value = "POSTGRESQL_WAL_LEVEL"
  }
  set {
    name  = "extraEnv[0].value"
    value = "logical"
  }
}

resource "helm_release" "mongodb-contact-center" {
  name       = "mongodb-contact-center"

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

resource "helm_release" "mongodb-customer-single-view" {
  name       = "mongodb-customer-single-view"

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