CREATE TABLE IF NOT EXISTS books (
  id SERIAL,
  title VARCHAR(255),
  synopsis TEXT,
  publish_date DATE,
  page_count INTEGER,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);