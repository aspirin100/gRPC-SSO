CREATE TABLE IF NOT EXISTS refresh_session
(
    sessionID uuid PRIMARY KEY,
    refreshToken BLOB NOT NULL,
    expiresAt timestamptz NOT NULL
);
CREATE TABLE IF NOT EXISTS users
(
    id        UUID PRIMARY KEY,
    refreshSessionID UUID,
    email     TEXT NOT NULL UNIQUE,
    passHash BLOB NOT NULL,
    FOREIGN KEY (refreshSessionID) REFERENCES refresh_session(sessionID)
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);
CREATE TABLE IF NOT EXISTS apps
(
    id     INTEGER PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE
);

