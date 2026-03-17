-- Routes: predefined paths (LineString) for route deviation detection.

-- 1. Routes table — base columns created by Frame auto-migrate from Go model.
-- Add PostGIS and deviation columns that Frame doesn't know about.
ALTER TABLE routes ADD COLUMN IF NOT EXISTS geometry_json TEXT;
ALTER TABLE routes ADD COLUMN IF NOT EXISTS geom geometry(LineString, 4326);
ALTER TABLE routes ADD COLUMN IF NOT EXISTS length_m DOUBLE PRECISION;
ALTER TABLE routes ADD COLUMN IF NOT EXISTS deviation_threshold_m DOUBLE PRECISION;
ALTER TABLE routes ADD COLUMN IF NOT EXISTS deviation_consecutive_count INTEGER;
ALTER TABLE routes ADD COLUMN IF NOT EXISTS deviation_cooldown_sec INTEGER;

-- Trigger function: compute route length when geom is updated.
CREATE OR REPLACE FUNCTION compute_route_length()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.geom IS NOT NULL THEN
        NEW.length_m := ST_Length(NEW.geom::geography);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_route_length ON routes;
CREATE TRIGGER trg_route_length
    BEFORE INSERT OR UPDATE OF geom ON routes
    FOR EACH ROW
    EXECUTE FUNCTION compute_route_length();

-- GIST index on route geometry for distance queries.
CREATE INDEX IF NOT EXISTS idx_routes_geom
    ON routes USING GIST (geom);

-- Index for owner-based route lookups.
CREATE INDEX IF NOT EXISTS idx_routes_owner_state
    ON routes (owner_id, state);

-- Partial index for active routes with deviation config.
CREATE INDEX IF NOT EXISTS idx_routes_active_deviation
    ON routes (state)
    WHERE state = 2 AND deviation_threshold_m IS NOT NULL;

-- 2. Route assignments table
CREATE TABLE IF NOT EXISTS route_assignments (
    id         VARCHAR(40)  PRIMARY KEY,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    subject_id VARCHAR(40)  NOT NULL,
    route_id   VARCHAR(40)  NOT NULL REFERENCES routes(id),
    valid_from TIMESTAMPTZ,
    valid_until TIMESTAMPTZ,
    state      SMALLINT     NOT NULL DEFAULT 0,
    extras     JSONB        DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_ra_subject_state
    ON route_assignments (subject_id, state);

CREATE INDEX IF NOT EXISTS idx_ra_route
    ON route_assignments (route_id);

-- 3. Route deviation states table (composite PK, no BaseModel)
CREATE TABLE IF NOT EXISTS route_deviation_states (
    subject_id              VARCHAR(40)      NOT NULL,
    route_id                VARCHAR(40)      NOT NULL,
    deviated                BOOLEAN          NOT NULL DEFAULT FALSE,
    consecutive_off_route   INTEGER          NOT NULL DEFAULT 0,
    last_deviation_event_at TIMESTAMPTZ,
    last_point_ts           TIMESTAMPTZ,
    last_lat                DOUBLE PRECISION NOT NULL DEFAULT 0,
    last_lon                DOUBLE PRECISION NOT NULL DEFAULT 0,
    updated_at              TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    PRIMARY KEY (subject_id, route_id)
);

CREATE INDEX IF NOT EXISTS idx_rds_route
    ON route_deviation_states (route_id);
