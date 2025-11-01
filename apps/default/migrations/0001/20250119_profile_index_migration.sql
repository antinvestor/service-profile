-- Recreate with 'simple' configuration and handle empty jsonb_to_tsv properly
ALTER TABLE profiles
    ADD COLUMN searchable tsvector GENERATED ALWAYS AS (
        jsonb_to_tsv(COALESCE(properties, '{}'::jsonb))
        ) STORED;

-- Recreate the GIN index
CREATE INDEX idx_profiles_searchable ON profiles USING GIN (searchable);
