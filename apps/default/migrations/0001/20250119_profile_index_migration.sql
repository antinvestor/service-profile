ALTER TABLE profiles
    ADD COLUMN search_column tsvector GENERATED ALWAYS AS ( jsonb_to_tsv(properties) ) STORED;

CREATE INDEX idx_profiles_search_column ON profiles USING GIN (search_column);