CREATE INDEX idx_games_starts_at ON games(starts_at);
CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_game_members_game_id ON game_members(game_id);
CREATE INDEX idx_game_members_game_status ON game_members(game_id, status);
CREATE INDEX idx_users_telegram_id ON users(telegram_id);

CREATE UNIQUE INDEX uniq_game_member_user_active
ON game_members(game_id, user_id)
WHERE user_id IS NOT NULL AND status IN ('active', 'waitlist');
