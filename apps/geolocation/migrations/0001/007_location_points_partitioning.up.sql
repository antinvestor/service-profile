-- Partition location_points by ts (timestamp) using PostgreSQL declarative range partitioning.
-- This allows efficient pruning of old data and keeps per-partition index sizes manageable.
--
-- Strategy: monthly partitions. The application (or a cron job) is responsible for creating
-- future partitions before they are needed. A DEFAULT partition catches any rows that don't
-- match an existing partition to prevent INSERT failures.
--
-- IMPORTANT: This migration converts the existing table to a partitioned table. This requires
-- PostgreSQL 11+ and will rewrite the table. Run during a maintenance window.
--
-- Migration approach:
--   1. Rename existing table to a temporary name.
--   2. Create the new partitioned table with the same schema.
--   3. Create an initial set of monthly partitions.
--   4. Copy data from the old table.
--   5. Drop the old table.

-- Step 1: Rename existing table.
ALTER TABLE IF EXISTS location_points RENAME TO location_points_old;

-- Step 2: Create partitioned table with same schema.
CREATE TABLE IF NOT EXISTS location_points (
    id VARCHAR(50) NOT NULL,
    tenant_id VARCHAR(50),
    partition_id VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    subject_id VARCHAR(40) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    ingested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    altitude DOUBLE PRECISION,
    accuracy DOUBLE PRECISION NOT NULL DEFAULT 0,
    speed DOUBLE PRECISION,
    bearing DOUBLE PRECISION,
    source SMALLINT NOT NULL DEFAULT 0,
    extras JSONB DEFAULT '{}',
    geom geometry(Point, 4326),
    PRIMARY KEY (id, ts)
) PARTITION BY RANGE (ts);

-- Step 3: Create a DEFAULT partition for safety (catches rows outside defined ranges).
CREATE TABLE IF NOT EXISTS location_points_default
    PARTITION OF location_points DEFAULT;

-- Step 4: Create monthly partitions for the current and next 3 months.
-- Additional partitions should be created by a scheduled job.
DO $$
DECLARE
    start_date DATE;
    end_date DATE;
    partition_name TEXT;
    i INTEGER;
BEGIN
    -- Create partitions for last month through next 3 months (5 total).
    FOR i IN -1..3 LOOP
        start_date := date_trunc('month', NOW()) + (i || ' months')::interval;
        end_date := start_date + '1 month'::interval;
        partition_name := 'location_points_' || to_char(start_date, 'YYYY_MM');

        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF location_points
             FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
    END LOOP;
END $$;

-- Step 5: Copy data from old table if it exists and has rows.
INSERT INTO location_points
    SELECT * FROM location_points_old
    WHERE EXISTS (SELECT 1 FROM location_points_old LIMIT 1);

-- Step 6: Drop old table.
DROP TABLE IF EXISTS location_points_old;

-- Step 7: Recreate indexes on the partitioned table.
-- These indexes propagate to all child partitions automatically.
CREATE INDEX IF NOT EXISTS idx_location_points_geom
    ON location_points USING GIST (geom);

CREATE INDEX IF NOT EXISTS idx_lp_subject_ts_desc
    ON location_points (subject_id, ts DESC);

CREATE INDEX IF NOT EXISTS idx_lp_ingested_at
    ON location_points (ingested_at);

-- Step 8: Recreate trigger for geom computation.
DROP TRIGGER IF EXISTS trg_location_point_geom ON location_points;
CREATE TRIGGER trg_location_point_geom
    BEFORE INSERT ON location_points
    FOR EACH ROW
    EXECUTE FUNCTION compute_location_point_geom();

-- Step 9: Helper function to create monthly partitions (called by a cron job or startup routine).
-- Usage: SELECT create_location_points_partitions(6) to create 6 months ahead.
CREATE OR REPLACE FUNCTION create_location_points_partitions(months_ahead INTEGER DEFAULT 3)
RETURNS VOID AS $$
DECLARE
    start_date DATE;
    end_date DATE;
    partition_name TEXT;
    i INTEGER;
BEGIN
    FOR i IN 0..months_ahead LOOP
        start_date := date_trunc('month', NOW()) + (i || ' months')::interval;
        end_date := start_date + '1 month'::interval;
        partition_name := 'location_points_' || to_char(start_date, 'YYYY_MM');

        BEGIN
            EXECUTE format(
                'CREATE TABLE IF NOT EXISTS %I PARTITION OF location_points
                 FOR VALUES FROM (%L) TO (%L)',
                partition_name, start_date, end_date
            );
        EXCEPTION WHEN duplicate_table THEN
            -- Partition already exists, skip.
            NULL;
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;
