package authz //nolint:testpackage // tests access unexported authz internals

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrantedRelationAndRolePermissions(t *testing.T) {
	t.Parallel()

	require.Equal(t, "granted_geolocation_view", GrantedRelation(PermissionGeolocationView))

	perms := RolePermissions()
	require.Contains(t, perms[RoleOwner], PermissionGeolocationManage)
	require.Contains(t, perms[RoleViewer], PermissionGeolocationView)
	require.NotContains(t, perms[RoleViewer], PermissionLocationIngest)
}

func TestTupleBuilders(t *testing.T) {
	t.Parallel()

	access := BuildAccessTuple("tenant/partition", "profile-1")
	require.Equal(t, NamespaceTenancyAccess, access.Object.Namespace)
	require.Equal(t, RoleMember, access.Relation)
	require.Equal(t, NamespaceProfileUser, access.Subject.Namespace)

	service := BuildServiceAccessTuple("tenant/partition", "svc-profile")
	require.Equal(t, RoleService, service.Relation)

	inheritance := BuildServiceInheritanceTuples("tenant/partition")
	require.Len(t, inheritance, 1)
	require.Equal(t, NamespaceProfile, inheritance[0].Object.Namespace)
	require.Equal(t, RoleService, inheritance[0].Relation)
}
