CREATE TABLE IF NOT EXISTS users (
    id      BIGINT PRIMARY KEY,
    name    VARCHAR(100),
    skill   DOUBLE PRECISION,
    latency DOUBLE PRECISION,
    added   TIME
);