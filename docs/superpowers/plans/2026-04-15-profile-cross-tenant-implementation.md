# Profile Service: Cross-Tenant Architecture Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the profile service a proper cross-tenant platform service with append-only property tracking and multi-list rosters.

**Architecture:** Three layers of changes — repository tenancy scoping (reads unscoped for identity data, tenant-scoped for organizational data), a property ledger model with JSONB cache, and a named roster model. All changes follow existing Frame/GORM patterns with TDD using tenant-scoped auth claims.

**Tech Stack:** Go 1.24+, Frame v1.94+, GORM, Connect RPC, protobuf, testcontainers (PostgreSQL + Keto)

**Spec:** `docs/superpowers/specs/2026-04-15-profile-cross-tenant-design.md`

---

### Task 1: Fix Cross-Tenant Repository Scoping for Profiles

**Files:**
- Modify: `apps/default/service/repository/profiles.go`
- Modify: `apps/default/service/business/profiles.go` (remove debug logging)
- Test: `apps/default/service/business/profiles_test.go`

- [ ] **Step 1: Add tenant-context test for CreateProfile (all types)**

Add test `Test_profileBusiness_CreateProfile_WithTenantContext` that creates PERSON, INSTITUTION, and BOT profiles with real auth claims in the context. This test should fail on current code due to the `GetByID` empty-claims issue.

```go
func (pts *ProfileTestSuite) Test_profileBusiness_CreateProfile_WithTenantContext() {
	t := pts.T()
	requestProp, _ := structpb.NewStruct(data.JSONMap{"au_name": "Tenant Context Tester"})

	testcases := []struct {
		name    string
		request *profilev1.CreateRequest
	}{
		{
			name: "Create person profile",
			request: &profilev1.CreateRequest{
				Type:       profilev1.ProfileType_PERSON,
				Contact:    "tenant.person@testing.com",
				Properties: requestProp,
			},
		},
		{
			name: "Create institution profile",
			request: &profilev1.CreateRequest{
				Type:       profilev1.ProfileType_INSTITUTION,
				Contact:    "tenant.org@testing.com",
				Properties: requestProp,
			},
		},
		{
			name: "Create bot profile",
			request: &profilev1.CreateRequest{
				Type:       profilev1.ProfileType_BOT,
				Contact:    "tenant.bot@testing.com",
				Properties: requestProp,
			},
		},
	}

	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		tenantID := util.IDString()
		partitionID := util.IDString()
		callerProfileID := util.IDString()
		ctx = pts.WithAuthClaims(ctx, tenantID, partitionID, callerProfileID)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				got, err := pb.CreateProfile(ctx, tt.request)
				require.NoError(t, err, "CreateProfile should succeed with tenant-scoped claims")
				require.Len(t, got.GetContacts(), 1)
				require.Equal(t, tt.request.GetProperties().AsMap(), got.GetProperties().AsMap())
			})
		}
	})
}
```

- [ ] **Step 2: Add cross-tenant GetByID test**

```go
func (pts *ProfileTestSuite) Test_profileBusiness_GetByID_WithTenantContext() {
	t := pts.T()
	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)

		tenantA := util.IDString()
		partitionA := util.IDString()
		ctxA := pts.WithAuthClaims(ctx, tenantA, partitionA, util.IDString())
		pb, _ := pts.getProfileBusiness(ctxA, svc)
		created, err := pb.CreateProfile(ctxA, &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "crosslookup@testing.com",
		})
		require.NoError(t, err)

		tenantB := util.IDString()
		partitionB := util.IDString()
		ctxB := pts.WithAuthClaims(ctx, tenantB, partitionB, util.IDString())
		pbB, _ := pts.getProfileBusiness(ctxB, svc)
		got, err := pbB.GetByID(ctxB, created.GetId())
		require.NoError(t, err, "GetByID should work cross-tenant")
		require.Equal(t, created.GetId(), got.GetId())
	})
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestProfileSuite/Test_profileBusiness_(CreateProfile_WithTenant|GetByID_WithTenant)" -v -count=1`

Expected: FAIL — `GetTypeByUID` and/or `GetByID` return "record not found" due to tenancy scoping on NULL/mismatched partition.

- [ ] **Step 4: Fix repository scoping**

In `apps/default/service/repository/profiles.go`, replace empty-claims pattern with `SkipTenancyChecksOnClaims`:

```go
func (pr *profileRepository) GetTypeByID(ctx context.Context, profileTypeID string) (*models.ProfileType, error) {
	profileType := &models.ProfileType{}
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	err := pr.Pool().DB(unscopedCtx, true).First(profileType, "id = ?", profileTypeID).Error
	return profileType, err
}

func (pr *profileRepository) GetTypeByUID(
	ctx context.Context,
	profileType profilev1.ProfileType,
) (*models.ProfileType, error) {
	profileTypeUID := models.ProfileTypeIDMap[profileType]
	profileTypeM := &models.ProfileType{}
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	err := pr.Pool().DB(unscopedCtx, true).First(profileTypeM, "uid = ?", profileTypeUID).Error
	return profileTypeM, err
}

func (pr *profileRepository) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	profile := &models.Profile{}
	err := pr.Pool().DB(unscopedCtx, true).Preload(clause.Associations).First(profile, "id = ?", id).Error
	return profile, err
}
```

