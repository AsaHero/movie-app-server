CREATE TABLE IF NOT EXISTS genres(
    id serial PRIMARY KEY,
    name varchar(100) UNIQUE NOT NULL
);

