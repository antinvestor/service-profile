DROP INDEX IF EXISTS idx_route_deviation_states_scope_pair;
DROP INDEX IF EXISTS idx_latest_positions_scope_subject;
DROP INDEX IF EXISTS idx_geofence_states_scope_pair;
DROP INDEX IF EXISTS idx_lp_processing_state;

ALTER TABLE IF EXISTS route_deviation_states
    DROP COLUMN IF EXISTS id,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS modified_at,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS partition_id,
    DROP COLUMN IF EXISTS access_id,
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE IF EXISTS latest_positions
    DROP COLUMN IF EXISTS id,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS modified_at,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS partition_id,
    DROP COLUMN IF EXISTS access_id,
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE IF EXISTS geofence_states
    DROP COLUMN IF EXISTS id,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS modified_at,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS partition_id,
    DROP COLUMN IF EXISTS access_id,
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE IF EXISTS location_points
    DROP COLUMN IF EXISTS processing_state,
    DROP COLUMN IF EXISTS processed_at,
    DROP COLUMN IF EXISTS processing_error;