Remove the `security` import line for `AuthenticationClaims` if it was only used for empty claims (it's still needed for `SkipTenancyChecksOnClaims`).

Also in `apps/default/service/business/profiles.go`, remove the debug logging from `CreateProfile` — revert to the clean version without `util.Log(ctx).With(...)` calls.

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestProfileSuite/Test_profileBusiness_(CreateProfile|GetByID_WithTenant)" -v -count=1`

Expected: ALL PASS

- [ ] **Step 6: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/repository/profiles.go apps/default/service/business/profiles.go apps/default/service/business/profiles_test.go
git commit -m "fix: use SkipTenancyChecksOnClaims for cross-tenant profile reads

Profile types (global seed data) and profiles (cross-tenant identity data)
must be readable regardless of caller's tenant. Replace empty-claims
workaround with SkipTenancyChecksOnClaims which correctly bypasses the
TenancyPartition scope.

Tests now use WithAuthClaims to match production conditions."
```

---

### Task 2: Fix Cross-Tenant Repository Scoping for Contacts

**Files:**
- Modify: `apps/default/service/repository/contacts.go`
- Test: `apps/default/service/business/profiles_test.go`

- [ ] **Step 1: Add cross-tenant GetByContact test**

```go
func (pts *ProfileTestSuite) Test_profileBusiness_GetByContact_WithTenantContext() {
	t := pts.T()
	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)

		tenantA := util.IDString()
		partitionA := util.IDString()
		ctxA := pts.WithAuthClaims(ctx, tenantA, partitionA, util.IDString())
		pb, _ := pts.getProfileBusiness(ctxA, svc)
		created, err := pb.CreateProfile(ctxA, &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_INSTITUTION,
			Contact: "crosscontact@testing.com",
		})
		require.NoError(t, err)

		tenantB := util.IDString()
		partitionB := util.IDString()
		ctxB := pts.WithAuthClaims(ctx, tenantB, partitionB, util.IDString())
		pbB, _ := pts.getProfileBusiness(ctxB, svc)
		got, err := pbB.GetByContact(ctxB, "crosscontact@testing.com")
		require.NoError(t, err, "GetByContact should work cross-tenant")
		require.Equal(t, created.GetId(), got.GetId())
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestProfileSuite/Test_profileBusiness_GetByContact_WithTenantContext" -v -count=1`

Expected: FAIL — contact lookup token query scoped by tenant.

- [ ] **Step 3: Fix contact repository read methods**

In `apps/default/service/repository/contacts.go`, add `security` import and apply `SkipTenancyChecksOnClaims` to cross-tenant read methods:

```go
import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)
```

```go
func (cr *contactRepository) GetByProfileID(ctx context.Context, profileID string) ([]*models.Contact, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	contactList := make([]*models.Contact, 0)
	err := cr.Pool().DB(unscopedCtx, true).Where("profile_id = ?", profileID).Find(&contactList).Error
	return contactList, err
}

func (cr *contactRepository) GetByLookupToken(
	ctx context.Context,
	lookupTokenList ...[]byte,
) ([]*models.Contact, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	var contactList []*models.Contact
	if err := cr.Pool().
		DB(unscopedCtx, true).
		Where(" look_up_token IN ?", lookupTokenList).
		Find(&contactList).
		Error; err != nil {
		return nil, err
	}
	return contactList, nil
}
```

Write methods (`DelinkFromProfile`, `VerificationSave`, `VerificationAttemptSave`) keep raw `ctx`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestProfileSuite/Test_profileBusiness_(CreateProfile|GetByID_WithTenant|GetByContact_WithTenant)" -v -count=1`

Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/repository/contacts.go apps/default/service/business/profiles_test.go
git commit -m "fix: use SkipTenancyChecksOnClaims for cross-tenant contact reads

Contact lookups by profile ID and lookup token are cross-tenant operations.
Apply SkipTenancyChecksOnClaims to read methods while keeping write methods
tenant-stamped."
```

---

### Task 3: Fix Cross-Tenant Repository Scoping for Addresses and Seed Data

**Files:**
- Modify: `apps/default/service/repository/addresses.go`
- Modify: `apps/default/service/repository/relationship.go`
- Test: `apps/default/service/business/address_test.go`

- [ ] **Step 1: Add cross-tenant address test**

In `apps/default/service/business/address_test.go`, add a test that creates a profile with an address under tenant A, then retrieves it from tenant B:

```go
func (ats *AddressTestSuite) Test_addressBusiness_GetByProfile_WithTenantContext() {
	t := ats.T()
	ats.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := ats.CreateService(t, dep)

		tenantA := util.IDString()
		partitionA := util.IDString()
		ctxA := ats.WithAuthClaims(ctx, tenantA, partitionA, util.IDString())
		pb, _ := ats.getProfileBusiness(ctxA, svc)

		profile, err := pb.CreateProfile(ctxA, &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "addr.cross@testing.com",
		})
		require.NoError(t, err)

		_, err = pb.AddAddress(ctxA, &profilev1.AddAddressRequest{
			Id: profile.GetId(),
			Address: &profilev1.AddressObject{
				Name:    "Test Office",
				Street:  "123 Main St",
				City:    "Nairobi",
				Country: "Kenya",
			},
		})
		require.NoError(t, err)

		// Read from tenant B
		tenantB := util.IDString()
		partitionB := util.IDString()
		ctxB := ats.WithAuthClaims(ctx, tenantB, partitionB, util.IDString())
		pbB, _ := ats.getProfileBusiness(ctxB, svc)

		got, err := pbB.GetByID(ctxB, profile.GetId())
		require.NoError(t, err, "profile with address should be readable cross-tenant")
		require.NotEmpty(t, got.GetAddresses(), "addresses should be visible cross-tenant")
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestAddressSuite/Test_addressBusiness_GetByProfile_WithTenantContext" -v -count=1`

Expected: FAIL — address repository `GetByProfileID` uses raw ctx.

- [ ] **Step 3: Fix address repository read methods**

