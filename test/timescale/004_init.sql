CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE IF NOT EXISTS devices (
    ts timestamp NOT NULL,
    -- device tdevice NOT NULL, -- we won't use tdevice yet
    device uuid NOT NULL,
    payload jsonb NOT NULL
);

SELECT create_hypertable('devices', 'ts');

CREATE MATERIALIZED VIEW IF NOT EXISTS device_cpu_by_minute 
WITH (timescaledb.continuous) AS 
SELECT 
    time_bucket('1 minute', ts) AS minute,
    device,
    max(payload->>'cpu') AS cpu
FROM devices
GROUP BY minute, device;