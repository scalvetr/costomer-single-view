-- CREATE new stream event_customer_entity
CREATE STREAM event_customer_entity
WITH (
    KEY_FORMAT='KAFKA',
    VALUE_FORMAT='AVRO',
    KAFKA_TOPIC='event_customer_entity'
);
-- CREATE new stream event_core_banking_bookings
CREATE STREAM event_core_banking_bookings
WITH (
    KEY_FORMAT='KAFKA',
    VALUE_FORMAT='AVRO',
    KAFKA_TOPIC='event_core_banking_bookings'
);
-- CREATE new stream event_core_banking_accounts
CREATE STREAM event_core_banking_accounts
WITH (
    KEY_FORMAT='KAFKA',
    VALUE_FORMAT='AVRO',
    KAFKA_TOPIC='event_core_banking_accounts'
);
-- CREATE new stream event_contact_center_customer_cases
CREATE STREAM event_contact_center_customer_cases
WITH (
    KEY_FORMAT='KAFKA',
    VALUE_FORMAT='AVRO',
    KAFKA_TOPIC='event_contact_center_customer_cases'
);