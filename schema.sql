CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE packages (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  version TEXT NOT NULL,
  description TEXT,
  authors TEXT[],
  dependencies JSONB,
  checksum TEXT NOT NULL,
  filename TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  UNIQUE (name, version)
);