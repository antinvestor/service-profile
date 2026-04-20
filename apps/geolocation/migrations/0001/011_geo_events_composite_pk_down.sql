-- Copyright 2023-2026 Ant Investor Ltd

ALTER TABLE geo_events DROP CONSTRAINT IF EXISTS geo_events_pkey;
ALTER TABLE geo_events ADD PRIMARY KEY (id);
