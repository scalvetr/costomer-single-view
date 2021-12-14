-- CREATE new stream keyed by ACCOUNT_ID
CREATE STREAM "event_core_banking_accounts"
WITH (
    KEY_FORMAT='KAFKA',
    VALUE_FORMAT='AVRO'
)
AS
SELECT
    account_id AS "account_id",
    customer_id AS "customer_id",
    iban AS "iban",
    balance AS "balance",
    creation_date AS "creation_date",
    cancellation_date AS "cancellation_date",
    status AS "status"
FROM "raw_core_banking_accounts"
    -- KEY = ACCOUNT_ID
    PARTITION BY account_id;

-- CREATE new stream, as topic is keyed by account_id _id field has
-- been added, adn will be used as object identifier in the downstream mongodb
CREATE STREAM "event_core_banking_bookings"
WITH (
    KEY_FORMAT='KAFKA',
    VALUE_FORMAT='AVRO'
)
AS
SELECT
    booking_id AS "_id",
    booking_id AS "booking_id",
    account_id AS "account_id",
    amount AS "amount",
    description AS "description",
    booking_date AS "booking_date",
    value_date AS "value_date",
    fee AS "fee",
    taxes AS  "taxes"
FROM "raw_core_banking_bookings"
    -- KEY = ACCOUNT_ID
PARTITION BY account_id;

-- CREATE new stream keyed by CASE_ID
CREATE STREAM "event_contact_center_customer_cases"
WITH (
    KEY_FORMAT='KAFKA',
    VALUE_FORMAT='AVRO'
)
AS
SELECT
    case_id AS "case_id",
    customer_id AS "customer_id",
    title AS "title",
    creation_timestamp AS "creation_timestamp",
    communications AS "communications"
FROM "raw_contact_center_customer_cases"
    -- KEY = ACCOUNT_ID
    PARTITION BY case_id;
