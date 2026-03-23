DROP INDEX IF EXISTS idx_lp_device_ts_desc;

ALTER TABLE IF EXISTS latest_positions
    DROP COLUMN IF EXISTS device_id;

ALTER TABLE IF EXISTS location_points
    DROP COLUMN IF EXISTS device_id;
