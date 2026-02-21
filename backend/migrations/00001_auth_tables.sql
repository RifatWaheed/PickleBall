-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  mobile text NOT NULL UNIQUE,
  email text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  role text NOT NULL DEFAULT 'user',
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  revoked_at timestamptz NULL,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;