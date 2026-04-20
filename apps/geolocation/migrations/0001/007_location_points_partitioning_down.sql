-- Rollback: convert partitioned location_points back to a regular table.
-- WARNING: This will rewrite the entire table. Run during a maintenance window.

-- Drop the helper function.
DROP FUNCTION IF EXISTS create_location_points_partitions(INTEGER);

-- Rename partitioned table.
ALTER TABLE IF EXISTS location_points RENAME TO location_points_partitioned;

-- Create regular (non-partitioned) table.
CREATE TABLE IF NOT EXISTS location_points (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50),
    partition_id VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    subject_id VARCHAR(40) NOT NULL,
    device_id VARCHAR(80) NOT NULL,
    true_created_at TIMESTAMPTZ NOT NULL,
    ingested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    altitude DOUBLE PRECISION,
    accuracy DOUBLE PRECISION NOT NULL DEFAULT 0,
    speed DOUBLE PRECISION,
    bearing DOUBLE PRECISION,
    source SMALLINT NOT NULL DEFAULT 0,
    extras JSONB DEFAULT '{}',
    processing_state SMALLINT NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ,
    processing_error TEXT NOT NULL DEFAULT '',
    geom geometry(Point, 4326)
);

-- Copy data from partitioned table.
INSERT INTO location_points
    SELECT * FROM location_points_partitioned;

-- Drop partitioned table (cascades to all partitions).
DROP TABLE IF EXISTS location_points_partitioned CASCADE;

-- Recreate indexes.
CREATE INDEX IF NOT EXISTS idx_location_points_geom
    ON location_points USING GIST (geom);

CREATE INDEX IF NOT EXISTS idx_lp_subject_true_created_at_desc
    ON location_points (subject_id, true_created_at DESC);

CREATE INDEX IF NOT EXISTS idx_lp_device_true_created_at_desc
    ON location_points (device_id, true_created_at DESC);

CREATE INDEX IF NOT EXISTS idx_lp_ingested_at
    ON location_points (ingested_at);

-- Recreate trigger.
DROP TRIGGER IF EXISTS trg_location_point_geom ON location_points;
CREATE TRIGGER trg_location_point_geom
    BEFORE INSERT ON location_points
    FOR EACH ROW
    EXECUTE FUNCTION compute_location_point_geom();
