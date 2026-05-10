CREATE TABLE game_members (
    id BIGSERIAL PRIMARY KEY,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id),
    added_by_user_id BIGINT REFERENCES users(id),
    name TEXT NOT NULL,
    status TEXT NOT NULL
        CHECK (status IN ('active', 'waitlist', 'cancelled')),
    source TEXT NOT NULL
        CHECK (source IN ('telegram', 'web', 'admin')),
    position_number INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    cancelled_at TIMESTAMPTZ
);
