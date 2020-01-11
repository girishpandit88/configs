CREATE TABLE IF NOT EXISTS configs
(
    key        text unique PRIMARY KEY,
    values     jsonb,
    createdAt  timestamp with time zone,
    modifiedAt timestamp with time zone
);
CREATE UNIQUE INDEX IF NOT EXISTS keyidx ON configs (key);
CREATE INDEX IF NOT EXISTS idxgin ON configs USING gin (values);
