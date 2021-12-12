-- CREATE a table
CREATE STREAM "event_core_banking_accounts"
AS
SELECT *
FROM "raw_core_banking_accounts"
    -- KEY = ACCOUNT_ID
    PARTITION BY ACCOUNT_ID;

-- CREATE a table
CREATE STREAM "event_core_banking_bookings"
AS
SELECT *
FROM "raw_core_banking_bookings"
    -- KEY = ACCOUNT_ID
    PARTITION BY ACCOUNT_ID;

-- CREATE a table
CREATE STREAM "event_contact_center_customer_cases"
AS
SELECT *
FROM "raw_contact_center_customer_cases"
    -- KEY = ACCOUNT_ID
    PARTITION BY CASE_ID;
