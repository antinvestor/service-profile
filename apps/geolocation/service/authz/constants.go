package authz

const (
	NamespaceProfile       = "service_profile"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile/user"
)

const (
	PermissionManageGeolocation = "manage_geolocation"
	PermissionViewGeolocation   = "view_geolocation"
	PermissionIngestLocation    = "ingest_location"
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
			PermissionManageGeolocation, PermissionViewGeolocation, PermissionIngestLocation,
		},
		RoleAdmin: {
			PermissionManageGeolocation, PermissionViewGeolocation, PermissionIngestLocation,
		},
		RoleOperator: {
			PermissionViewGeolocation, PermissionIngestLocation,
		},
		RoleViewer: {
			PermissionViewGeolocation,
		},
		RoleMember: {
			PermissionViewGeolocation,
		},
		RoleService: {
			PermissionManageGeolocation, PermissionViewGeolocation, PermissionIngestLocation,
		},
	}
}
