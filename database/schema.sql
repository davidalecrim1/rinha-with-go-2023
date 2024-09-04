CREATE TABLE IF NOT EXISTS persons (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    dob TEXT NOT NULL,
    stack TEXT,
    search TEXT GENERATED ALWAYS AS (
        nickname || ' ' || name || ' ' || stack
    ) STORED
);

CREATE INDEX idx_persons_search ON persons (search);