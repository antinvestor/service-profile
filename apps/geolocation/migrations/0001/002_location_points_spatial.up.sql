-- Add PostGIS POINT geometry column to location_points.
-- SRID 4326 = WGS 84 (standard GPS coordinate system).
-- This column is populated by a trigger on INSERT that reads latitude/longitude.

ALTER TABLE location_points
    ADD COLUMN IF NOT EXISTS geom geometry(Point, 4326);

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
