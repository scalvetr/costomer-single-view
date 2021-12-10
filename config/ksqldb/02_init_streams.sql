CREATE STREAM raw_core_banking_accounts WITH (
    kafka_topic = 'raw.core_banking.accounts',
    value_format = 'avro'
);
CREATE STREAM raw_core_banking_bookings WITH (
    kafka_topic = 'raw.core_banking.bookings',
    value_format = 'avro'
);

CREATE STREAM raw_contact_center_customer_cases WITH (
    kafka_topic = 'raw.contact_center.customer_cases',
    value_format = 'avro'
);

CREATE STREAM event_customer_entity WITH (
    kafka_topic = 'event.customer.entity',
    value_format = 'avro'
);
