BEGIN;

CREATE TABLE games (
    id serial PRIMARY KEY,
    white_user_id TEXT,
    black_user_id TEXT
);

CREATE TABLE moves (
    id serial PRIMARY KEY,
    game_id serial,
    user_id TEXT,
    move TEXT
);

COMMIT;
