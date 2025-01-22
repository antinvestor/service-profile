
CREATE INDEX roster_search_idx ON rosters
    USING bm25 (id, properties, profile_id, contact_id, created_at)
    WITH (key_field='id');
