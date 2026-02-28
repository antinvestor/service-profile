package authz_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/frametests"
	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
	"github.com/stretchr/testify/suite"

	"github.com/antinvestor/service-profile/apps/default/service/authz"
	"github.com/antinvestor/service-profile/apps/default/tests/testketo"
)

const (
	testTenantID    = "tenant1"
	testPartitionID = "partition1"
)

func testTenancyPath() string { return testTenantID + "/" + testPartitionID }

type MiddlewareTestSuite struct {
	frametests.FrameBaseTestSuite
	ketoReadURI  string
	ketoWriteURI string
}

func initMiddlewareResources(_ context.Context) []definition.TestResource {
	pg := testpostgres.NewWithOpts("authz_middleware_test",
		definition.WithUserName("ant"),
		definition.WithCredential("s3cr3t"),
		definition.WithEnableLogging(false),
		definition.WithUseHostMode(false),
	)
	keto := testketo.NewWithOpts(
		definition.WithDependancies(pg),
		definition.WithEnableLogging(false),
	)
	return []definition.TestResource{pg, keto}
}

func (s *MiddlewareTestSuite) SetupSuite() {
	s.InitResourceFunc = initMiddlewareResources
	s.FrameBaseTestSuite.SetupSuite()

	ctx := s.T().Context()
	var ketoDep definition.DependancyConn
	for _, res := range s.Resources() {
		if res.Name() == testketo.ImageName {
			ketoDep = res
			break
		}
	}
	s.Require().NotNil(ketoDep, "keto dependency should be available")

	writeURL, err := url.Parse(string(ketoDep.GetDS(ctx)))
	s.Require().NoError(err)
	s.ketoWriteURI = writeURL.Host

	readPort, err := ketoDep.PortMapping(ctx, "4466/tcp")
	s.Require().NoError(err)
	s.ketoReadURI = writeURL.Hostname() + ":" + readPort
}

func (s *MiddlewareTestSuite) newAuthorizer() security.Authorizer {
	cfg := &config.ConfigurationDefault{
		AuthorizationServiceReadURI:  s.ketoReadURI,
		AuthorizationServiceWriteURI: s.ketoWriteURI,
	}
	return authorizer.NewKetoAdapter(cfg, nil)
}

func (s *MiddlewareTestSuite) ctxWithClaims(subjectID string) context.Context {
	claims := &security.AuthenticationClaims{
		TenantID:    testTenantID,
		PartitionID: testPartitionID,
	}
	claims.Subject = subjectID
	return claims.ClaimsToContext(context.Background())
}

func (s *MiddlewareTestSuite) ctxWithSystemInternalClaims(subjectID string) context.Context {
	claims := &security.AuthenticationClaims{
		TenantID:    testTenantID,
		PartitionID: testPartitionID,
		Roles:       []string{"system_internal"},
	}
	claims.Subject = subjectID
	return claims.ClaimsToContext(context.Background())
}

func (s *MiddlewareTestSuite) seedRole(auth security.Authorizer, tenancyPath, profileID, role string) {
	permissions := authz.RolePermissions()[role]
	tuples := make([]security.RelationTuple, 0, 1+len(permissions))

	tuples = append(tuples, security.RelationTuple{
		Object:   security.ObjectRef{Namespace: authz.NamespaceProfile, ID: tenancyPath},
		Relation: role,
		Subject:  security.SubjectRef{Namespace: authz.NamespaceProfileUser, ID: profileID},
	})

	for _, perm := range permissions {
		tuples = append(tuples, security.RelationTuple{
			Object:   security.ObjectRef{Namespace: authz.NamespaceProfile, ID: tenancyPath},
			Relation: perm,
			Subject:  security.SubjectRef{Namespace: authz.NamespaceProfileUser, ID: profileID},
		})
	}

	err := auth.WriteTuples(s.T().Context(), tuples)
	s.Require().NoError(err)
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) TestOwnerHasAllPermissions() {
	auth := s.newAuthorizer()
	s.seedRole(auth, testTenancyPath(), "user1", authz.RoleOwner)

	mw := authz.NewMiddleware(auth)
	ctx := s.ctxWithClaims("user1")

	s.NoError(mw.CanViewProfile(ctx))
	s.NoError(mw.CanCreateProfile(ctx))
	s.NoError(mw.CanUpdateProfile(ctx))
	s.NoError(mw.CanMergeProfiles(ctx))
	s.NoError(mw.CanManageContacts(ctx))
	s.NoError(mw.CanManageRoster(ctx))
	s.NoError(mw.CanManageRelationships(ctx))
}

