-- Enable PostGIS extension for spatial operations.
-- This must be run before any spatial columns or functions are used.
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