In `apps/default/service/repository/addresses.go`, add `security` import and apply `SkipTenancyChecksOnClaims` to read methods. Country lookups are global seed data.

```go
import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"
	"gorm.io/gorm/clause"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)
```

Apply `SkipTenancyChecksOnClaims` to all read methods: `GetByNameAdminUnitAndCountry`, `GetByProfileID`, `CountryGetByISO3`, `CountryGetByAny`, `CountryGetByName`. Keep `SaveLink` and `DeleteLink` with raw ctx.

For each read method, the pattern is the same — add `unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)` as the first line and replace `ctx` with `unscopedCtx` in the `Pool().DB()` call.

- [ ] **Step 4: Fix relationship repository seed data lookups**

In `apps/default/service/repository/relationship.go`, add `security` import and apply `SkipTenancyChecksOnClaims` to `RelationshipType` and `RelationshipTypeByID` (global seed data). Leave `List` with raw ctx (tenant-scoped per spec).

```go
func (ar *relationshipRepository) RelationshipType(
	ctx context.Context,
	profileType profilev1.RelationshipType,
) (*models.RelationshipType, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	relationshipType := &models.RelationshipType{}
	relationshipTypeUID := models.RelationshipTypeIDMap[profileType]
	err := ar.Pool().DB(unscopedCtx, true).First(relationshipType, "uid = ?", relationshipTypeUID).Error
	return relationshipType, err
}

func (ar *relationshipRepository) RelationshipTypeByID(
	ctx context.Context,
	profileTypeID string,
) (*models.RelationshipType, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	relationshipType := &models.RelationshipType{}
	err := ar.Pool().DB(unscopedCtx, true).First(relationshipType, "id = ?", profileTypeID).Error
	return relationshipType, err
}
```

- [ ] **Step 5: Run all cross-tenant tests**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "WithTenantContext" -v -count=1`

Expected: ALL PASS

- [ ] **Step 6: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/repository/addresses.go apps/default/service/repository/relationship.go apps/default/service/business/address_test.go
git commit -m "fix: cross-tenant reads for addresses, countries, relationship types

Apply SkipTenancyChecksOnClaims to address reads (cross-tenant identity data)
and relationship type/country lookups (global seed data). Relationship List
stays tenant-scoped per design."
```

---

### Task 4: PropertyEntry Model and Migration

**Files:**
- Modify: `apps/default/service/models/models.go`
- Create: `apps/default/migrations/0001/20260415_property_entries.sql`
- Modify: `apps/default/service/repository/migrate.go`

- [ ] **Step 1: Add PropertyEntry model**

In `apps/default/service/models/models.go`, add the new model after the existing `Profile` struct:

```go
// PropertyEntry is an append-only ledger of property changes on a profile.
// The latest entry per (profile_id, key) determines the current value.
// Scoped entries are tenant-private and excluded from the JSONB cache.
type PropertyEntry struct {
	data.BaseModel
	ProfileID string `gorm:"type:varchar(50);not null;index:idx_prop_profile,priority:1"`
	Key       string `gorm:"type:varchar(255);not null;index:idx_prop_profile_key,priority:1"`
	Value     string `gorm:"type:text;not null"`
	Scoped    bool   `gorm:"not null;default:false;index:idx_prop_tenant_scoped"`
}
```

- [ ] **Step 2: Add PropertyEntry to GORM auto-migration**

Check `apps/default/service/repository/migrate.go` — the `Migrate` function registers models for auto-migration. Add `&models.PropertyEntry{}` to the model list.

```go
func Migrate(ctx context.Context, dsm datastore.Manager, migrationDir string) error {
	dbPool := dsm.GetPool(ctx, datastore.DefaultPoolName)
	return dbPool.SaveMigration(ctx,
		migration.WithAutoMigrate(
			// ... existing models ...
			&models.PropertyEntry{},
		),
		migration.WithMigrationDir(migrationDir),
	)
}
```

- [ ] **Step 3: Create SQL migration for composite indexes**

Create `apps/default/migrations/0001/20260415_property_entries.sql`:

```sql
CREATE INDEX IF NOT EXISTS idx_prop_profile_created
    ON property_entries (profile_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_prop_profile_key_created
    ON property_entries (profile_id, key, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_prop_tenant_scoped_lookup
    ON property_entries (profile_id, tenant_id, scoped)
    WHERE scoped = TRUE;

-- Explode existing Profile.Properties JSONB into property_entries
INSERT INTO property_entries (id, created_at, modified_at, created_by, modified_by, version, tenant_id, partition_id, access_id, profile_id, key, value, scoped)
SELECT
    gen_random_uuid()::varchar(50),
    p.created_at,
    p.created_at,
    COALESCE(p.created_by, ''),
    COALESCE(p.created_by, ''),
    0,
    COALESCE(p.tenant_id, ''),
    COALESCE(p.partition_id, ''),
    COALESCE(p.access_id, ''),
    p.id,
    kv.key,
    kv.value::text,
    FALSE
FROM profiles p,
     LATERAL jsonb_each_text(p.properties) AS kv(key, value)
WHERE p.properties IS NOT NULL
  AND p.properties != '{}'::jsonb
  AND p.deleted_at IS NULL;
```

- [ ] **Step 4: Verify migration runs**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestProfileSuite/Test_profileBusiness_CreateProfile$" -v -count=1`

Expected: PASS — migration creates the table, existing tests still work.

- [ ] **Step 5: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/models/models.go apps/default/service/repository/migrate.go apps/default/migrations/0001/20260415_property_entries.sql
git commit -m "feat: add PropertyEntry model and data migration

