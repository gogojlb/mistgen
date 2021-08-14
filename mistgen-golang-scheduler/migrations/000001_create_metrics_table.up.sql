BEGIN;

CREATE SCHEMA timeseries AUTHORIZATION mistgen; 

CREATE TABLE IF NOT EXISTS timeseries.measurements(
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NULL,
    temperature DOUBLE PRECISION NULL,
    humidity DOUBLE PRECISION NULL 
);
CREATE index on timeseries.measurements ("time");
