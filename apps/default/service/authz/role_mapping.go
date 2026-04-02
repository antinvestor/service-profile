package authz

import "github.com/pitabwire/frame/security"

const (
	NamespaceProfile       = "service_profile"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile_user"
)

const (
	PermissionProfileView         = "profile_view"
	PermissionProfileCreate       = "profile_create"
	PermissionProfileUpdate       = "profile_update"
	PermissionProfilesMerge       = "profile_merge"
	PermissionContactManage       = "contact_manage"
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
// OPL direct grant relations.
func GrantedRelation(permission string) string {
	return "granted_" + permission
}

// RolePermissions returns the permissions granted by each role.
func RolePermissions() map[string][]string {
	return map[string][]string{
		RoleOwner: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionProfilesMerge, PermissionContactManage, PermissionRosterManage,
			PermissionRelationshipsManage,
		},
		RoleAdmin: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionProfilesMerge, PermissionContactManage, PermissionRosterManage,
			PermissionRelationshipsManage,
		},
		RoleOperator: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionContactManage, PermissionRosterManage,
		},
		RoleViewer: {
			PermissionProfileView,
		},
		RoleMember: {
			PermissionProfileView,
		},
		RoleService: {
			PermissionProfileView, PermissionProfileCreate, PermissionProfileUpdate,
			PermissionProfilesMerge, PermissionContactManage, PermissionRosterManage,
			PermissionRelationshipsManage,
		},
	}
}

// BuildAccessTuple creates a tenancy_access#member tuple for a user.
func BuildAccessTuple(tenancyPath, profileID string) security.RelationTuple {
	return security.RelationTuple{
		Object:   security.ObjectRef{Namespace: NamespaceTenancyAccess, ID: tenancyPath},
		Relation: RoleMember,
		Subject:  security.SubjectRef{Namespace: NamespaceProfileUser, ID: profileID},
	}
}

// BuildServiceAccessTuple creates a tenancy_access#service tuple for a service bot.
func BuildServiceAccessTuple(tenancyPath, profileID string) security.RelationTuple {
	return security.RelationTuple{
		Object:   security.ObjectRef{Namespace: NamespaceTenancyAccess, ID: tenancyPath},
		Relation: RoleService,
		Subject:  security.SubjectRef{Namespace: NamespaceProfileUser, ID: profileID},
	}
}

// BuildServiceInheritanceTuples creates bridge tuples that allow service bots
// to inherit functional permissions via subject sets.
// Only the bridge tuple is needed — the OPL permits already check the service
// role directly, so explicit granted_* tuples per permission are redundant.
func BuildServiceInheritanceTuples(tenancyPath string) []security.RelationTuple {
	return []security.RelationTuple{{
		Object:   security.ObjectRef{Namespace: NamespaceProfile, ID: tenancyPath},
		Relation: RoleService,
		Subject: security.SubjectRef{
			Namespace: NamespaceTenancyAccess,
			ID:        tenancyPath,
			Relation:  RoleService,
		},
	}}
}
