-- Additional indexes for geo_events queries.
-- The table is defined by GORM auto-migrate; these are supplemental spatial/temporal indexes.

-- Composite index for querying events by subject within a time range.
CREATE INDEX IF NOT EXISTS idx_geo_events_subject_ts
    ON geo_events (subject_id, ts DESC);

-- Composite index for querying events by area within a time range.
CREATE INDEX IF NOT EXISTS idx_geo_events_area_ts
    ON geo_events (area_id, ts DESC);

-- Index for event type filtering.
CREATE INDEX IF NOT EXISTS idx_geo_events_type
    ON geo_events (event_type);

-- Partial index for HasDwellEvent check: (subject_id, area_id, ts) WHERE event_type = DWELL(2).
-- Supports the geofence engine's dedup check for dwell events.
CREATE INDEX IF NOT EXISTS idx_geo_events_dwell_dedup
    ON geo_events (subject_id, area_id, ts DESC)
    WHERE event_type = 2;
