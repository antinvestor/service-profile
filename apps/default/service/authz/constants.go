package authz

const (
	NamespaceProfile       = "service_profile"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile/user"
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
	RoleMember   = "member"
	RoleService  = "service"
)

// RolePermissions returns the permissions granted by each role.
func RolePermissions() map[string][]string {
	return map[string][]string{
		RoleOwner: {
			PermissionViewProfile, PermissionCreateProfile, PermissionUpdateProfile,
			PermissionMergeProfiles, PermissionManageContacts, PermissionManageRoster,
			PermissionManageRelationships,
		},
		RoleAdmin: {
			PermissionViewProfile, PermissionCreateProfile, PermissionUpdateProfile,
			PermissionMergeProfiles, PermissionManageContacts, PermissionManageRoster,
			PermissionManageRelationships,
		},
		RoleOperator: {
			PermissionViewProfile, PermissionCreateProfile, PermissionUpdateProfile,
			PermissionManageContacts, PermissionManageRoster,
		},
		RoleViewer: {
			PermissionViewProfile,
		},
		RoleMember: {
			PermissionViewProfile,
		},
		RoleService: {
			PermissionViewProfile, PermissionCreateProfile, PermissionUpdateProfile,
			PermissionMergeProfiles, PermissionManageContacts, PermissionManageRoster,
			PermissionManageRelationships,
		},
	}
}
