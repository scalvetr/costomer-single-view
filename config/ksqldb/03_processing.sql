CREATE STREAM event_core_banking_accounts
AS
SELECT c.ACCOUNT_ID        AS ACCOUNT_ID,
       c.CUSTOMER_ID       AS CUSTOMER_ID,
       c.IBAN              AS IBAN,
       c.BALANCE           AS BALANCE,
       c.CREATION_DATE     AS CREATION_DATE,
       c.CANCELLATION_DATE AS CANCELLATION_DATE,
       c.STATUS            AS STATUS
FROM raw_core_banking_accounts c PARTITION BY c.CUSTOMER_ID;

CREATE STREAM event_contact_center_customer_cases
AS
SELECT *
FROM raw_contact_center_customer_cases c PARTITION BY c.CUSTOMER_ID;