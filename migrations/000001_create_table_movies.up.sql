CREATE TABLE IF NOT EXISTS movies(
    id bigserial PRIMARY KEY,
    title varchar(255) NOT NULL,
    release date,
    plot text,
    duration_minutes int,
    poster_url text,
    trailer_url text,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
);

