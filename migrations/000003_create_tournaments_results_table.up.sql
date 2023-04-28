CREATE TABLE IF NOT EXISTS tournaments_results (
    id bigserial PRIMARY KEY,
    tournament text NOT NULL,
    rikishi text REFERENCES rikishis(shikona),
    rank text NOT NULL,
    wins integer NOT NULL,
    losses integer NOT NULL,
    absent integer NOT NULL,
    version integer NOT NULL DEFAULT 1
);