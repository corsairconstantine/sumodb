CREATE INDEX IF NOT EXISTS rikishis_shikona_idx ON rikishis USING GIN (to_tsvector('simple', shikona));
