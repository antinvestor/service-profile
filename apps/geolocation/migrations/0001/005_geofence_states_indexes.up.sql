-- Additional indexes for geofence_states lookups.
-- The table is defined by GORM auto-migrate with composite PK (subject_id, area_id).

-- Index for "who is inside this area?" queries.
CREATE INDEX IF NOT EXISTS idx_geofence_states_area_inside
    ON geofence_states (area_id) WHERE inside = true;

-- Index for "what areas is this subject inside?" queries.
CREATE INDEX IF NOT EXISTS idx_geofence_states_subject_inside
    ON geofence_states (subject_id) WHERE inside = true;
