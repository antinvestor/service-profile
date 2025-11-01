

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_contacts_detail_trgm ON contacts USING GIN (detail gin_trgm_ops);

-- Recreate with 'simple' configuration and handle empty jsonb_to_tsv properly
ALTER TABLE contacts
    ADD COLUMN searchable tsvector GENERATED ALWAYS AS (
        jsonb_to_tsv(COALESCE(properties, '{}'::jsonb))
        ) STORED;

-- Recreate the GIN index
CREATE INDEX idx_contacts_searchable ON contacts USING GIN (searchable);


