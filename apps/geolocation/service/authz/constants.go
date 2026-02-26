package authz

const (
	NamespaceTenant  = "profile_tenant"
	NamespaceProfile = "profile"
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
)
