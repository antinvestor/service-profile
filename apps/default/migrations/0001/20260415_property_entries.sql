-- Composite indexes for property entry queries
CREATE INDEX IF NOT EXISTS idx_prop_profile_created
    ON property_entries (profile_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_prop_profile_key_created
    ON property_entries (profile_id, key, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_prop_tenant_scoped_lookup
    ON property_entries (profile_id, tenant_id, scoped)
    WHERE scoped = TRUE;

-- Explode existing Profile.Properties JSONB into property_entries
INSERT INTO property_entries (id, created_at, modified_at, created_by, modified_by, version, tenant_id, partition_id, access_id, profile_id, key, value, scoped)
SELECT
    gen_random_uuid()::varchar(50),
    p.created_at,
    p.created_at,
    COALESCE(p.created_by, ''),
    COALESCE(p.created_by, ''),
    0,
    COALESCE(p.tenant_id, ''),
    COALESCE(p.partition_id, ''),
    COALESCE(p.access_id, ''),
    p.id,
    kv.key,
    kv.value::text,
    FALSE
FROM profiles p,
     LATERAL jsonb_each_text(p.properties) AS kv(key, value)
WHERE p.properties IS NOT NULL
  AND p.properties != '{}'::jsonb
  AND p.deleted_at IS NULL;
