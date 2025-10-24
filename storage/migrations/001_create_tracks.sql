CREATE TABLE IF NOT EXISTS tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    genre TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_path TEXT UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_artist ON tracks (artist);
CREATE INDEX IF NOT EXISTS idx_genre ON tracks (genre);
