CREATE TABLE IF NOT EXISTS persons (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    dob TEXT NOT NULL,
    stack TEXT,
    searchable TEXT GENERATED ALWAYS AS (
        nickname || ' ' || name || ' ' || lower(stack)
    ) STORED
);

CREATE INDEX CONCURRENTLY idx_persons_searchable ON persons (searchable);