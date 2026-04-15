# Profile Service: Cross-Tenant Architecture, Property Ledger, and Multi-List Roster

**Date:** 2026-04-15
**Status:** Approved design, pending implementation

## Problem Statement

The profile service is a shared platform service consumed by many services across tenants. Three issues require architectural changes:

1. **Tenancy scoping mismatch:** Frame's `TenancyPartition` GORM scope filters reads by the caller's tenant/partition. Profile data (profiles, contacts, addresses) and global seed data (profile_types, relationship_types, countries) must be readable cross-tenant, but the current repository layer inconsistently handles this — some methods use raw ctx (broken), some use empty claims (broken for NULL tenant rows), some use `SkipTenancyChecksOnClaims` (correct but scattered).

2. **No property change tracking:** `Profile.Properties` is a mutable JSONB blob. When multiple tenants/services write properties on the same profile, there is no provenance, no diff, and no way to see who changed what. Tenants also need private properties on shared profiles that other tenants cannot see.

3. **Flat roster model:** Rosters link a profile to contacts without categorization. Users need multiple named lists (friends, colleagues, vip) with the same contact appearing in multiple lists. Rosters must remain strictly tenant-scoped.

## Design

### 1. Cross-Tenant Repository Layer

**Rule:** Data about identity is shared. Data about how a tenant organizes identity is private.