Append-only ledger table for profile property changes. Existing
Profile.Properties JSONB data is exploded into initial entries
with scoped=false and creator provenance from the profile."
```

---

### Task 5: PropertyEntry Repository

**Files:**
- Create: `apps/default/service/repository/property_entries.go`
- Modify: `apps/default/service/repository/interfaces.go`

- [ ] **Step 1: Define the repository interface**

In `apps/default/service/repository/interfaces.go`, add:

```go
type PropertyEntryRepository interface {
	datastore.BaseRepository[*models.PropertyEntry]
	// AppendEntries appends property entries for a profile. Uses raw ctx for provenance stamping.
	AppendEntries(ctx context.Context, entries []*models.PropertyEntry) error
	// LatestGlobalByProfile returns the latest global (scoped=false) entry per key.
	LatestGlobalByProfile(ctx context.Context, profileID string) ([]*models.PropertyEntry, error)
	// LatestScopedByProfileAndPartition returns the latest scoped entry per key for a partition.
	LatestScopedByProfileAndPartition(ctx context.Context, profileID, partitionID string) ([]*models.PropertyEntry, error)
	// HistoryByKey returns all entries for a profile+key, most recent first.
	// Scoped entries are filtered to the caller's tenant.
	HistoryByKey(ctx context.Context, profileID, key, callerTenantID string) ([]*models.PropertyEntry, error)
}
```

- [ ] **Step 2: Implement the repository**

Create `apps/default/service/repository/property_entries.go`:

```go
package repository

import (
	"context"

	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/datastore/pool"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/default/service/models"
)

type propertyEntryRepository struct {
	datastore.BaseRepository[*models.PropertyEntry]
}

func NewPropertyEntryRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) PropertyEntryRepository {
	return &propertyEntryRepository{
		BaseRepository: datastore.NewBaseRepository[*models.PropertyEntry](
			ctx, dbPool, workMan, func() *models.PropertyEntry { return &models.PropertyEntry{} },
		),
	}
}

func (r *propertyEntryRepository) AppendEntries(ctx context.Context, entries []*models.PropertyEntry) error {
	// Write with raw ctx to stamp tenant provenance via BaseModel.GenID
	return r.Pool().DB(ctx, false).Create(&entries).Error
}

func (r *propertyEntryRepository) LatestGlobalByProfile(ctx context.Context, profileID string) ([]*models.PropertyEntry, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	var entries []*models.PropertyEntry
	err := r.Pool().DB(unscopedCtx, true).
		Raw(`SELECT DISTINCT ON (key) * FROM property_entries
			 WHERE profile_id = ? AND scoped = FALSE AND deleted_at IS NULL
			 ORDER BY key, created_at DESC`, profileID).
		Scan(&entries).Error
	return entries, err
}

func (r *propertyEntryRepository) LatestScopedByProfileAndPartition(ctx context.Context, profileID, partitionID string) ([]*models.PropertyEntry, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	var entries []*models.PropertyEntry
	err := r.Pool().DB(unscopedCtx, true).
		Raw(`SELECT DISTINCT ON (key) * FROM property_entries
			 WHERE profile_id = ? AND scoped = TRUE AND partition_id = ? AND deleted_at IS NULL
			 ORDER BY key, created_at DESC`, profileID, partitionID).
		Scan(&entries).Error
	return entries, err
}

func (r *propertyEntryRepository) HistoryByKey(ctx context.Context, profileID, key, callerTenantID string) ([]*models.PropertyEntry, error) {
	unscopedCtx := security.SkipTenancyChecksOnClaims(ctx)
	var entries []*models.PropertyEntry
	err := r.Pool().DB(unscopedCtx, true).
		Where("profile_id = ? AND key = ? AND (scoped = FALSE OR tenant_id = ?) AND deleted_at IS NULL",
			profileID, key, callerTenantID).
		Order("created_at DESC").
		Find(&entries).Error
	return entries, err
}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd /home/j/code/antinvestor/service-profile && go build ./apps/default/...`

Expected: Success

- [ ] **Step 4: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/repository/property_entries.go apps/default/service/repository/interfaces.go
git commit -m "feat: add PropertyEntryRepository with cross-tenant reads

Append-only property entry repository. Writes use raw ctx for provenance
stamping. Reads use SkipTenancyChecksOnClaims for cross-tenant access.
LatestGlobalByProfile returns current global properties.
LatestScopedByProfileAndPartition returns partition-scoped properties."
```

---

### Task 6: Property Business Layer

**Files:**
- Modify: `apps/default/service/business/profiles.go`
- Test: `apps/default/service/business/profiles_test.go`

- [ ] **Step 1: Write tests for property operations**

