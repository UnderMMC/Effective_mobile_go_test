CREATE TABLE music (
    song_id SERIAL PRIMARY KEY,
    performer TEXT NOT NULL,
    song_name TEXT NOT NULL,
    release_date TEXT,
    lyric TEXT,
    link TEXT
);