func (s *MiddlewareTestSuite) TestOperatorPermissions() {
	auth := s.newAuthorizer()
	s.seedRole(auth, testTenancyPath(), "user2", authz.RoleOperator)

	mw := authz.NewMiddleware(auth)
	ctx := s.ctxWithClaims("user2")

	s.NoError(mw.CanViewProfile(ctx))
	s.NoError(mw.CanCreateProfile(ctx))
	s.NoError(mw.CanUpdateProfile(ctx))
	s.NoError(mw.CanManageContacts(ctx))
	s.NoError(mw.CanManageRoster(ctx))

	// Operator cannot merge profiles or manage relationships
	s.Require().Error(mw.CanMergeProfiles(ctx))
	s.Require().Error(mw.CanManageRelationships(ctx))
}

func (s *MiddlewareTestSuite) TestViewerPermissions() {
	auth := s.newAuthorizer()
	s.seedRole(auth, testTenancyPath(), "user3", authz.RoleViewer)

	mw := authz.NewMiddleware(auth)
	ctx := s.ctxWithClaims("user3")

	s.Require().NoError(mw.CanViewProfile(ctx))

	s.Require().Error(mw.CanCreateProfile(ctx))
	s.Require().Error(mw.CanUpdateProfile(ctx))
	s.Require().Error(mw.CanMergeProfiles(ctx))
	s.Require().Error(mw.CanManageContacts(ctx))
	s.Require().Error(mw.CanManageRoster(ctx))
	s.Require().Error(mw.CanManageRelationships(ctx))
}

func (s *MiddlewareTestSuite) TestNoClaims() {
	auth := s.newAuthorizer()
	mw := authz.NewMiddleware(auth)

	err := mw.CanViewProfile(context.Background())
	s.ErrorIs(err, authorizer.ErrInvalidSubject)
}

func (s *MiddlewareTestSuite) TestNoTenant() {
	auth := s.newAuthorizer()
	mw := authz.NewMiddleware(auth)

	claims := &security.AuthenticationClaims{}
	claims.Subject = "user1"
	ctx := claims.ClaimsToContext(context.Background())
	err := mw.CanViewProfile(ctx)
	s.ErrorIs(err, authorizer.ErrInvalidObject)
}

func (s *MiddlewareTestSuite) TestSelfBypass() {
	auth := s.newAuthorizer()
	mw := authz.NewMiddleware(auth)

	// User with no roles can still access their own profile via self-bypass
	ctx := s.ctxWithClaims("self-user")
	s.NoError(mw.CanViewProfileSelf(ctx, "self-user"))
	s.NoError(mw.CanUpdateProfileSelf(ctx, "self-user"))
	s.NoError(mw.CanManageContactsSelf(ctx, "self-user"))

	// Accessing someone else's profile without permission fails
	s.Require().Error(mw.CanViewProfileSelf(ctx, "other-user"))
	s.Require().Error(mw.CanUpdateProfileSelf(ctx, "other-user"))
}

func (s *MiddlewareTestSuite) TestAccessChecker_MemberAllowed() {
	auth := s.newAuthorizer()
	checker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)

	err := auth.WriteTuple(s.T().Context(), authz.BuildAccessTuple(testTenancyPath(), "member-user"))
	s.Require().NoError(err)

	ctx := s.ctxWithClaims("member-user")
	s.NoError(checker.CheckAccess(ctx))
}

func (s *MiddlewareTestSuite) TestAccessChecker_ServiceBotAllowed() {
	auth := s.newAuthorizer()
	checker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)

	err := auth.WriteTuple(s.T().Context(), authz.BuildServiceAccessTuple(testTenancyPath(), "bot-user"))
	s.Require().NoError(err)

	ctx := s.ctxWithSystemInternalClaims("bot-user")
	s.NoError(checker.CheckAccess(ctx))
}

func (s *MiddlewareTestSuite) TestAccessChecker_NoTupleDenied() {
	auth := s.newAuthorizer()
	checker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)

	ctx := s.ctxWithClaims("unknown-user")
	s.Require().Error(checker.CheckAccess(ctx))
}

func (s *MiddlewareTestSuite) seedServiceBridgeTuples(auth security.Authorizer, tenancyPath string) {
	tuples := authz.BuildServiceInheritanceTuples(tenancyPath)
	err := auth.WriteTuples(s.T().Context(), tuples)
	s.Require().NoError(err)
}

func (s *MiddlewareTestSuite) TestServiceBotViaSubjectSets() {
	auth := s.newAuthorizer()
	mw := authz.NewMiddleware(auth)
	accessChecker := authorizer.NewTenancyAccessChecker(auth, authz.NamespaceTenancyAccess)

	s.seedServiceBridgeTuples(auth, testTenancyPath())

	err := auth.WriteTuple(s.T().Context(), authz.BuildServiceAccessTuple(testTenancyPath(), "service-bot"))
	s.Require().NoError(err)

	botCtx := s.ctxWithSystemInternalClaims("service-bot")

	s.NoError(accessChecker.CheckAccess(botCtx))
	s.NoError(mw.CanViewProfile(botCtx))
	s.NoError(mw.CanCreateProfile(botCtx))
	s.NoError(mw.CanManageContacts(botCtx))
}
