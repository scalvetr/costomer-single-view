CREATE STREAM event_core_banking_accounts
AS
SELECT *
FROM raw_core_banking_accounts
PARTITION BY CUSTOMER_ID;

CREATE STREAM event_contact_center_customer_cases
AS
SELECT CUSTOMER_ID,
       MAP(
               'CASE_ID' := CASE_ID,
               'TITLE' := TITLE,
               'CREATION_TIMESTAMP' := TIMESTAMPTOSTRING(CREATION_TIMESTAMP, 'yyyy-MM-dd HH:mm:ss.SSS'),
               'COMMUNICATIONS' := CAST(
                       COMMUNICATIONS
                   AS STRING)
           ) AS `CASE`
FROM raw_contact_center_customer_cases
PARTITION BY CUSTOMER_ID;