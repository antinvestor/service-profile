-- Rollback: remove spatial columns, trigger, and indexes from areas.
DROP TRIGGER IF EXISTS trg_area_metrics ON areas;
DROP FUNCTION IF EXISTS compute_area_metrics();

DROP INDEX IF EXISTS idx_areas_active;
DROP INDEX IF EXISTS idx_areas_owner_state;
DROP INDEX IF EXISTS idx_areas_bbox;
DROP INDEX IF EXISTS idx_areas_geom;

ALTER TABLE areas
    DROP COLUMN IF EXISTS bbox,
    DROP COLUMN IF EXISTS geom;
