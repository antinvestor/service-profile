package business_test

import (
	"context"
	"fmt"
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

// The database refuses a second profile with a handle, through the index
// createWithHandle recognises, deleted profiles included.
func (pts *ProfileTestSuite) Test_profileHandles_UniqueIndex() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		held := pts.createNamed(ctx, t, pb, "unique.a@testing.com", data.JSONMap{"name": "Unique Holder"})
		other := pts.createNamed(ctx, t, pb, "unique.b@testing.com", data.JSONMap{"name": "Someone Else"})

		db := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName).
			DB(security.SkipTenancyChecksOnClaims(ctx), false)
		err := db.Exec("UPDATE profiles SET handle = ? WHERE id = ?", held.GetHandle(), other.GetId()).Error
		require.ErrorContains(t, err, "idx_profiles_handle")

		require.NoError(t, db.Exec("UPDATE profiles SET deleted_at = now() WHERE id = ?", held.GetId()).Error)
		err = db.Exec("UPDATE profiles SET handle = ? WHERE id = ?", held.GetHandle(), other.GetId()).Error
		require.ErrorContains(t, err, "idx_profiles_handle", "a deleted profile keeps its handle")

		_, err = pb.GetByHandle(ctx, held.GetHandle())
		require.Equal(t, connect.CodeNotFound, connect.CodeOf(err), "a deleted profile is not found by handle")
	})
}

// Profiles seeded by migrations carry handles too: no profile is ever
// without one.
func (pts *ProfileTestSuite) Test_profileHandles_SeededProfilesHaveHandles() {
	t := pts.T()

	pts.WithTestDependancies(t, func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		db := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName).
			DB(security.SkipTenancyChecksOnClaims(ctx), true)
		var seeded, missing int64
		require.NoError(t, db.Raw("SELECT count(*) FROM profiles").Scan(&seeded).Error)
		require.Positive(t, seeded, "the migrations seed profiles")
		require.NoError(t, db.Raw("SELECT count(*) FROM profiles WHERE handle IS NULL OR handle = ''").
			Scan(&missing).Error)
		require.Zero(t, missing)

		bot, err := pb.GetByHandle(ctx, "@service.notification.whatsapp")
		require.NoError(t, err)
		require.Equal(t, "daenc7kpf2t8pa1q04hg", bot.GetId())
	})
}
