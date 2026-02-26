package authz

const (
	NamespaceTenant  = "profile_tenant"
	NamespaceProfile = "profile"
)

const (
	PermissionViewProfile         = "view_profile"
	PermissionCreateProfile       = "create_profile"
	PermissionUpdateProfile       = "update_profile"
	PermissionMergeProfiles       = "merge_profiles"
	PermissionManageContacts      = "manage_contacts"
	PermissionManageRoster        = "manage_roster"
	PermissionManageRelationships = "manage_relationships"
)

const (
	RoleOwner    = "owner"
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)
