ALTER TABLE contacts
    ADD COLUMN search_column tsvector GENERATED ALWAYS AS ( jsonb_to_tsv(properties) || to_tsvector('english', coalesce(detail, '')) ) STORED;

CREATE INDEX idx_contacts_search_column ON contacts USING GIN (search_column);