-- Active: 1705213922423@@127.0.0.1@5432@linux
ALTER TABLE movies ADD CONSTRAINT movie_runtime_check CHECK(runtime>=0);
ALTER TABLE movies ADD CONSTRAINT movie_year_check CHECK(year BETWEEN 1888 AND DATE_PART('year', now()));

ALTER TABLE movies ADD CONSTRAINT genres_length_check CHECK(array_length(genres,1) BETWEEN 1 AND 5);