package authz_test

import (
	"context"
	"testing"

	"github.com/pitabwire/frame/frametests/definition"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/geolocation/service/authz"
	geotests "github.com/antinvestor/service-profile/apps/geolocation/tests"
)

const (
	authzTenantID    = "tenant-authz"
	authzPartitionID = "partition-authz"
)

type MiddlewareSuite struct {
	geotests.GeolocationBaseTestSuite
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}

func (s *MiddlewareSuite) authzContext(ctx context.Context, subjectID string) context.Context {
	return s.WithAuthClaims(ctx, authzTenantID, authzPartitionID, subjectID)
}

func (s *MiddlewareSuite) TestMiddlewarePermissionChecks() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		ctx := s.authzContext(baseCtx, "subject-1")
		mw := authz.NewMiddleware(svc.SecurityManager().GetAuthorizer(ctx))

		s.SeedTenantAccess(ctx, svc, authzTenantID, authzPartitionID, "subject-1")
		s.SeedTenantRole(ctx, svc, authzTenantID, authzPartitionID, "subject-1", authz.RoleOwner)

		require.NoError(t, mw.CanGeolocationViewSelf(ctx, "subject-1"))
		require.NoError(t, mw.CanLocationIngestSelf(ctx, "subject-1"))
		require.NoError(t, mw.CanGeolocationViewSelf(ctx, "subject-2"))
		require.NoError(t, mw.CanLocationIngestSelf(ctx, "subject-2"))
		require.NoError(t, mw.CanGeolocationManage(ctx))
		require.NoError(t, mw.CanGeolocationView(ctx))
		require.NoError(t, mw.CanLocationIngest(ctx))
	})
}

func (s *MiddlewareSuite) TestMiddlewarePermissionDenied() {
	s.WithTestDependencies(s.T(), func(t *testing.T, dep *definition.DependencyOption) {
		baseCtx, svc := s.CreateService(t, dep)
		ctx := s.authzContext(baseCtx, "subject-2")
		mw := authz.NewMiddleware(svc.SecurityManager().GetAuthorizer(ctx))

		s.SeedTenantAccess(ctx, svc, authzTenantID, authzPartitionID, "subject-2")
		s.SeedTenantRole(ctx, svc, authzTenantID, authzPartitionID, "subject-2", authz.RoleViewer)

		require.NoError(t, mw.CanGeolocationView(ctx))
		require.NoError(t, mw.CanGeolocationViewSelf(ctx, "subject-2"))
		require.Error(t, mw.CanGeolocationManage(ctx))
		require.Error(t, mw.CanLocationIngest(ctx))
		require.Error(t, mw.CanLocationIngestSelf(ctx, "subject-3"))
	})
}
