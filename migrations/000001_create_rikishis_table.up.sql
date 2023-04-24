CREATE TABLE IF NOT EXISTS rikishis (
    shikona text PRIMARY KEY NOT NULL,
    highest_rank text NOT NULL,
    heya text NOT NULL,
    shikona_history text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);