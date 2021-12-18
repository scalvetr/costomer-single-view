SHOW STREAMS;

CREATE
SOURCE CONNECTOR `contact-center-source`
WITH (
    'connection.uri'='mongodb://user:password@mongodb-contact-center/contact-center',
    'connector.class'='com.mongodb.kafka.connect.MongoSourceConnector',
    'database'='contact-center',
    'collection'='cases',
    'topic.namespace.map'='{"contact-center":"event_contact_center_default","contact-center.cases":"event_contact_center_customer_cases"}',
    'publish.full.document.only'= true,
    'output.format.value'='schema',
    'output.schema.value'='{"namespace":"edu.uoc.scalvetr","type":"record","name":"Case","fields":[{"name":"case_id","type":"string"},{"name":"customer_id","type":"string"},{"name":"title","type":"string"},{"name":"creation_timestamp","type":"long","logicalType":"timestamp-millis"},{"name":"communications","type":{"type":"array","items":{"name":"Communication","type":"record","fields":[{"name":"communication_id","type":"string"},{"name":"type","type":"string"},{"name":"text","type":"string"},{"name":"notes","type":"string"},{"name":"timestamp","type":"long","logicalType":"timestamp-millis"}]}}}]}',
    'value.converter.schemas.enable'='true',
    'value.converter'='io.confluent.connect.avro.AvroConverter',
    'value.converter.schema.registry.url'='http://schema-registry:8081',

    --'output.format.key'='schema',
    --'output.schema.key'='{"namespace":"edu.uoc.scalvetr","type":"record","name":"CaseId","fields":[{"name":"case_id","type":"string"}]}',
    --'output.schema.key'='{"namespace":"edu.uoc.scalvetr","type":"record","name":"CaseId","fields":[{"name":"fullDocument","type":{"name":"fullDocument","type":"record","fields":[{"name":"case_id","type":"string"}]}}]}',
    --'key.converter'='io.confluent.connect.avro.AvroConverter',
    --'key.converter.schema.registry.url'='http://schema-registry:8081',
    'key.converter'='org.apache.kafka.connect.storage.StringConverter',


    'transforms'='valueToKey,extractKey',
    'transforms.valueToKey.type'='org.apache.kafka.connect.transforms.ValueToKey',
    'transforms.valueToKey.fields'='case_id',
    'transforms.extractKey.type'='org.apache.kafka.connect.transforms.ExtractField$Key',
    'transforms.extractKey.field'='case_id'
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
    --https://debezium.io/documentation/reference/stable/connectors/postgresql.html#postgresql-decimal-types
    'decimal.handling.mode'='double',
    'time.precision.mode'='adaptive',

    'value.converter'='io.confluent.connect.avro.AvroConverter',
    'value.converter.schema.registry.url'='http://schema-registry:8081',
    -- Partition by customer: https://debezium.io/documentation/reference/stable/transformations/topic-routing.html#_example
    'message.key.columns'='(.*).accounts:account_id;(.*).bookings:account_id',
    'key.converter'='org.apache.kafka.connect.storage.StringConverter',
    --'key.converter'='io.confluent.connect.avro.AvroConverter',
    --'key.converter.schema.registry.url'='http://schema-registry:8081',

    'transforms'='reroute,unwrap,extractKey',
    'transforms.reroute.type'='io.debezium.transforms.ByLogicalTableRouter',
    'transforms.reroute.topic.regex'='core-banking.public\.(.*)',
    'transforms.reroute.topic.replacement'='event_core_banking_$1',
    'transforms.reroute.key.field.name'='table',
    'transforms.reroute.key.enforce.uniqueness'= false,
    'transforms.unwrap.type'='io.debezium.transforms.ExtractNewRecordState',
    'transforms.extractKey.type'='org.apache.kafka.connect.transforms.ExtractField$Key',
    'transforms.extractKey.field'='account_id'

);