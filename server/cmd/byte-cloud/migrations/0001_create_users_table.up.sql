CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name           TEXT        NOT NULL,
    email          TEXT UNIQUE NOT NULL,
    password  TEXT, -- NULL for Google users
    provider       TEXT        NOT NULL CHECK (provider IN ('local', 'google')),
    provider_id    TEXT, -- Google `sub` field, NULL for local users
    picture_url    TEXT,
    email_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
