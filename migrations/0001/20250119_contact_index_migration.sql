
CREATE INDEX contact_search_idx ON contacts
    USING bm25 (id, detail, properties, profile_id, created_at)
    WITH (key_field='id', text_fields='{"detail": {"tokenizer": {"type": "ngram", "min_gram": 3, "max_gram": 3, "prefix_only": false}}}');
