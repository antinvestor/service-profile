
-- Recreate with 'simple' configuration and handle empty jsonb_to_tsv properly
ALTER TABLE device_logs
    ADD COLUMN search_data tsvector GENERATED ALWAYS AS (
        jsonb_to_tsv(COALESCE(data, '{}'::jsonb))
        ) STORED;

-- Recreate the GIN index
CREATE INDEX idx_device_log_search_properties ON device_logs USING GIN (search_data);
