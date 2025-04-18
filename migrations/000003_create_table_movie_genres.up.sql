CREATE TABLE IF NOT EXISTS movie_genres(
    movie_id bigint NOT NULL,
    genre_id int NOT NULL,
    created_at timestamptz DEFAULT now(),
    PRIMARY KEY (movie_id, genre_id),
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_movie_genres_movie_id ON movie_genres(movie_id);

CREATE INDEX IF NOT EXISTS idx_movie_genres_genre_id ON movie_genres(genre_id);

