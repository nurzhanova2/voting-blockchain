-- +goose Up
-- Создаёт таблицы пользователей, голосований, голосов и блоков

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS elections (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    election_id INTEGER REFERENCES elections(id),
    choice TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(user_id, election_id)
);

CREATE TABLE blockchain (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    vote_hash TEXT NOT NULL,
    previous_hash TEXT,
    current_hash TEXT NOT NULL,
    election_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
-- Удаляет все таблицы

DROP TABLE IF EXISTS blockchain;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS elections;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;