```go
func (pts *ProfileTestSuite) Test_profileBusiness_UpdateProperties_Global() {
	t := pts.T()
	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		tenantID := util.IDString()
		partitionID := util.IDString()
		ctx = pts.WithAuthClaims(ctx, tenantID, partitionID, util.IDString())
		pb, _ := pts.getProfileBusiness(ctx, svc)

		profile, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "props.global@testing.com",
		})
		require.NoError(t, err)

		// Update global properties
		updateReq := &profilev1.UpdateRequest{
			Id:         profile.GetId(),
			Properties: data.JSONMap{"name": "Updated Name", "org": "Acme"}.ToProtoStruct(),
		}
		updated, err := pb.UpdateProfile(ctx, updateReq)
		require.NoError(t, err)
		require.Equal(t, "Updated Name", updated.GetProperties().AsMap()["name"])
		require.Equal(t, "Acme", updated.GetProperties().AsMap()["org"])

		// Read from different tenant — should see global properties
		ctxB := pts.WithAuthClaims(ctx, util.IDString(), util.IDString(), util.IDString())
		pbB, _ := pts.getProfileBusiness(ctxB, svc)
		got, err := pbB.GetByID(ctxB, profile.GetId())
		require.NoError(t, err)
		require.Equal(t, "Updated Name", got.GetProperties().AsMap()["name"])
	})
}

func (pts *ProfileTestSuite) Test_profileBusiness_UpdateProperties_Scoped() {
	t := pts.T()
	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)

		tenantA := util.IDString()
		partitionA := util.IDString()
		ctxA := pts.WithAuthClaims(ctx, tenantA, partitionA, util.IDString())
		pb, _ := pts.getProfileBusiness(ctxA, svc)

		profile, err := pb.CreateProfile(ctxA, &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "props.scoped@testing.com",
		})
		require.NoError(t, err)

		// Write scoped property
		scopedReq := &profilev1.UpdateRequest{
			Id:         profile.GetId(),
			Properties: data.JSONMap{"credit_score": "750"}.ToProtoStruct(),
			Scoped:     true,
		}
		_, err = pb.UpdateProfile(ctxA, scopedReq)
		require.NoError(t, err)

		// GetByID should NOT include scoped property
		got, err := pb.GetByID(ctxA, profile.GetId())
		require.NoError(t, err)
		_, hasCreditScore := got.GetProperties().AsMap()["credit_score"]
		require.False(t, hasCreditScore, "scoped property should not appear in GetByID")

		// GetByIDAndPartition with same partition SHOULD include it
		gotPartition, err := pb.GetByIDAndPartition(ctxA, profile.GetId(), partitionA)
		require.NoError(t, err)
		require.Equal(t, "750", gotPartition.GetProperties().AsMap()["credit_score"])

		// GetByIDAndPartition with different partition should NOT include it
		tenantB := util.IDString()
		partitionB := util.IDString()
		ctxB := pts.WithAuthClaims(ctx, tenantB, partitionB, util.IDString())
		pbB, _ := pts.getProfileBusiness(ctxB, svc)
		gotOther, err := pbB.GetByIDAndPartition(ctxB, profile.GetId(), partitionB)
		require.NoError(t, err)
		_, hasCreditScore = gotOther.GetProperties().AsMap()["credit_score"]
		require.False(t, hasCreditScore, "other partition should not see scoped property")
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestProfileSuite/Test_profileBusiness_UpdateProperties" -v -count=1`

Expected: Compilation error — `UpdateRequest.Scoped` and `GetByIDAndPartition` don't exist yet.

- [ ] **Step 3: Add GetByIDAndPartition to ProfileBusiness interface**

In `apps/default/service/business/profiles.go`, add to the `ProfileBusiness` interface:

```go
GetByIDAndPartition(ctx context.Context, profileID, partitionID string) (*profilev1.ProfileObject, error)
```

- [ ] **Step 4: Wire PropertyEntryRepository into profileBusiness**

Add the `propertyEntryRepo` field to the `profileBusiness` struct and update `NewProfileBusiness`:

```go
type profileBusiness struct {
	cfg              *config.ProfileConfig
	dek              *config.DEK
	contactBusiness  ContactBusiness
	addressBusiness  AddressBusiness
	profileRepo      repository.ProfileRepository
	propertyEntryRepo repository.PropertyEntryRepository
	eventsMan        frevents.Manager
}

func NewProfileBusiness(_ context.Context, cfg *config.ProfileConfig, dek *config.DEK,
	eventsMan frevents.Manager,
	contactBusiness ContactBusiness, addressBusiness AddressBusiness,
	profileRepo repository.ProfileRepository,
	propertyEntryRepo repository.PropertyEntryRepository) ProfileBusiness {
	return &profileBusiness{
		cfg:              cfg,
		dek:              dek,
		contactBusiness:  contactBusiness,
		addressBusiness:  addressBusiness,
		profileRepo:      profileRepo,
		propertyEntryRepo: propertyEntryRepo,
		eventsMan:        eventsMan,
	}
}
```

Update all callers (`handlers/profiles.go`, `tests/base_testsuite.go`, `business/profiles_test.go`) to pass the new repository.

- [ ] **Step 5: Implement UpdateProfile with property entries**

Modify `UpdateProfile` in `apps/default/service/business/profiles.go`. When `request.GetScoped()` is true, append scoped entries without rebuilding cache. When false, append global entries and rebuild cache.

```go
func (pb *profileBusiness) UpdateProfile(
	ctx context.Context,
	request *profilev1.UpdateRequest) (*profilev1.ProfileObject, error) {

	profile, err := pb.profileRepo.GetByID(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	requestProperties := data.JSONMap{}
	requestProperties = requestProperties.FromProtoStruct(request.GetProperties())

	// Append property entries
	var entries []*models.PropertyEntry
	for key, value := range requestProperties {
		entry := &models.PropertyEntry{
			ProfileID: profile.GetID(),
			Key:       key,
			Value:     fmt.Sprintf("%v", value),
			Scoped:    request.GetScoped(),
		}
		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		if err := pb.propertyEntryRepo.AppendEntries(ctx, entries); err != nil {
			return nil, data.ErrorConvertToAPI(err)
		}
	}

	// For global properties, rebuild the JSONB cache
	if !request.GetScoped() {
		latestEntries, err := pb.propertyEntryRepo.LatestGlobalByProfile(ctx, profile.GetID())
		if err != nil {
			return nil, data.ErrorConvertToAPI(err)
		}

		newProps := data.JSONMap{}
		for _, e := range latestEntries {
			newProps[e.Key] = e.Value
		}
		profile.Properties = newProps

		_, err = pb.profileRepo.Update(ctx, profile, "properties")
		if err != nil {
			return nil, data.ErrorConvertToAPI(err)
		}
	}

	return pb.ToAPI(ctx, profile)
}
```

