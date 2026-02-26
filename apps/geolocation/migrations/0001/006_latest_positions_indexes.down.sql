-- Rollback: remove spatial column, trigger, and indexes from latest_positions.
DROP TRIGGER IF EXISTS trg_latest_position_geom ON latest_positions;
DROP FUNCTION IF EXISTS compute_latest_position_geom();

DROP INDEX IF EXISTS idx_latest_positions_ts;
DROP INDEX IF EXISTS idx_latest_positions_geom;

ALTER TABLE latest_positions DROP COLUMN IF EXISTS geom;
