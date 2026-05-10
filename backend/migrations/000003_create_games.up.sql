CREATE TABLE games (
    id BIGSERIAL PRIMARY KEY,
    venue_id BIGINT NOT NULL REFERENCES venues(id),
    starts_at TIMESTAMPTZ NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 120,
    min_players INT NOT NULL DEFAULT 10,
    max_players INT NOT NULL DEFAULT 15,
    price_rub INT NOT NULL DEFAULT 220,
    notes TEXT,
    status TEXT NOT NULL DEFAULT 'open'
        CHECK (status IN ('open', 'closed', 'cancelled', 'finished')),
    created_by_user_id BIGINT REFERENCES users(id),
    telegram_chat_id BIGINT,
    telegram_message_id BIGINT,
    version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
