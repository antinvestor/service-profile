-- Drop old composite index (profile_id, contact_id only)
DROP INDEX IF EXISTS roster_composite_index;

-- GORM auto-migration will create the new index with (profile_id, contact_id, name)
-- Set default name for existing entries
UPDATE rosters SET name = 'default' WHERE name IS NULL OR name = '';
