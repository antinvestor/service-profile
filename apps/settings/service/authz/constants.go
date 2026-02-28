package authz

const (
	NamespaceProfile       = "service_profile"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile/user"
)

const (
	PermissionManageSettings = "manage_settings"
	PermissionViewSettings   = "view_settings"
)

const (
	RoleOwner    = "owner"
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
	RoleMember   = "member"
	RoleService  = "service"
)

// RolePermissions returns the permissions granted by each role.
func RolePermissions() map[string][]string {
	return map[string][]string{
		RoleOwner: {
			PermissionManageSettings, PermissionViewSettings,
		},
		RoleAdmin: {
			PermissionManageSettings, PermissionViewSettings,
		},
		RoleOperator: {
			PermissionViewSettings,
		},
		RoleViewer: {
			PermissionViewSettings,
		},
		RoleMember: {
			PermissionViewSettings,
		},
		RoleService: {
			PermissionManageSettings, PermissionViewSettings,
		},
	}
}
