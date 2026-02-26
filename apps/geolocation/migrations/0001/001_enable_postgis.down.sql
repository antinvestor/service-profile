-- Rollback: drop PostGIS extensions.
-- WARNING: This will drop ALL spatial columns and indexes in the database.
-- Only run if no other tables/schemas depend on PostGIS.
DROP EXTENSION IF EXISTS postgis_topology CASCADE;
DROP EXTENSION IF EXISTS postgis CASCADE;
