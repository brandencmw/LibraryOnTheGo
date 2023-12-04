-- Add name_search_vector column to authors table
ALTER TABLE authors ADD COLUMN name_search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', COALESCE(first_name, '') || ' ' || COALESCE(last_name, ''))) STORED;

-- Populate column for existing entries
UPDATE authors SET name_search_vector = DEFAULT;

-- Create index for name search
CREATE INDEX name_search_idx ON authors USING GIN(name_search_vector);