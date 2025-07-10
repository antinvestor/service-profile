CREATE INDEX profile_search_idx ON profiles
    USING bm25 (id, properties, created_at)
    WITH (
        key_field='id',
        json_fields = '{ "properties": {"fast": true}}'
    );
