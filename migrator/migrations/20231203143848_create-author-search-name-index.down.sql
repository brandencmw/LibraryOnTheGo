-- Remove name_search_idx
DROP INDEX IF EXISTS name_search_idx;

-- Remove name_search_vector column from authors table
ALTER TABLE authors DROP COLUMN IF EXISTS name_search_vector;