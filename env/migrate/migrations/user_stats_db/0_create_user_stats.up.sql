BEGIN;

CREATE TABLE user_game_stats (
    user_id TEXT NOT NULL,
    game_id TEXT NOT NULL,
    outcome TEXT NOT NULL
);

COMMIT;
