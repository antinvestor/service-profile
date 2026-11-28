package business_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	settingsv1 "buf.build/gen/go/antinvestor/settingz/protocolbuffers/go/settings/v1"
	"connectrpc.com/connect/v2"
	"github.com/pitabwire/frame/v3/frametests/definition"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

func settingKey(name string) *settingsv1.Setting {
	return &settingsv1.Setting{
		Name:     name,
		Object:   "profile",
		ObjectId: "obj-1",
		Lang:     "en",
		Module:   "UpdateTests",
	}
}

func (ts *SettingsTestSuite) TestSet_UpdatesExistingKey() {
	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		ctx, svc := ts.CreateService(t, depOpt)
		nb, _, valRepo := ts.getSettingBusiness(ctx, svc)

		key := settingKey("set.update.existing")

		first, err := nb.Set(ctx, &settingsv1.SetRequest{Key: key, Value: "first"})
		require.NoError(t, err, "first Set must create the value")

		for _, value := range []string{"second", "third"} {
			got, setErr := nb.Set(ctx, &settingsv1.SetRequest{Key: key, Value: value})
			require.NoError(t, setErr, "Set on an existing key must update, not insert")
			require.Equal(t, first.GetData().GetId(), got.GetData().GetId(),
				"update must keep the same setting value row")
			require.Equal(t, value, got.GetData().GetValue())
		}

		resp, err := nb.Get(ctx, &settingsv1.GetRequest{Key: key})
		require.NoError(t, err)
		require.Equal(t, "third", resp.GetData().GetValue(), "Get must return the latest value")

		stored, err := valRepo.GetByID(ctx, first.GetData().GetId())
		require.NoError(t, err)
		require.Equal(t, uint(3), stored.GetVersion(),
			"version is the optimistic-lock counter: 1 on create, +1 per update")

		vals, err := valRepo.GetByRef(ctx, stored.Ref)
		require.NoError(t, err)
		require.Len(t, vals, 1, "updates must not insert extra value rows")
	})
}

func (ts *SettingsTestSuite) TestSet_UpdatesLegacyZeroVersionRow() {
	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		ctx, svc := ts.CreateService(t, depOpt)
		nb, _, valRepo := ts.getSettingBusiness(ctx, svc)

		key := settingKey("set.update.legacy")
		first, err := nb.Set(ctx, &settingsv1.SetRequest{Key: key, Value: "legacy"})
		require.NoError(t, err)

		// Rows written before the fix were persisted with version 0.
		res := valRepo.Pool().DB(ctx, false).
			Exec("UPDATE setting_vals SET version = 0 WHERE id = ?", first.GetData().GetId())
		require.NoError(t, res.Error)
		require.Equal(t, int64(1), res.RowsAffected)

		got, err := nb.Set(ctx, &settingsv1.SetRequest{Key: key, Value: "migrated"})
		require.NoError(t, err, "a version-0 row from before the fix must still be updatable")
		require.Equal(t, "migrated", got.GetData().GetValue())

		stored, err := valRepo.GetByID(ctx, first.GetData().GetId())
		require.NoError(t, err)
		require.Positive(t, stored.GetVersion(), "the update must move a legacy row onto a real lock version")
		require.Equal(t, "migrated", stored.Detail)
	})
}

func (ts *SettingsTestSuite) TestSettingValUpdate_RejectsStaleVersion() {
	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		ctx, svc := ts.CreateService(t, depOpt)
		nb, _, valRepo := ts.getSettingBusiness(ctx, svc)

		created, err := nb.Set(ctx, &settingsv1.SetRequest{Key: settingKey("set.update.stale"), Value: "v1"})
		require.NoError(t, err)

		readerA, err := valRepo.GetByID(ctx, created.GetData().GetId())
		require.NoError(t, err)
		readerB, err := valRepo.GetByID(ctx, created.GetData().GetId())
		require.NoError(t, err)

		readerA.Detail = "from A"
		require.NoError(t, valRepo.UpdateDetail(ctx, readerA))
		require.Equal(t, uint(2), readerA.GetVersion())

		readerB.Detail = "from B"
		err = valRepo.UpdateDetail(ctx, readerB)
		require.ErrorIs(t, err, repository.ErrSettingValVersionConflict,
			"an update based on a stale version must not overwrite the newer value")

		stored, err := valRepo.GetByID(ctx, created.GetData().GetId())
		require.NoError(t, err)
		require.Equal(t, "from A", stored.Detail)
	})
}

func (ts *SettingsTestSuite) TestSet_ConcurrentUpdatesNeverLoseAcknowledgedWrites() {
	ts.WithTestDependancies(ts.T(), func(t *testing.T, depOpt *definition.DependencyOption) {
		ctx, svc := ts.CreateService(t, depOpt)
		nb, _, valRepo := ts.getSettingBusiness(ctx, svc)

		key := settingKey("set.update.concurrent")
		created, err := nb.Set(ctx, &settingsv1.SetRequest{Key: key, Value: "initial"})
		require.NoError(t, err)

		const writers = 8
		var (
			wg        sync.WaitGroup
			mu        sync.Mutex
			succeeded []string
			aborted   int
		)
		start := make(chan struct{})
		for i := range writers {
			wg.Add(1)
			go func(value string) {
				defer wg.Done()
				<-start
				_, setErr := nb.Set(ctx, &settingsv1.SetRequest{Key: key, Value: value})
				mu.Lock()
				defer mu.Unlock()
				if setErr == nil {
					succeeded = append(succeeded, value)
					return
				}
				var connectErr *connect.Error
				if errors.As(setErr, &connectErr) && connectErr.Code() == connect.CodeAborted {
					aborted++
					return
				}
				t.Errorf("unexpected Set error: %v", setErr)
			}(fmt.Sprintf("writer-%d", i))
		}
		close(start)
		wg.Wait()

		require.NotEmpty(t, succeeded, "at least one concurrent writer must win")
		require.Equal(t, writers, len(succeeded)+aborted)

		stored, err := valRepo.GetByID(ctx, created.GetData().GetId())
		require.NoError(t, err)
		require.Equal(t, uint(1+len(succeeded)), stored.GetVersion(),
			"every acknowledged write must bump the version exactly once")
		require.Contains(t, succeeded, stored.Detail, "the stored value must come from an acknowledged write")

		vals, err := valRepo.GetByRef(ctx, stored.Ref)
		require.NoError(t, err)
		require.Len(t, vals, 1)
	})
}
