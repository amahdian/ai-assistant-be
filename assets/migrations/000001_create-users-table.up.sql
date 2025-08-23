BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     email TEXT NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Explicit index for faster lookups by email
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

COMMIT;
