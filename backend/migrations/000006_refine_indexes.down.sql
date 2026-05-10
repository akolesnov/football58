DROP INDEX idx_games_status_starts_at;

CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_game_members_game_id ON game_members(game_id);
CREATE INDEX idx_users_telegram_id ON users(telegram_id);
