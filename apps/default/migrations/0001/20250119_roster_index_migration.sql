
ALTER TABLE rosters
    ADD COLUMN search_column tsvector GENERATED ALWAYS AS ( jsonb_to_tsv(properties) ) STORED;

CREATE INDEX idx_rosters_search_column ON rosters USING GIN (search_column);
