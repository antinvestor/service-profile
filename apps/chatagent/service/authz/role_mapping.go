package authz

import "github.com/pitabwire/frame/v2/security"

const (
	NamespaceChatAgent     = "service_chat_agent"
	NamespaceTenancyAccess = "tenancy_access"
	NamespaceProfileUser   = "profile_user"
)

const (
	PermissionChatAgentView   = "chat_agent_view"
	PermissionChatAgentManage = "chat_agent_manage"
	PermissionChatAgentTurn   = "chat_agent_turn"
)

const (
	RoleOwner    = "owner"
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
	RoleMember   = "member"
	RoleService  = "service"
)

// GrantedRelation returns the relation name prefixed with "granted_".
func GrantedRelation(permission string) string {
	return "granted_" + permission
}

// RolePermissions returns the permissions granted by each role.
func RolePermissions() map[string][]string {
	return map[string][]string{
		RoleOwner: {
			PermissionChatAgentView, PermissionChatAgentManage, PermissionChatAgentTurn,
		},
		RoleAdmin: {
			PermissionChatAgentView, PermissionChatAgentManage, PermissionChatAgentTurn,
		},
		RoleOperator: {
			PermissionChatAgentView, PermissionChatAgentTurn,
		},
		RoleViewer: {
			PermissionChatAgentView,
		},
		RoleMember: {
			PermissionChatAgentView, PermissionChatAgentTurn,
		},
		RoleService: {
			PermissionChatAgentView, PermissionChatAgentManage, PermissionChatAgentTurn,
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

// BuildServiceInheritanceTuples bridges service bots to chat-agent permissions.
func BuildServiceInheritanceTuples(tenancyPath string) []security.RelationTuple {
	return []security.RelationTuple{{
		Object:   security.ObjectRef{Namespace: NamespaceChatAgent, ID: tenancyPath},
		Relation: RoleService,
		Subject: security.SubjectRef{
			Namespace: NamespaceTenancyAccess,
			ID:        tenancyPath,
			Relation:  RoleService,
		},
	}}
}
