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

    -- Id Strategy
    'post.processor.chain'='com.mongodb.kafka.connect.sink.processor.DocumentIdAdder',
    'document.id.strategy'='com.mongodb.kafka.connect.sink.processor.id.strategy.ProvidedInKeyStrategy', -- _id field in key


    -- Collection names
    'collection'='customers',
    'topic.override.event_core_banking_accounts.collection'='accounts',
    'topic.override.event_core_banking_bookings.collection'='account_bookings',
    'topic.override.event_contact_center_customer_cases.collection'='cases',

    'transforms'='valueToKey,replaceField',
    'transforms.valueToKey.type'='org.apache.kafka.connect.transforms.ValueToKey',
    'transforms.replaceField.type'='org.apache.kafka.connect.transforms.ReplaceField$Key',

    -- topic.override.event_customer_entity.
    'transforms.valueToKey.fields'='customerId',
    'transforms.replaceField.fields'='customerId:_id',
    -- topic.override.event_core_banking_accounts.
    'topic.override.event_core_banking_accounts.transforms.valueToKey.fields'='account_id',
    'topic.override.event_core_banking_accounts.transforms.replaceField.fields'='account_id:_id',
    -- topic.override.event_core_banking_bookings.
    'topic.override.event_core_banking_bookings.transforms.valueToKey.fields'='booking_id',
    'topic.override.event_core_banking_bookings.transforms.replaceField.fields'='booking_id:_id',
    -- topic.override.event_contact_center_customer_cases.
    'topic.override.event_contact_center_customer_cases.transforms.valueToKey.fields'='case_id',
    'topic.override.event_contact_center_customer_cases.transforms.replaceField.fields'='case_id:_id'

);