- [ ] **Step 6: Implement GetByIDAndPartition**

```go
func (pb *profileBusiness) GetByIDAndPartition(
	ctx context.Context,
	profileID, partitionID string) (*profilev1.ProfileObject, error) {

	profileObj, err := pb.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	// Merge partition-scoped properties into the response
	scopedEntries, err := pb.propertyEntryRepo.LatestScopedByProfileAndPartition(ctx, profileID, partitionID)
	if err != nil {
		return nil, data.ErrorConvertToAPI(err)
	}

	if len(scopedEntries) > 0 {
		merged := profileObj.GetProperties().AsMap()
		for _, e := range scopedEntries {
			merged[e.Key] = e.Value
		}
		profileObj.Properties = data.JSONMap(merged).ToProtoStruct()
	}

	return profileObj, nil
}
```

- [ ] **Step 7: Run tests**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestProfileSuite/Test_profileBusiness_(UpdateProperties|CreateProfile)" -v -count=1`

Expected: ALL PASS

- [ ] **Step 8: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/business/profiles.go apps/default/service/business/profiles_test.go apps/default/service/handlers/profiles.go apps/default/tests/base_testsuite.go
git commit -m "feat: append-only property ledger with global and scoped entries

UpdateProfile now appends property entries instead of direct JSONB mutation.
Global properties rebuild the JSONB cache. Scoped properties are tenant-private
and only visible via GetByIDAndPartition. Includes tests with tenant auth claims."
```

---

### Task 7: Proto Changes and Handler Wiring

**Files:**
- Modify: `proto/profile/profile/v1/profile.proto`
- Modify: `apps/default/service/handlers/profiles.go`

- [ ] **Step 1: Add scoped field to UpdateRequest**

In `proto/profile/profile/v1/profile.proto`, add to the existing `UpdateRequest`:

```protobuf
message UpdateRequest {
    string id = 1;
    google.protobuf.Struct properties = 2;
    common.v1.STATE state = 3;
    bool scoped = 4;
}
```

- [ ] **Step 2: Add GetByIDAndPartition RPC messages**

```protobuf
message GetByIDAndPartitionRequest {
    string id = 1;
    string partition_id = 2;
}

message GetByIDAndPartitionResponse {
    ProfileObject data = 1;
}
```

- [ ] **Step 3: Add PropertyHistory RPC messages**

```protobuf
message PropertyHistoryRequest {
    string id = 1;
    string key = 2;
}

message PropertyEntryObject {
    string key = 1;
    string value = 2;
    string tenant_id = 3;
    string created_by = 4;
    google.protobuf.Timestamp created_at = 5;
    bool scoped = 6;
}

message PropertyHistoryResponse {
    repeated PropertyEntryObject entries = 1;
}
```

- [ ] **Step 4: Add RPCs to the service definition**

```protobuf
service ProfileService {
    // ... existing RPCs ...
    rpc GetByIDAndPartition(GetByIDAndPartitionRequest) returns (GetByIDAndPartitionResponse);
    rpc PropertyHistory(PropertyHistoryRequest) returns (PropertyHistoryResponse);
}
```

- [ ] **Step 5: Regenerate proto**

Run: `cd /home/j/code/antinvestor/service-profile && make buf-generate` (or the project's proto generation command)

- [ ] **Step 6: Implement GetByIDAndPartition handler**

In `apps/default/service/handlers/profiles.go`:

```go
func (ps *ProfileServer) GetByIDAndPartition(
	ctx context.Context,
	request *connect.Request[profilev1.GetByIDAndPartitionRequest],
) (*connect.Response[profilev1.GetByIDAndPartitionResponse], error) {
	profileObj, err := ps.profileBusiness.GetByIDAndPartition(ctx, request.Msg.GetId(), request.Msg.GetPartitionId())
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	auditlib.WithResource(ctx, auditlib.ResourceProfile, request.Msg.GetId())

	return connect.NewResponse(&profilev1.GetByIDAndPartitionResponse{Data: profileObj}), nil
}
```

- [ ] **Step 7: Implement PropertyHistory handler**

```go
func (ps *ProfileServer) PropertyHistory(
	ctx context.Context,
	request *connect.Request[profilev1.PropertyHistoryRequest],
) (*connect.Response[profilev1.PropertyHistoryResponse], error) {
	claims := security.ClaimsFromContext(ctx)
	callerTenantID := ""
	if claims != nil {
		callerTenantID = claims.GetTenantID()
	}

	entries, err := ps.profileBusiness.GetPropertyHistory(ctx, request.Msg.GetId(), request.Msg.GetKey(), callerTenantID)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}

	var protoEntries []*profilev1.PropertyEntryObject
	for _, e := range entries {
		protoEntries = append(protoEntries, &profilev1.PropertyEntryObject{
			Key:       e.Key,
			Value:     e.Value,
			TenantId:  e.TenantID,
			CreatedBy: e.CreatedBy,
			CreatedAt: timestamppb.New(e.CreatedAt),
			Scoped:    e.Scoped,
		})
	}

	return connect.NewResponse(&profilev1.PropertyHistoryResponse{Entries: protoEntries}), nil
}
```

- [ ] **Step 8: Add GetPropertyHistory to business interface and implement**

In `apps/default/service/business/profiles.go`:

```go
// In ProfileBusiness interface:
GetPropertyHistory(ctx context.Context, profileID, key, callerTenantID string) ([]*models.PropertyEntry, error)

