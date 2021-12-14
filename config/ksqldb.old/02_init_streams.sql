SHOW STREAMS;

CREATE STREAM "raw_core_banking_accounts" WITH (
    kafka_topic = 'raw_core_banking_accounts',
    value_format = 'avro'
);
CREATE STREAM "raw_core_banking_bookings" WITH (
    kafka_topic = 'raw_core_banking_bookings',
    value_format = 'avro'
);

CREATE STREAM "raw_contact_center_customer_cases" WITH (
    kafka_topic = 'raw_contact_center_customer_cases',
    value_format = 'avro'
);

CREATE STREAM "event_customer_entity" WITH (
    kafka_topic = 'event_customer_entity',
    value_format = 'avro'
);
