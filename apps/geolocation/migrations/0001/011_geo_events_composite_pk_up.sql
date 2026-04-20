-- Copyright 2023-2026 Ant Investor Ltd
--
-- Licensed under the Apache License, Version 2.0 (the "License").

-- geo_events is promoted to a TimescaleDB hypertable. TimescaleDB requires
-- the time-partition column to participate in every UNIQUE/PRIMARY
-- constraint, so replace the BaseModel-default PK (id) with a composite
-- (id, true_created_at). Client-captured event timestamps are the
-- correctness-critical ordering dimension for offline-batched uploads.

ALTER TABLE geo_events DROP CONSTRAINT IF EXISTS geo_events_pkey;
ALTER TABLE geo_events ADD PRIMARY KEY (id, true_created_at);
