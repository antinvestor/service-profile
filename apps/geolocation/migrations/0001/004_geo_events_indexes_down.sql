-- Rollback: remove supplemental indexes from geo_events.
DROP INDEX IF EXISTS idx_geo_events_dwell_dedup;
DROP INDEX IF EXISTS idx_geo_events_type;
DROP INDEX IF EXISTS idx_geo_events_area_ts;
DROP INDEX IF EXISTS idx_geo_events_subject_ts;
