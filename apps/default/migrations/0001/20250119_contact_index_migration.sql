

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Note: We don't create index on encrypted_detail as it's encrypted data
-- The searchable column will handle text search functionality

-- Recreate with 'simple' configuration and handle empty jsonb_to_tsv properly
ALTER TABLE contacts
    ADD COLUMN searchable tsvector GENERATED ALWAYS AS (
        jsonb_to_tsv(COALESCE(properties, '{}'::jsonb))
        ) STORED;

-- Recreate the GIN index
CREATE INDEX idx_contacts_searchable ON contacts USING GIN (searchable);


