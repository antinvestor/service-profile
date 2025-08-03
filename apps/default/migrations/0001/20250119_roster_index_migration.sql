
ALTER TABLE rosters
    ADD COLUMN search_properties tsvector GENERATED ALWAYS AS ( jsonb_to_tsv(properties) ) STORED;

CREATE INDEX idx_rosters_search_properties ON rosters USING GIN (search_properties);
