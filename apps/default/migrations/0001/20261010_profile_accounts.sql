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

-- A profile has at most one primary account per family and account version;
-- accounts moved in by a merge are secondary and not constrained.
CREATE UNIQUE INDEX IF NOT EXISTS idx_profile_accounts_primary
    ON profile_accounts (profile_id, family, account_version)
    WHERE is_primary AND deleted_at IS NULL;
