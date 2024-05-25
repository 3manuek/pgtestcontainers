-- This SQL script generates a CSV data with random data.
-- 

-- This is case study, not for being a good practice for production modeling.
-- Even further, it should be using at least uuidv7. 
DO $createDevices$
DECLARE
    r record;
BEGIN
  WITH uuids AS (
    SELECT $$'$$ || gen_random_uuid()::text || $$'$$ dev FROM generate_series(1,100)
  )
  SELECT string_agg(dev, ',') as devices INTO r FROM uuids;
  EXECUTE FORMAT('CREATE TYPE tdevice AS ENUM (%s);', r.devices); 
END;
$createDevices$;

CREATE FUNCTION payloadGen() RETURNS table(temp double precision, cpu double precision)
    AS $$
    SELECT round(random_normal(0.9,0.8)*20)  , 
            round(random_normal(0.9,0.8)*10+50) 
$$  LANGUAGE SQL VOLATILE;

COPY (
    WITH tss(ts,n, payload) AS (
        SELECT t.ts, d.n ,  
            row_to_json(payloadGen()) as payload
            FROM generate_series('2021-01-01 00:00:00'::timestamp, 
                                    '2021-01-02 00:00:00'::timestamp, 
                                    '1 minute') t(ts),
                -- we do lateral so we iterate for each ts all the devices
                LATERAL (SELECT enumlabel FROM pg_enum WHERE enumtypid = 'tdevice'::regtype) d(n)
    )
    select ts,n, payload::jsonb from tss
) TO '/tmp/devices.csv' DELIMITER ',' CSV HEADER;

DROP FUNCTION IF EXISTS payloadGen();

