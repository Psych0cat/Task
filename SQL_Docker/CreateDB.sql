CREATE TABLE if not exists balance (
id              SERIAL UNIQUE,
balance         decimal(15,2),
PRIMARY KEY (id)
);

CREATE TABLE if not exists source (
id              SERIAL UNIQUE,
source_name     varchar(50),
PRIMARY KEY (id)
);

CREATE TABLE if not exists transaction_queue (
id              SERIAL,
source          int NOT NULL,
state           varchar NOT NULL,
amount          float NOT NULL,
unix_timestamp  bigint NOT NULL,
transactionId   varchar NOT NULL,
cancelled        boolean,
balance_id      int NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (source) REFERENCES source(id),
FOREIGN KEY (balance_id) REFERENCES balance(id)
);



