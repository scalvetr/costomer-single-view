#!/bin/bash
echo "[INIT] Installing Additional Connectors"
echo "[INIT] confluent-hub install --no-prompt debezium/debezium-connector-postgresql:1.7.1"
confluent-hub install --no-prompt debezium/debezium-connector-postgresql:1.7.1
#echo "[INIT] confluent-hub install --no-prompt debezium/debezium-connector-mongodb:1.7.1"
#confluent-hub install --no-prompt debezium/debezium-connector-mongodb:1.7.1
echo "[INIT] confluent-hub install --no-prompt mongodb/kafka-connect-mongodb:1.6.1"
confluent-hub install --no-prompt mongodb/kafka-connect-mongodb:1.6.1
#
echo "Launching Kafka Connect worker"
/etc/confluent/docker/run &
#
sleep infinity