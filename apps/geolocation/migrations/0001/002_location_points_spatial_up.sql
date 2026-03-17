-- PostGIS trigger and spatial indexes for location_points.
-- The geom column is created by GORM auto-migrate from the Go model.
-- This migration adds the trigger that computes geom from lat/lon on INSERT.

-- Trigger function: compute geom from latitude/longitude on INSERT.
-- Must be VOLATILE (not IMMUTABLE) because it reads from NEW row fields.
CREATE OR REPLACE FUNCTION compute_location_point_geom()
RETURNS TRIGGER AS $$
BEGIN
    NEW.geom := ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_location_point_geom ON location_points;
CREATE TRIGGER trg_location_point_geom
    BEFORE INSERT ON location_points
    FOR EACH ROW
    EXECUTE FUNCTION compute_location_point_geom();

-- Spatial index for point-in-polygon and proximity queries.
CREATE INDEX IF NOT EXISTS idx_location_points_geom
    ON location_points USING GIST (geom);

-- Composite index for time-series queries by subject.
CREATE INDEX IF NOT EXISTS idx_lp_subject_ts_desc
    ON location_points (subject_id, ts DESC);

-- Index for ingestion time queries (used by data retention).
CREATE INDEX IF NOT EXISTS idx_lp_ingested_at
    ON location_points (ingested_at);
