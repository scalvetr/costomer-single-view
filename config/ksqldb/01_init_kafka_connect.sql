SHOW STREAMS;

CREATE
SOURCE CONNECTOR `contact-center-source`
WITH (
    'connection.uri'='mongodb://user:password@mongodb-contact-center/contact-center',
    'connector.class'='com.mongodb.kafka.connect.MongoSourceConnector',
    'database'='contact-center',
    'collection'='cases',
    'topic.namespace.map'='{"contact-center":"raw.contact_center.default","contact-center.cases":"raw.contact_center.customer_cases"}',
    'publish.full.document.only'= true,
    -- Partition by customer: https://www.mongodb.com/blog/post/mongo-db-connector-for-apache-kafka-1-3-available-now
    'output.format.key'='schema',
    'output.schema.key'='{"name":"CustomerId","type":"record","namespace":"edu.uoc.scalvetr","fields":[{"name":"fullDocument","type":{"name":"fullDocument","type":"record","fields":[{"name":"customer_id","type":"string"}]}}]}',
    'output.format.value'='schema',
    'output.schema.value'='{"namespace":"edu.uoc.scalvetr","type":"record","name":"Case","fields":[{"name":"_id","type":"string"},{"name":"case_id","type":"string"},{"name":"customer_id","type":"string"},{"name":"title","type":"string"},{"name":"creation_timestamp","type":"string"},{"name":"communications","type":{"type":"array","items":{"name":"Communication","type":"record","fields":[{"name":"communication_id","type":"string"},{"name":"type","type":"string"},{"name":"text","type":"string"},{"name":"notes","type":"string"}]}}}]}',
    'value.converter.schemas.enable'='true',
    'key.converter'='io.confluent.connect.avro.AvroConverter',
    'key.converter.schema.registry.url'='http://schema-registry:8081',
    'value.converter'='io.confluent.connect.avro.AvroConverter',
    'value.converter.schema.registry.url'='http://schema-registry:8081'
);

CREATE
SOURCE CONNECTOR `core-banking-source`
WITH (
    'connector.class'='io.debezium.connector.postgresql.PostgresConnector',
    'database.server.name'='core-banking',
    'database.hostname'='postgresql-core-banking',
    'database.port'='5432',
    'database.user'='user',
    'database.password'='password',
    'database.dbname'='core-banking',
    'plugin.name'='wal2json',
    'table.include.list'='public.accounts,public.bookings',
    'topic.creation.default.partitions'='1',
    'topic.creation.default.replication.factor'='1',
    'topic.creation.enable'='true',
    'transforms'='reroute,unwrap',
    'transforms.reroute.type'='io.debezium.transforms.ByLogicalTableRouter',
    'transforms.reroute.topic.regex'='core-banking.public\.(.*)',
    'transforms.reroute.topic.replacement'='raw.core_banking.$1',
    'transforms.reroute.key.field.name'='table',
    'transforms.reroute.key.enforce.uniqueness'= false,
    'transforms.unwrap.type'='io.debezium.transforms.ExtractNewRecordState',
    'transforms.unwrap.drop.tombstones'='false',
    'transforms.unwrap.delete.handling.mode'='rewrite',
    'value.converter'='io.confluent.connect.avro.AvroConverter',
    'value.converter.schema.registry.url'='http://schema-registry:8081',
    -- Partition by customer: https://debezium.io/documentation/reference/stable/transformations/topic-routing.html#_example
    'message.key.columns'='(.*).accounts:customer_id,customer_id;(.*).bookings:customer_id,customer_id',
    'key.converter'='io.confluent.connect.avro.AvroConverter',
    'key.converter.schema.registry.url'='http://schema-registry:8081'
);