--- Copyright 2023-2026 Ant Investor Ltd
---
--- Licensed under the Apache License, Version 2.0 (the "License");
--- you may not use this file except in compliance with the License.
--- You may obtain a copy of the License at
---
---      http://www.apache.org/licenses/LICENSE-2.0
---
--- Unless required by applicable law or agreed to in writing, software
--- distributed under the License is distributed on an "AS IS" BASIS,
--- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
--- See the License for the specific language governing permissions and
--- limitations under the License.

-- Every profile gets a stable handle for tagging and selecting it. New
-- profiles get one when created (business/handle.go); this gives existing
-- profiles theirs, by the same rules:
--   * the "name" (else "au_name") property, lowercased, with every run of
--     characters that are not a-z or 0-9 turned into one ".", cut to 24
--     characters;
--   * the oldest profile with that base gets it bare; the others get
--     "<base>.<last 8 characters of the profile id>";
--   * a profile without a usable name gets "user.<last 8 of its id>".
-- Handles are unique across all tenants and never reused, so the index
-- includes deleted profiles.
WITH slugs AS (
    SELECT id,
           created_at,
           rtrim(left(trim(BOTH '.' FROM regexp_replace(lower(coalesce(properties ->> 'name', '')),
                                                       '[^a-z0-9]+', '.', 'g')), 24), '.') AS name_slug,
           rtrim(left(trim(BOTH '.' FROM regexp_replace(lower(coalesce(properties ->> 'au_name', '')),
                                                       '[^a-z0-9]+', '.', 'g')), 24), '.') AS au_name_slug,
           right(regexp_replace(lower(id), '[^a-z0-9]', '', 'g'), 8)                      AS suffix
    FROM profiles
    WHERE handle IS NULL OR handle = ''
),
bases AS (
    SELECT id,
           created_at,
           suffix,
           CASE
               WHEN length(name_slug) >= 3 THEN name_slug
               WHEN length(au_name_slug) >= 3 THEN au_name_slug
               ELSE 'user'
           END AS base
    FROM slugs
),
ranked AS (
    SELECT id,
           base,
           suffix,
           row_number() OVER (PARTITION BY base ORDER BY created_at, id) AS rn
    FROM bases
)
UPDATE profiles p
SET handle = CASE
                 WHEN r.base <> 'user' AND r.rn = 1
                     AND NOT EXISTS (SELECT 1 FROM profiles q WHERE q.handle = r.base)
                     THEN r.base
                 ELSE r.base || '.' || r.suffix
             END
FROM ranked r
WHERE p.id = r.id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_profiles_handle
    ON profiles (handle)
    WHERE handle <> '';
