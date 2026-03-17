-- PostGIS trigger and indexes for routes, plus route_deviation_states table.
-- Route, RouteAssignment, and RouteDeviationState tables are created by
-- GORM auto-migrate from the Go models. Geom, length_m, and deviation
-- columns are in the Go model. This migration adds the trigger + indexes.

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

-- Indexes for route_assignments.
CREATE INDEX IF NOT EXISTS idx_ra_subject_state
    ON route_assignments (subject_id, state);

CREATE INDEX IF NOT EXISTS idx_ra_route
    ON route_assignments (route_id);

-- Index for route_deviation_states.
CREATE INDEX IF NOT EXISTS idx_rds_route
    ON route_deviation_states (route_id);
