-- Additional columns and indexes for latest_positions proximity queries.
-- The table is defined by GORM auto-migrate with PK on subject_id.

-- Add PostGIS POINT geometry column for efficient spatial proximity queries.
ALTER TABLE latest_positions
    ADD COLUMN IF NOT EXISTS geom geometry(Point, 4326);

-- Trigger function: compute geom from latitude/longitude on INSERT/UPDATE.
CREATE OR REPLACE FUNCTION compute_latest_position_geom()
RETURNS TRIGGER AS $$
BEGIN
    NEW.geom := ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_latest_position_geom ON latest_positions;
CREATE TRIGGER trg_latest_position_geom
    BEFORE INSERT OR UPDATE OF latitude, longitude ON latest_positions
    FOR EACH ROW
    EXECUTE FUNCTION compute_latest_position_geom();

-- Backfill existing rows that have lat/lon but no geom.
UPDATE latest_positions
    SET geom = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
    WHERE geom IS NULL AND latitude != 0 AND longitude != 0;

-- GIST index for PostGIS ST_DWithin proximity queries.
CREATE INDEX IF NOT EXISTS idx_latest_positions_geom
    ON latest_positions USING GIST (geom);

-- Index on timestamp for staleness filtering in proximity queries.
CREATE INDEX IF NOT EXISTS idx_latest_positions_ts
    ON latest_positions (ts DESC);
