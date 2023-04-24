CREATE TABLE IF NOT EXISTS bouts (
    id bigserial PRIMARY KEY,
    tournament date NOT NULL,
    day text NOT NULL,
    winner text REFERENCES rikishis(shikona),
    loser text REFERENCES rikishis(shikona),
    kimarite text,
    version integer NOT NULL DEFAULT 1
);