// Implementation:
func (pb *profileBusiness) GetPropertyHistory(ctx context.Context, profileID, key, callerTenantID string) ([]*models.PropertyEntry, error) {
	return pb.propertyEntryRepo.HistoryByKey(ctx, profileID, key, callerTenantID)
}
```

- [ ] **Step 9: Build and verify**

Run: `cd /home/j/code/antinvestor/service-profile && go build ./apps/default/...`

Expected: Success

- [ ] **Step 10: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add proto/ apps/default/service/handlers/profiles.go apps/default/service/business/profiles.go sdk/
git commit -m "feat: add GetByIDAndPartition and PropertyHistory RPCs

New RPCs for partition-scoped profile reads and property change history.
UpdateRequest gains a scoped flag for tenant-private properties."
```

---

### Task 8: Multi-List Roster Model and Migration

**Files:**
- Modify: `apps/default/service/models/models.go`
- Create: `apps/default/migrations/0001/20260415_roster_name.sql`

- [ ] **Step 1: Add Name field to Roster model**

In `apps/default/service/models/models.go`, update the `Roster` struct:

```go
type Roster struct {
	data.BaseModel
	ProfileID  string       `gorm:"type:varchar(50);uniqueIndex:roster_composite_index,priority:1"`
	ContactID  string       `gorm:"type:varchar(50);uniqueIndex:roster_composite_index,priority:2"`
	Name       string       `gorm:"type:varchar(255);not null;default:'default';uniqueIndex:roster_composite_index,priority:3"`
	Contact    *Contact     `gorm:"foreignKey:ContactID"`
	Properties data.JSONMap `gorm:"type:JSONB"`
}
```

- [ ] **Step 2: Create SQL migration for index change**

Create `apps/default/migrations/0001/20260415_roster_name.sql`:

```sql
-- Drop old composite index (profile_id, contact_id only)
DROP INDEX IF EXISTS roster_composite_index;

-- GORM auto-migration will create the new index with (profile_id, contact_id, name)
-- Set default name for existing entries
UPDATE rosters SET name = 'default' WHERE name IS NULL OR name = '';
```

- [ ] **Step 3: Update Roster.ToAPI to include Name**

In `apps/default/service/models/models.go`, update the `ToAPI` method:

```go
func (r *Roster) ToAPI(dek *config.DEK) (*profilev1.RosterObject, error) {
	rosterObj := &profilev1.RosterObject{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
	}
	// ... rest of existing ToAPI logic ...
}
```

- [ ] **Step 4: Verify migration runs and build succeeds**

Run: `cd /home/j/code/antinvestor/service-profile && go build ./apps/default/...`

Expected: Success (proto changes for RosterObject.name field needed — may need proto update first)

- [ ] **Step 5: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/models/models.go apps/default/migrations/0001/20260415_roster_name.sql
git commit -m "feat: add Name field to Roster model for multi-list support

Roster entries now have a freeform Name field. Existing entries default
to 'default'. Unique index changes to (profile_id, contact_id, name)."
```

---

### Task 9: Multi-List Roster Business and Repository

**Files:**
- Modify: `apps/default/service/repository/roster.go`
- Modify: `apps/default/service/repository/interfaces.go`
- Modify: `apps/default/service/business/roster.go`
- Test: `apps/default/service/business/roster_test.go`

- [ ] **Step 1: Write test for multi-list roster with tenant isolation**

```go
func (rts *RosterTestSuite) TestRosterBusiness_MultiList_TenantIsolation() {
	t := rts.T()
	rts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := rts.CreateService(t, dep)

		tenantA := util.IDString()
		partitionA := util.IDString()
		ctxA := rts.WithAuthClaims(ctx, tenantA, partitionA, util.IDString())

		pb, _ := rts.getProfileBusiness(ctxA, svc)
		profile, err := pb.CreateProfile(ctxA, &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: "roster.multi@testing.com",
		})
		require.NoError(t, err)

		rb := rts.getRosterBusiness(ctxA, svc)

		// Add contacts to "friends" list
		friendsReq := &profilev1.AddRosterRequest{
			ProfileId: profile.GetId(),
			Name:      "friends",
			Data: []*profilev1.RawContact{
				{Contact: "friend1@testing.com"},
				{Contact: "friend2@testing.com"},
			},
		}
		friends, err := rb.CreateRoster(ctxA, friendsReq)
		require.NoError(t, err)
		require.Len(t, friends, 2)

		// Add one of the same contacts to "colleagues" list
		colleaguesReq := &profilev1.AddRosterRequest{
			ProfileId: profile.GetId(),
			Name:      "colleagues",
			Data: []*profilev1.RawContact{
				{Contact: "friend1@testing.com"},
				{Contact: "colleague1@testing.com"},
			},
		}
		colleagues, err := rb.CreateRoster(ctxA, colleaguesReq)
		require.NoError(t, err)
		require.Len(t, colleagues, 2)

		// Search "friends" — should return 2
		friendSearch := &profilev1.SearchRosterRequest{
			ProfileId: profile.GetId(),
			Name:      "friends",
		}
		// ... verify search returns 2 results

		// Tenant B should see nothing
		tenantB := util.IDString()
		partitionB := util.IDString()
		ctxB := rts.WithAuthClaims(ctx, tenantB, partitionB, util.IDString())
		rbB := rts.getRosterBusiness(ctxB, svc)
		allSearch := &profilev1.SearchRosterRequest{
			ProfileId: profile.GetId(),
		}
		// ... verify search returns 0 results for tenant B
	})
}
```

- [ ] **Step 2: Update roster repository interface**

In `apps/default/service/repository/interfaces.go`, update `RosterRepository` to add name-aware lookup:

```go
type RosterRepository interface {
	datastore.BaseRepository[*models.Roster]
	GetByContactAndProfileID(ctx context.Context, profileID, contactID string) (*models.Roster, error)
	GetByContactAndProfileIDAndName(ctx context.Context, profileID, contactID, name string) (*models.Roster, error)
	GetByContactIDsAndProfileID(ctx context.Context, contactIDs []string, profileID string) ([]*models.Roster, error)
	GetByContactIDsAndProfileIDAndName(ctx context.Context, contactIDs []string, profileID, name string) ([]*models.Roster, error)
	Search(ctx context.Context, query *data.SearchQuery) (workerpool.JobResultPipe[[]*models.Roster], error)
}
```

- [ ] **Step 3: Implement name-aware repository methods**

In `apps/default/service/repository/roster.go`:

```go
func (rr *rosterRepository) GetByContactAndProfileIDAndName(
	ctx context.Context,
	profileID, contactID, name string,
) (*models.Roster, error) {
	roster := &models.Roster{}
	err := rr.Pool().DB(ctx, true).
		Preload(clause.Associations).
		Where("profile_id = ? AND contact_id = ? AND name = ?", profileID, contactID, name).
		First(roster).
		Error
	return roster, err
}

