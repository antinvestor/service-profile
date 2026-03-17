-- Rollback: remove supplemental indexes from geofence_states.
DROP INDEX IF EXISTS idx_geofence_states_subject_inside;
DROP INDEX IF EXISTS idx_geofence_states_area_inside;
