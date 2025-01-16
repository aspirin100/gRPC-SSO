CREATE TABLE IF NOT EXISTS users
(
    id        UUID PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    passHash BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS refresh_session
(
    userID uuid PRIMARY KEY,
    refreshToken BLOB NOT NULL,
    expiresAt timestamptz NOT NULL,
    isUsed BOOL DEFAULT FALSE,
    FOREIGN KEY (userID) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS apps
(
    id     INTEGER PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE
);

