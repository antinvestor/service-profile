ALTER TABLE IF EXISTS location_points
    ADD COLUMN IF NOT EXISTS device_id VARCHAR(80);

UPDATE location_points
SET device_id = COALESCE(NULLIF(access_id, ''), subject_id || '-unknown-device')
WHERE device_id IS NULL OR device_id = '';

ALTER TABLE IF EXISTS location_points
    ALTER COLUMN device_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_lp_device_true_created_at_desc
    ON location_points (device_id, true_created_at DESC);

ALTER TABLE IF EXISTS latest_positions
    ADD COLUMN IF NOT EXISTS device_id VARCHAR(80);

UPDATE latest_positions lp
SET device_id = source_points.device_id
FROM (
    SELECT DISTINCT ON (tenant_id, partition_id, subject_id)
        tenant_id,
        partition_id,
        subject_id,
        device_id,
        true_created_at
    FROM location_points
    WHERE deleted_at IS NULL
      AND device_id IS NOT NULL
      AND device_id <> ''
    ORDER BY tenant_id, partition_id, subject_id, true_created_at DESC, modified_at DESC
) AS source_points
WHERE lp.subject_id = source_points.subject_id
  AND COALESCE(lp.tenant_id, '') = COALESCE(source_points.tenant_id, '')
  AND COALESCE(lp.partition_id, '') = COALESCE(source_points.partition_id, '')
  AND (lp.device_id IS NULL OR lp.device_id = '');

UPDATE latest_positions
SET device_id = subject_id || '-unknown-device'
WHERE device_id IS NULL OR device_id = '';

ALTER TABLE IF EXISTS latest_positions
    ALTER COLUMN device_id SET NOT NULL;
