DROP INDEX IF EXISTS idx_rds_route;
DROP INDEX IF EXISTS idx_ra_route;
DROP INDEX IF EXISTS idx_ra_subject_state;
DROP INDEX IF EXISTS idx_routes_active_deviation;
DROP INDEX IF EXISTS idx_routes_owner_state;
DROP INDEX IF EXISTS idx_routes_geom;

DROP TRIGGER IF EXISTS trg_route_length ON routes;
DROP FUNCTION IF EXISTS compute_route_length();
