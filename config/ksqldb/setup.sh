#!/bin/sh

echo "Run \"01_kafka_connect_sources.sql\""
/bin/ksql --file /config/01_kafka_connect_sources.sql -- http://ksqldb-server:8088

echo "Wait 1m"
sleep 1m

echo "Run \"02_create_streams.sql\""
/bin/ksql --file /config/02_create_streams.sql -- http://ksqldb-server:8088

echo "Run \"03_kafka_connect_sinks.sql\""
/bin/ksql --file /config/03_kafka_connect_sinks.sql -- http://ksqldb-server:8088