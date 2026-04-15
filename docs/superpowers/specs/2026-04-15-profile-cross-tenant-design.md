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

#### 2.1 New Model: `PropertyEntry`

Table created via GORM auto-migration from the Go model:

```go
type PropertyEntry struct {
    data.BaseModel
    ProfileID string `gorm:"type:varchar(50);not null;index:idx_prop_profile,priority:1"`
    Key       string `gorm:"type:varchar(255);not null;index:idx_prop_profile_key,priority:1"`
    Value     string `gorm:"type:text;not null"`
    Scoped    bool   `gorm:"not null;default:false;index:idx_prop_tenant_scoped"`
}
```

`data.BaseModel` provides `ID`, `CreatedAt`, `ModifiedAt`, `CreatedBy`, `ModifiedBy`, `Version`, `TenantID`, `PartitionID`, `AccessID`, `DeletedAt`.

GORM will create appropriate indexes. Additional composite indexes added via migration:
- `(profile_id, created_at DESC)` for current-state queries
- `(profile_id, key, created_at DESC)` for key history queries
- `(profile_id, tenant_id, scoped)` partial index where `scoped = true` for tenant property lookups

- `scoped = false`: Global property. Visible to all callers. Included in `Profile.Properties` JSONB cache.
- `scoped = true`: Tenant-private property. Visible only when caller's tenant matches `tenant_id`. Excluded from JSONB cache.

Entries are immutable once written. The latest entry per (profile_id, key, scoped, tenant_id) determines current state.

#### 2.2 Write Flow

1. Caller invokes `UpdateProperties` with a `scoped` flag indicating global or tenant-private.
2. Authorization check: caller must have `profile_update` permission (ReBAC — profile owner, partition owner/admin, or authorized agent).
3. For each key-value pair, append a row to `property_entries` with provenance (tenant_id, partition_id, access_id, created_by from claims) and the `scoped` flag.
4. For global properties (`scoped = false`): rebuild `Profile.Properties` JSONB cache from latest global entries per key. Save the profile.
5. For tenant properties (`scoped = true`): no cache rebuild.

#### 2.3 Read Flow

| Scenario | Behavior |
|---|---|
| `GetByID(profileID)` | Return `Profile.Properties` (global cache only) |
| `GetByIDAndPartition(profileID, partitionID)` | Return global properties + partition's tenant-scoped entries merged |
| `GetPropertiesByTenant(profileID, tenantID)` | Return only properties written by that tenant (global + scoped) |
| `GetPropertyHistory(profileID, key)` | Full change log for that key. Scoped entries filtered to caller's tenant. |

#### 2.4 Proto Changes

```protobuf
message ProfileObject {
    // ... existing fields ...
    google.protobuf.Struct properties = 4;  // Global properties (from JSONB cache)
}

// Existing UpdateRequest gains a scoped flag
message UpdateRequest {
    string id = 1;
    google.protobuf.Struct properties = 2;
    bool scoped = 3;  // false = global (default), true = tenant-private
}

// New: get profile with partition-scoped properties merged in
message GetByIDAndPartitionRequest {
    string id = 1;
    string partition_id = 2;
}
// Returns ProfileObject with properties = global + partition's scoped entries merged

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

#### 3.1 Model Change

Add `Name` field to the existing `Roster` GORM model:

```go
type Roster struct {
    data.BaseModel
    ProfileID  string          `gorm:"type:varchar(50);uniqueIndex:roster_composite_index,priority:1"`
    ContactID  string          `gorm:"type:varchar(50);uniqueIndex:roster_composite_index,priority:2"`
    Name   string          `gorm:"type:varchar(255);not null;default:'default';uniqueIndex:roster_composite_index,priority:3"`
    Contact    Contact         `gorm:"foreignKey:ContactID"`
    Properties data.JSONMap    `gorm:"type:JSONB"`
}
```

The unique index changes from `(profile_id, contact_id)` to `(profile_id, contact_id, name)`. GORM auto-migration handles the column addition. A SQL migration drops the old index and the new one is created by GORM. The tenant_id is already part of the tenancy scoping at query time, so it does not need to be in the unique index.

#### 3.2 Behavior

- List names are freeform strings chosen by the user.
- A contact can appear in multiple lists within the same tenant.
- Roster reads and writes are always tenant-scoped (raw ctx, `TenancyPartition` enforced).
- Different tenants have completely isolated roster lists on the same profile.

#### 3.3 API Changes

```protobuf
message RosterEntry {
    // ... existing fields ...
    string name = 6;
}

message ProcessRosterBatchRequest {
    string profile_id = 1;
    string name = 2;       // Required, e.g., "friends"
    repeated RosterItem items = 3;
}

message SearchRosterRequest {
    string profile_id = 1;
    string name = 2;       // Empty = all lists for this tenant
    int32 count = 3;
    int32 page = 4;
}

message RemoveRosterRequest {
    string profile_id = 1;
    string contact_id = 2;
    string name = 3;       // Required
}
```

#### 3.4 Migration

Existing roster entries get `name = 'default'`.

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
