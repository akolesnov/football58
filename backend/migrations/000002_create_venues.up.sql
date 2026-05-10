CREATE TABLE venues (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    address TEXT
);
