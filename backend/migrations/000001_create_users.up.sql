CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    telegram_id BIGINT UNIQUE,
    telegram_username TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
