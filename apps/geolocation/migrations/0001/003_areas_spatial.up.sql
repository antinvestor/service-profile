-- Add PostGIS geometry and bounding box columns to areas.
-- geom: the actual polygon/multipolygon geometry for containment testing.
-- bbox: the bounding box envelope for fast pre-filtering via GIST index.

ALTER TABLE areas
    ADD COLUMN IF NOT EXISTS geom geometry(Geometry, 4326),
    ADD COLUMN IF NOT EXISTS bbox geometry(Polygon, 4326);

-- Trigger function: compute bbox and metrics when geom is updated.
-- Computes: bounding box envelope, area in square meters, perimeter in meters.
-- Uses geography cast for accurate metric computation on the WGS 84 ellipsoid.
-- Must be VOLATILE (not IMMUTABLE) because it reads from NEW row fields.
CREATE OR REPLACE FUNCTION compute_area_metrics()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.geom IS NOT NULL THEN
        -- Compute bounding box.
        NEW.bbox := ST_Envelope(NEW.geom);

        -- Compute area in square meters using geography cast.
        NEW.area_m2 := ST_Area(NEW.geom::geography);

        -- Compute perimeter in meters using geography cast.
        NEW.perimeter_m := ST_Perimeter(NEW.geom::geography);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_area_metrics ON areas;
CREATE TRIGGER trg_area_metrics
    BEFORE INSERT OR UPDATE OF geom ON areas
    FOR EACH ROW
    EXECUTE FUNCTION compute_area_metrics();

-- GIST index on geom for containment (ST_Contains) queries.
CREATE INDEX IF NOT EXISTS idx_areas_geom
    ON areas USING GIST (geom);

-- GIST index on bbox for fast bounding box intersection pre-filtering.
CREATE INDEX IF NOT EXISTS idx_areas_bbox
    ON areas USING GIST (bbox);

-- Index for owner-based area lookups.
CREATE INDEX IF NOT EXISTS idx_areas_owner_state
    ON areas (owner_id, state);

-- Partial index for active areas only (state=2 is ACTIVE).
CREATE INDEX IF NOT EXISTS idx_areas_active
    ON areas (state) WHERE state = 2;
