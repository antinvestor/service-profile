
-- Recreate with 'simple' configuration and handle empty jsonb_to_tsv properly
ALTER TABLE rosters
    ADD COLUMN searchable tsvector GENERATED ALWAYS AS (
        jsonb_to_tsv(COALESCE(properties, '{}'::jsonb)) || ' ' ||
        to_tsvector('english', COALESCE(description, ''))
        ) STORED;

-- Recreate the GIN index
CREATE INDEX idx_rosters_searchable ON rosters USING GIN (searchable);
