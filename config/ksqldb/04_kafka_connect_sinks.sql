SHOW STREAMS;

-- https://github.com/mongodb/mongo-kafka/blob/master/config/MongoSinkConnector.properties
CREATE
SINK CONNECTOR `customer-one-view-sink`
WITH (
    'connection.uri'='mongodb://user:password@mongodb-customer-single-view/customer-single-view',
    'connector.class'='com.mongodb.kafka.connect.MongoSinkConnector',
    'database'='customer-single-view',

    'key.converter'='org.apache.kafka.connect.storage.StringConverter',
    'value.converter'='io.confluent.connect.avro.AvroConverter',
    'value.converter.schema.registry.url'='http://schema-registry:8081',

    'key.projection.type'='none',
    'key.projection.list'='',
    'value.projection.type'='none',
    'value.projection.list'='',

    'topics'='event_customer_entity,event_core_banking_accounts,event_core_banking_bookings,event_contact_center_customer_cases',
    'collection'='customers',

    -- Collection names
    'topic.override.event_core_banking_accounts.collection'='accounts',
    'topic.override.event_core_banking_bookings.collection'='account_bookings',
    'topic.override.event_contact_center_customer_cases.collection'='cases',

    --# Write configuration
    -- 'delete.on.null.values'='true',
    -- 'writemodel.strategy'='com.mongodb.kafka.connect.sink.writemodel.strategy.ReplaceOneDefaultStrategy',
    -- Id Strategy
    'post.processor.chain'='com.mongodb.kafka.connect.sink.processor.DocumentIdAdder',
    'document.id.strategy'='com.mongodb.kafka.connect.sink.processor.id.strategy.FullKeyStrategy',
    'topic.override.event_core_banking_bookings.document.id.strategy'='com.mongodb.kafka.connect.sink.processor.id.strategy.ProvidedInValueStrategy'
);