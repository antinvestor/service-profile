package authz

const (
	NamespaceProfile       = "service_profile"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile/user"
)

const (
	PermissionManageDevices = "manage_devices"
	PermissionViewDevices   = "view_devices"
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
			PermissionManageDevices, PermissionViewDevices,
		},
		RoleAdmin: {
			PermissionManageDevices, PermissionViewDevices,
		},
		RoleOperator: {
			PermissionViewDevices,
		},
		RoleViewer: {
			PermissionViewDevices,
		},
		RoleMember: {
			PermissionViewDevices,
		},
		RoleService: {
			PermissionManageDevices, PermissionViewDevices,
		},
	}
}
