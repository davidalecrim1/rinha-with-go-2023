CREATE TABLE IF NOT EXISTS persons (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL,
    name TEXT NOT NULL,
    dob TEXT NOT NULL,
    stack TEXT[]
);