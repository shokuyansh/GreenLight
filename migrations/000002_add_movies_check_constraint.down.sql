ALTER TABLE movies DROP CONSTRAINT IF exists movies_runtime_check;
ALTER TABLE movies DROP CONSTRAINT IF exists movies_year_check;
ALTER TABLE movies DROP CONSTRAINT IF exists genres_length_check;
