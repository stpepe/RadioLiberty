CREATE TABLE IF NOT EXISTS queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_name TEXT NOT NULL unique,
    audio_name TEXT NOT NULL,
    artist TEXT NOT NULL
);