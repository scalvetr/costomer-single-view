SHOW STREAMS;

CREATE SOURCE CONNECTOR `contact-center-source` WITH(
    'connection.uri'='mongodb://user:password@mongodb-contact-center/contact-center',
    'connector.class'='com.mongodb.kafka.connect.MongoSourceConnector',
    'database'='contact-center',
    'collection'='cases',
    'topic.namespace.map'='{"contact-center":"raw.contact_center.default","contact-center.cases":"raw.contact_center.customer_cases"}',
    'publish.full.document.only'=true,
    'output.format.value'='schema',
    'output.schema.value'='{"namespace":"edu.uoc.scalvetr","type":"record","name":"Case","fields":[{"name":"_id","type":"string"},{"name":"case_id","type":"string"},{"name":"customer_id","type":"string"},{"name":"title","type":"string"},{"name":"creation_timestamp","type":"string"},{"name":"communications","type":{"type":"array","items":{"name":"Communication","type":"record","fields":[{"name":"communication_id","type":"string"},{"name":"type","type":"string"},{"name":"text","type":"string"},{"name":"notes","type":"string"}]}}}]}',
    'value.converter.schemas.enable'='true',
    'key.converter'='org.apache.kafka.connect.storage.StringConverter',
    'value.converter'='io.confluent.connect.avro.AvroConverter',
    'value.converter.schema.registry.url'='http://schema-registry:8081');

CREATE SOURCE CONNECTOR `core-banking-source` WITH(
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
    'transforms'='Reroute',
    'transforms.Reroute.type'='io.debezium.transforms.ByLogicalTableRouter',
    'transforms.Reroute.topic.regex'='core-banking.public\.(.*)',
    'transforms.Reroute.topic.replacement'='raw.core_banking.$1',
    'key.converter'='org.apache.kafka.connect.storage.StringConverter',
    'value.converter'='io.confluent.connect.avro.AvroConverter',
    'value.converter.schema.registry.url'='http://schema-registry:8081'
);