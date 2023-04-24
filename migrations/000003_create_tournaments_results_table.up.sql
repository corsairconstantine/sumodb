CREATE TABLE IF NOT EXISTS tournaments_results (
    id bigserial PRIMARY KEY,
    tournament date NOT NULL,
    rikishi text REFERENCES rikishis(shikona),
    rank text NOT NULL,
    result text NOT NULL,
    version integer NOT NULL DEFAULT 1
);