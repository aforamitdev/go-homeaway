-- Active: 1705213922423@@127.0.0.1@5432@linux
ALTER TABLE movies ALTER COLUMN id add generated always as identity;