
CREATE TABLE IF NOT EXISTS devices (
    ts timestamp NOT NULL,
    -- device tdevice NOT NULL, -- we won't use tdevice yet
    device uuid NOT NULL,
    payload jsonb NOT NULL
);

COPY devices FROM '/tmp/devices.csv' DELIMITER ',' CSV HEADER;