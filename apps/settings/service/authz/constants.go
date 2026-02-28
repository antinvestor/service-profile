package authz

const (
	NamespaceProfile       = "service_profile"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile_user"
)

const (
	PermissionSettingsManage = "settings_manage"
	PermissionSettingsView   = "settings_view"
)

const (
	RoleOwner    = "owner"
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
	RoleMember   = "member"
	RoleService  = "service"
)

// GrantedRelation returns the relation name prefixed with "granted_" for use in
// OPL direct grant relations.
func GrantedRelation(permission string) string {
	return "granted_" + permission
}

// RolePermissions returns the permissions granted by each role.
func RolePermissions() map[string][]string {
	return map[string][]string{
		RoleOwner: {
			PermissionSettingsManage, PermissionSettingsView,
		},
		RoleAdmin: {
			PermissionSettingsManage, PermissionSettingsView,
		},
		RoleOperator: {
			PermissionSettingsView,
		},
		RoleViewer: {
			PermissionSettingsView,
		},
		RoleMember: {
			PermissionSettingsView,
		},
		RoleService: {
			PermissionSettingsManage, PermissionSettingsView,
		},
	}
}
