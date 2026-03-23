ALTER TABLE IF EXISTS location_points
    ADD COLUMN IF NOT EXISTS processing_state SMALLINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS processed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS processing_error TEXT NOT NULL DEFAULT '';

UPDATE location_points
SET processing_state = 1,
    processed_at = COALESCE(processed_at, modified_at, ingested_at),
    processing_error = ''
WHERE processing_state = 0
  AND processed_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_lp_processing_state
    ON location_points (processing_state, ingested_at);

ALTER TABLE IF EXISTS geofence_states
    ADD COLUMN IF NOT EXISTS id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS partition_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS access_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

UPDATE geofence_states gs
SET tenant_id = a.tenant_id,
    partition_id = a.partition_id,
    access_id = a.access_id,
    created_at = COALESCE(gs.created_at, gs.enter_ts, gs.last_transition, gs.last_point_ts, NOW()),
    modified_at = COALESCE(gs.modified_at, gs.last_point_ts, gs.last_transition, NOW()),
    version = GREATEST(gs.version, 1)
FROM areas a
WHERE gs.area_id = a.id
  AND (gs.tenant_id IS NULL OR gs.partition_id IS NULL);

UPDATE geofence_states
SET id = SUBSTR(MD5(subject_id || ':' || area_id), 1, 24)
WHERE id IS NULL OR id = '';

DELETE FROM geofence_states
WHERE tenant_id IS NULL OR partition_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_geofence_states_scope_pair
    ON geofence_states (tenant_id, partition_id, subject_id, area_id);

ALTER TABLE IF EXISTS latest_positions
    ADD COLUMN IF NOT EXISTS id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS partition_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS access_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

TRUNCATE TABLE latest_positions;

CREATE UNIQUE INDEX IF NOT EXISTS idx_latest_positions_scope_subject
    ON latest_positions (tenant_id, partition_id, subject_id);

ALTER TABLE IF EXISTS route_deviation_states
    ADD COLUMN IF NOT EXISTS id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS partition_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS access_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

UPDATE route_deviation_states rds
SET tenant_id = r.tenant_id,
    partition_id = r.partition_id,
    access_id = r.access_id,
    created_at = COALESCE(rds.created_at, rds.last_deviation_event_at, rds.last_point_ts, NOW()),
    modified_at = COALESCE(rds.modified_at, rds.last_point_ts, rds.last_deviation_event_at, NOW()),
    version = GREATEST(rds.version, 1)
FROM routes r
WHERE rds.route_id = r.id
  AND (rds.tenant_id IS NULL OR rds.partition_id IS NULL);

UPDATE route_deviation_states
SET id = SUBSTR(MD5(subject_id || ':' || route_id), 1, 24)
WHERE id IS NULL OR id = '';

DELETE FROM route_deviation_states
WHERE tenant_id IS NULL OR partition_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_route_deviation_states_scope_pair
    ON route_deviation_states (tenant_id, partition_id, subject_id, route_id);
