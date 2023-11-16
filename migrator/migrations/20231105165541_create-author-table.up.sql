CREATE TABLE IF NOT EXISTS authors (
  id SERIAL,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  bio TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);