func (rr *rosterRepository) GetByContactIDsAndProfileIDAndName(
	ctx context.Context,
	contactIDs []string,
	profileID, name string,
) ([]*models.Roster, error) {
	rosterList := make([]*models.Roster, 0, len(contactIDs))
	err := rr.Pool().DB(ctx, true).
		Where("profile_id = ? AND contact_id IN ? AND name = ?", profileID, contactIDs, name).
		Find(&rosterList).
		Error
	return rosterList, err
}
```

- [ ] **Step 4: Update roster business layer**

In `apps/default/service/business/roster.go`, update `CreateRoster` and `processRosterBatch` to accept and use the `Name` field from the request. Update `Search` to filter by name when provided.

Key changes:
- `CreateRoster`: pass `request.GetName()` (default to `"default"` if empty) to `processRosterBatch`
- `processRosterBatch`: use `GetByContactIDsAndProfileIDAndName` instead of `GetByContactIDsAndProfileID`
- `findRostersToCreate`: set `Name` on each new `Roster`
- `Search`: add `name = ?` filter to `SearchQuery` when name is not empty

- [ ] **Step 5: Update proto for roster Name field**

In `proto/profile/profile/v1/profile.proto`, update relevant messages:

```protobuf
message RosterObject {
    string id = 1;
    string profile_id = 2;
    ContactObject contact = 3;
    google.protobuf.Struct extra = 4;
    string name = 5;
}

message AddRosterRequest {
    repeated RawContact data = 1;
    string name = 2;
}

message SearchRosterRequest {
    // ... existing fields ...
    string name = 9;
}

message RemoveRosterRequest {
    string id = 1;
    string name = 2;
}
```

Regenerate proto.

- [ ] **Step 6: Run tests**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/service/business/ -run "TestRosterSuite" -v -count=1`

Expected: ALL PASS

- [ ] **Step 7: Commit**

```bash
cd /home/j/code/antinvestor/service-profile
git add apps/default/service/repository/roster.go apps/default/service/repository/interfaces.go apps/default/service/business/roster.go apps/default/service/business/roster_test.go apps/default/service/models/models.go proto/ sdk/
git commit -m "feat: multi-list roster with name field and tenant isolation

Rosters now support freeform named lists. Same contact can appear in
multiple lists within the same tenant. Roster reads and writes remain
strictly tenant-scoped. Search supports filtering by list name."
```

---

### Task 10: Integration Verification and Cleanup

**Files:**
- Test: all test files
- Modify: `apps/default/service/business/profiles.go` (remove any remaining debug logging)

- [ ] **Step 1: Run full test suite**

Run: `cd /home/j/code/antinvestor/service-profile && go test ./apps/default/... -v -count=1`

Verify all tests pass, including pre-existing tests.

- [ ] **Step 2: Run linter**

Run: `cd /home/j/code/antinvestor/service-profile && make lint` (or `golangci-lint run ./...`)

Fix any lint issues.

- [ ] **Step 3: Verify build**

Run: `cd /home/j/code/antinvestor/service-profile && go build ./...`

Expected: Clean build

- [ ] **Step 4: Final commit and release**

```bash
cd /home/j/code/antinvestor/service-profile
git add -A
git commit -m "chore: integration cleanup and lint fixes"
```

Create a release that includes all changes:

```bash
gh release create v1.29.0 --title "v1.29.0 — Cross-tenant architecture, property ledger, multi-list roster" --target main --notes "$(cat <<'EOF'
## Cross-Tenant Repository Layer
- Profile, contact, and address reads use SkipTenancyChecksOnClaims (cross-tenant)
- Roster and relationship reads remain tenant-scoped
- All writes stamp creator provenance (tenant/partition/access)

## Append-Only Property Ledger
- PropertyEntry model tracks all property changes with full provenance
- Global properties (scoped=false) rebuild Profile.Properties JSONB cache
- Tenant-scoped properties (scoped=true) visible only via GetByIDAndPartition
- PropertyHistory RPC returns full change log per key

## Multi-List Roster
- Roster entries gain a Name field for freeform list categorization
- Same contact can appear in multiple lists within a tenant
- Roster isolation: tenant A's lists invisible to tenant B

## Migration
- Existing Profile.Properties exploded into property_entries
- Existing rosters get name='default'
EOF
)"
```
