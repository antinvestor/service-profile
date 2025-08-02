
CREATE OR REPLACE FUNCTION jsonb_to_tsv(jdoc jsonb)
    RETURNS tsvector AS $$
WITH RECURSIVE all_text AS (
    -- First level: key/value pairs
    SELECT value
    FROM jsonb_each_text(jdoc)

    UNION ALL

    -- If nested JSON objects/arrays exist, go deeper
    SELECT v.value
    FROM (
             SELECT value::jsonb AS val
             FROM jsonb_each(jdoc)
             WHERE jsonb_typeof(value) IN ('object', 'array')
         ) nested
             CROSS JOIN LATERAL jsonb_each_text(nested.val) v
)
SELECT
    -- Combine English + Simple configs for safety
    to_tsvector('english', string_agg(value, ' ')) ||
    to_tsvector('simple',  string_agg(value, ' '))
FROM all_text;
$$ LANGUAGE SQL IMMUTABLE;