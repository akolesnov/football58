DROP INDEX idx_users_telegram_id;
DROP INDEX idx_game_members_game_id;
DROP INDEX idx_games_status;

CREATE INDEX idx_games_status_starts_at ON games(status, starts_at);
