package business_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"connectrpc.com/connect/v2"
	"github.com/pitabwire/frame/v3/data"
	"github.com/pitabwire/frame/v3/datastore"
	"github.com/pitabwire/frame/v3/frametests/definition"
	"github.com/pitabwire/frame/v3/security"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/default/service/business"
)

func (pts *ProfileTestSuite) createNamed(
	ctx context.Context, t *testing.T, pb business.ProfileBusiness, contact string, props data.JSONMap,
) *profilev1.ProfileObject {
	t.Helper()
	p, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
		Type:       profilev1.ProfileType_PERSON,
		Contact:    contact,
		Properties: props.ToProtoStruct(),
	})
	require.NoError(t, err)
	return p
}

func (pts *ProfileTestSuite) Test_profileBusiness_Handles() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		first := pts.createNamed(ctx, t, pb, "handle.one@testing.com", data.JSONMap{"name": "Handle Person"})
		require.Equal(t, "handle.person", first.GetHandle(), "the first profile with a name gets it bare")

		second := pts.createNamed(ctx, t, pb, "handle.two@testing.com", data.JSONMap{"au_name": "Handle  Person!"})
		id := second.GetId()
		require.Equal(t, "handle.person."+id[len(id)-8:], second.GetHandle(),
			"the next one with the same name gets the tail of its id")

		nameless := pts.createNamed(ctx, t, pb, "handle.three@testing.com", data.JSONMap{})
		require.True(t, strings.HasPrefix(nameless.GetHandle(), "user."), nameless.GetHandle())

		// The handle is returned on every read and does not follow the name.
		renamed := data.JSONMap{"name": "Renamed Person"}
		_, err := pb.UpdateProfile(ctx, &profilev1.UpdateRequest{
			Id:         first.GetId(),
			Properties: renamed.ToProtoStruct(),
		})
		require.NoError(t, err)
		reread, err := pb.GetByID(ctx, first.GetId())
		require.NoError(t, err)
		require.Equal(t, "handle.person", reread.GetHandle(), "a handle never changes")

		byHandle, err := pb.GetByHandle(ctx, " @Handle.Person ")
		require.NoError(t, err)
		require.Equal(t, first.GetId(), byHandle.GetId())

		_, err = pb.GetByHandle(ctx, "no.such.handle")
		require.Equal(t, connect.CodeNotFound, connect.CodeOf(err), "%v", err)

		_, err = pb.GetByHandle(ctx, "@a")
		require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err), "%v", err)

		found := pts.searchIDs(ctx, t, pb, "@Handle.Pers")
		require.ElementsMatch(t, []string{first.GetId(), second.GetId()}, found)

		_, err = pb.SearchProfile(ctx, &profilev1.SearchRequest{Query: "@bad%", Count: 10})
		require.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err), "%v", err)
	})
}

func (pts *ProfileTestSuite) searchIDs(
	ctx context.Context, t *testing.T, pb business.ProfileBusiness, query string,
) []string {
	t.Helper()
	result, err := pb.SearchProfile(ctx, &profilev1.SearchRequest{Query: query, Count: 50})
	require.NoError(t, err)
	var ids []string
	for {
		res, ok := result.ReadResult(ctx)
		if !ok {
			return ids
		}
		require.NoError(t, res.Error())
		for _, p := range res.Item() {
			ids = append(ids, p.GetID())
		}
	}
}

// The migration gives profiles created before handles existed one, by the
// same rules as new profiles.
func (pts *ProfileTestSuite) Test_profileHandles_Backfill() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		taken := pts.createNamed(ctx, t, pb, "backfill.taken@testing.com", data.JSONMap{"name": "Backfill Taken"})
		older := pts.createNamed(ctx, t, pb, "backfill.a@testing.com", data.JSONMap{"name": "Backfill Twin"})
		newer := pts.createNamed(ctx, t, pb, "backfill.b@testing.com", data.JSONMap{"au_name": "backfill-twin"})
		clash := pts.createNamed(ctx, t, pb, "backfill.c@testing.com", data.JSONMap{"name": "Backfill Taken"})
		anon := pts.createNamed(ctx, t, pb, "backfill.d@testing.com", data.JSONMap{"name": "李"})

		db := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName).
			DB(security.SkipTenancyChecksOnClaims(ctx), false)
		ids := []string{older.GetId(), newer.GetId(), clash.GetId(), anon.GetId()}
		require.NoError(t, db.Exec("UPDATE profiles SET handle = '' WHERE id IN ?", ids).Error)

		sql, err := os.ReadFile("../../migrations/0001/20261103_profile_handles.sql")
		require.NoError(t, err)
		backfill, _, found := strings.Cut(string(sql), "CREATE UNIQUE INDEX")
		require.True(t, found)
		require.NoError(t, db.Exec(backfill).Error)

		tail := func(id string) string { return id[len(id)-8:] }
		handleOf := func(id string) string {
			p, getErr := pb.GetByID(ctx, id)
			require.NoError(t, getErr)
			return p.GetHandle()
		}
		require.Equal(t, "backfill.taken", handleOf(taken.GetId()), "an assigned handle is left alone")
		require.Equal(t, "backfill.twin", handleOf(older.GetId()), "the oldest gets the bare name")
		require.Equal(t, "backfill.twin."+tail(newer.GetId()), handleOf(newer.GetId()))
		require.Equal(t, "backfill.taken."+tail(clash.GetId()), handleOf(clash.GetId()),
			"a name another profile already holds gets the id tail")
		require.Equal(t, "user."+tail(anon.GetId()), handleOf(anon.GetId()))
	})
}

// Profiles with the same name created at the same moment all succeed, each
// with its own handle: a candidate claimed in between is skipped.
func (pts *ProfileTestSuite) Test_profileHandles_ConcurrentSameName() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		const n = 6
		type created struct {
			handle string
			err    error
		}
		results := make(chan created, n)
		props := data.JSONMap{"name": "Race Person"}
		for i := range n {
			go func() {
				p, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
					Type:       profilev1.ProfileType_PERSON,
					Contact:    fmt.Sprintf("race.%d@testing.com", i),
					Properties: props.ToProtoStruct(),
				})
				results <- created{handle: p.GetHandle(), err: err}
			}()
		}

		seen := map[string]bool{}
		bare := 0
		for range n {
			r := <-results
			require.NoError(t, r.err)
			require.False(t, seen[r.handle], "handle %q given twice", r.handle)
			seen[r.handle] = true
			if r.handle == "race.person" {
				bare++
			}
		}
		require.Equal(t, 1, bare, "exactly one profile gets the bare name")
	})
}
