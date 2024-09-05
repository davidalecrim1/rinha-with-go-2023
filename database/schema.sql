CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS people (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    dob TEXT NOT NULL,
    stack TEXT,
    searchable TEXT GENERATED ALWAYS AS (
        nickname || ' ' || name || ' ' || lower(stack)
    ) STORED
);

-- I will test using GiST and GIN
CREATE INDEX
    CONCURRENTLY IF NOT EXISTS idx_people_searchable ON public.people USING gist
    (searchable public.gist_trgm_ops);

-- CREATE INDEX
--     CONCURRENTLY IF NOT EXISTS idx_people_searchable ON public.people USING gin
--     (searchable public.gin_trgm_ops);