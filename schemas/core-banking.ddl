CREATE TABLE account IN NOT EXISTS
(
    account_id NUMBER NOT NULL PRIMARY KEY,
    customer_id VARCHAR(20) NOT NULL,
    iban VARCHAR(35) NOT NULL,
    balance NUMBER(20,3) NOT NULL,
    creation_date DATE NOT NULL,
    cancellation_date DATE,
    status VARCHAR(10) NOT NULL
)
CREATE TABLE booking IN NOT EXISTS
(
    booking_id NUMBER NOT NULL PRIMARY KEY,
    account_id NUMBER NOT NULL,
    ammount NUMBER(20,3) NOT NULL,
    description VARCHAR(250) NOT NULL,
    booking_date DATE NOT NULL,
    value_date DATE,
    fee NUMBER(20,3) NOT NULL DEFAULT 0,
    taxes NUMBER(20,3) NOT NULL DEFAULT 0
    )