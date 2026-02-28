package authz

const (
	NamespaceProfile       = "service_profile"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile_user"
)

const (
	PermissionProfileView         = "profile_view"
	PermissionProfileCreate       = "profile_create"
	PermissionProfileUpdate       = "profile_update"
	PermissionProfilesMerge       = "profiles_merge"
	PermissionContactsManage      = "contacts_manage"
	PermissionRosterManage        = "roster_manage"
	PermissionRelationshipsManage = "relationships_manage"
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
// OPL direct grant relations. This avoids name conflicts where Keto skips permit
// evaluation when a relation shares the same name as a permit function.
func GrantedRelation(permission string) string {
	return "granted_" + permission
}

// RolePermissions returns the permissions granted by each role.
func RolePermissions() map[string][]string {
	return map[string][]string{
		RoleOwner: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionProfilesMerge, PermissionContactsManage, PermissionRosterManage,
			PermissionRelationshipsManage,
		},
		RoleAdmin: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionProfilesMerge, PermissionContactsManage, PermissionRosterManage,
			PermissionRelationshipsManage,
		},
		RoleOperator: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionContactsManage, PermissionRosterManage,
		},
		RoleViewer: {
			PermissionProfileView,
		},
		RoleMember: {
			PermissionProfileView,
		},
		RoleService: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionProfilesMerge, PermissionContactsManage, PermissionRosterManage,
			PermissionRelationshipsManage,
		},
	}
}