| Data | Reads | Writes |
|---|---|---|
| Profiles | Cross-tenant (`SkipTenancyChecksOnClaims`) | Stamped with creator's tenant (raw ctx) |
| Contacts | Cross-tenant | Stamped with creator's tenant |
| Addresses | Cross-tenant | Stamped with creator's tenant |
| Profile types, relationship types, countries | Cross-tenant (global seed data) | Seed migrations only |
| Rosters | Tenant-scoped (raw ctx) | Tenant-scoped (raw ctx) |
| Relationships | Tenant-scoped (raw ctx) | Tenant-scoped (raw ctx) |
| Property entries (global) | Cross-tenant | Stamped with creator's tenant |
| Property entries (scoped) | Tenant-scoped (caller's tenant only) | Stamped with creator's tenant |

**Implementation pattern:** Every repository read method for cross-tenant data applies `SkipTenancyChecksOnClaims(ctx)` before querying. Every write method uses raw `ctx` so `BaseModel.GenID()` stamps provenance. No empty claims, no nil claims, no special cases.

### 2. Hybrid Append-Only Property Log

#### 2.1 New Table: `property_entries`

```sql
CREATE TABLE property_entries (
    id              VARCHAR(50) PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL,
    modified_at     TIMESTAMPTZ NOT NULL,
    created_by      VARCHAR(50) NOT NULL,
    modified_by     VARCHAR(50) NOT NULL,
    version         BIGINT DEFAULT 0,
    tenant_id       VARCHAR(50) NOT NULL,
    partition_id    VARCHAR(50) NOT NULL,
    access_id       VARCHAR(50) NOT NULL,
    deleted_at      TIMESTAMPTZ,

    profile_id      VARCHAR(50) NOT NULL,
    key             VARCHAR(255) NOT NULL,
    value           TEXT NOT NULL,
    scoped          BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_property_entries_profile ON property_entries (profile_id, created_at DESC);
CREATE INDEX idx_property_entries_profile_key ON property_entries (profile_id, key, created_at DESC);
CREATE INDEX idx_property_entries_tenant ON property_entries (profile_id, tenant_id, scoped) WHERE scoped = TRUE;
```

- `scoped = false`: Global property. Visible to all callers. Included in `Profile.Properties` JSONB cache.
- `scoped = true`: Tenant-private property. Visible only when caller's tenant matches `tenant_id`. Excluded from JSONB cache.

Entries are immutable once written. The latest entry per (profile_id, key, scoped, tenant_id) determines current state.

#### 2.2 Write Flow

1. Caller invokes `UpdateProperties` or `UpdateTenantProperties`.
2. Authorization check: caller must have `profile_update` permission (ReBAC — profile owner, partition owner/admin, or authorized agent).
3. For each key-value pair, append a row to `property_entries` with provenance (tenant_id, partition_id, access_id, created_by from claims).
4. For global properties (`scoped = false`): rebuild `Profile.Properties` JSONB cache from latest global entries per key. Save the profile.
5. For tenant properties (`scoped = true`): no cache rebuild.

#### 2.3 Read Flow

| Scenario | Behavior |
|---|---|
| `GetByID(profileID)` without tenant claims | Return `Profile.Properties` (global cache only) |
| `GetByID(profileID)` with tenant claims | Return `Profile.Properties` merged with caller's tenant-scoped entries |
| `GetPropertiesByTenant(profileID, tenantID)` | Return only properties written by that tenant (global + scoped) |
| `GetPropertyHistory(profileID, key)` | Full change log for that key. Scoped entries filtered to caller's tenant. |

#### 2.4 Proto Changes

```protobuf
message ProfileObject {
    // ... existing fields ...
    google.protobuf.Struct properties = 4;         // Global properties (existing)
    google.protobuf.Struct tenant_properties = 10;  // Caller's tenant-scoped properties
}

message UpdatePropertiesRequest {
    string id = 1;
    google.protobuf.Struct properties = 2;
}

message UpdateTenantPropertiesRequest {
    string id = 1;
    google.protobuf.Struct properties = 2;
}

message PropertyHistoryRequest {
    string id = 1;
    string key = 2;
}

message PropertyEntry {
    string key = 1;
    string value = 2;
    string tenant_id = 3;
    string created_by = 4;
    google.protobuf.Timestamp created_at = 5;
    bool scoped = 6;
}

message PropertyHistoryResponse {
    repeated PropertyEntry entries = 1;
}
```

#### 2.5 Migration

Existing `Profile.Properties` JSONB data is exploded into `property_entries` rows:
- One row per key-value pair
- `scoped = false` (all existing data is global)
- `tenant_id`, `partition_id`, `access_id` copied from the profile's own fields (creator provenance)
- `created_by` set to the profile's `created_by`
- `created_at` set to the profile's `created_at`

The `Profile.Properties` column is retained as the denormalized cache.

### 3. Multi-List Roster

#### 3.1 Schema Change

Add `list_name` column to `rosters` table:

```sql
ALTER TABLE rosters ADD COLUMN list_name VARCHAR(255) NOT NULL DEFAULT 'default';

-- Replace existing unique index
DROP INDEX IF EXISTS roster_composite_index;
CREATE UNIQUE INDEX roster_composite_index ON rosters (profile_id, contact_id, list_name, tenant_id);
```

#### 3.2 Behavior

- List names are freeform strings chosen by the user.
- A contact can appear in multiple lists within the same tenant.
- Roster reads and writes are always tenant-scoped (raw ctx, `TenancyPartition` enforced).
- Different tenants have completely isolated roster lists on the same profile.

#### 3.3 API Changes

```protobuf
message RosterEntry {
    // ... existing fields ...
    string list_name = 6;
}

message ProcessRosterBatchRequest {
    string profile_id = 1;
    string list_name = 2;       // Required, e.g., "friends"
    repeated RosterItem items = 3;
}

message SearchRosterRequest {
    string profile_id = 1;
    string list_name = 2;       // Empty = all lists for this tenant
    int32 count = 3;
    int32 page = 4;
}

message RemoveRosterRequest {
    string profile_id = 1;
    string contact_id = 2;
    string list_name = 3;       // Required
}
```

#### 3.4 Migration

Existing roster entries get `list_name = 'default'`.

### 4. Authorization

All profile mutation paths (property writes, contact changes, address changes) are gated by `profile_update` permission via the existing ReBAC system:

| Actor | Can Edit |
|---|---|
| Profile owner | Their own profile |
| Partition owner/admin | Any profile in their partition |
| Authorized agent (service account) | Profiles they are explicitly granted access to |

No new permissions are introduced. The property entries system, tenant properties, and existing profile update handlers all use the same `profile_update` check.

Roster mutations use existing roster permissions and are additionally scoped by tenancy at the query level.

### 5. Testing Strategy

All new tests must include variants **with tenant-scoped auth claims** to match production behavior. Tests without claims only verify the nil-claims bypass path, which is not the production path.

Key test scenarios:
- Create profile with all profile types under tenant context
- Read profile cross-tenant (created by tenant A, read by tenant B)
- Write global property, verify JSONB cache rebuilt
- Write tenant-scoped property, verify not in JSONB cache
- Read profile with tenant context, verify tenant properties merged
- Read profile without tenant context, verify only global properties
- Property history returns full chain
- Property history filters scoped entries by caller's tenant
- Roster CRUD with named lists under tenant context
- Roster isolation: tenant A's lists invisible to tenant B
- Same contact in multiple lists within same tenant
- Authorization: owner can update, random profile cannot, partition admin can
