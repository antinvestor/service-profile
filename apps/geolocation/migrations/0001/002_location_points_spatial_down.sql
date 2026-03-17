-- Rollback: remove spatial column, trigger, and indexes from location_points.
DROP TRIGGER IF EXISTS trg_location_point_geom ON location_points;
DROP FUNCTION IF EXISTS compute_location_point_geom();

DROP INDEX IF EXISTS idx_lp_ingested_at;
DROP INDEX IF EXISTS idx_lp_subject_ts_desc;
DROP INDEX IF EXISTS idx_location_points_geom;

ALTER TABLE location_points DROP COLUMN IF EXISTS geom;
