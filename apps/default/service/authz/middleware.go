package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the profile service.
type Middleware interface {
	CanViewProfile(ctx context.Context) error
	CanViewProfileSelf(ctx context.Context, targetProfileID string) error
	CanCreateProfile(ctx context.Context) error
	CanUpdateProfile(ctx context.Context) error
	CanUpdateProfileSelf(ctx context.Context, targetProfileID string) error
	CanMergeProfiles(ctx context.Context) error
	CanManageContacts(ctx context.Context) error
	CanManageContactsSelf(ctx context.Context, targetProfileID string) error
	CanManageRoster(ctx context.Context) error
	CanManageRosterSelf(ctx context.Context, targetProfileID string) error
	CanManageRelationships(ctx context.Context) error
	CanManageRelationshipsSelf(ctx context.Context, targetProfileID string) error
}

type middleware struct {
	authorizer security.Authorizer
}

// NewMiddleware creates a new authorisation middleware backed by the given authorizer.
func NewMiddleware(authorizer security.Authorizer) Middleware {
	return &middleware{authorizer: authorizer}
}

// isSelf checks if the target profile is the caller's own profile.
func isSelf(ctx context.Context, targetProfileID string) bool {
	claims := security.ClaimsFromContext(ctx)
	if claims == nil {
		return false
	}
	subjectID, err := claims.GetSubject()
	if err != nil {
		return false
	}
	return subjectID == targetProfileID
}

// --- Self-bypass methods ---

func (m *middleware) CanViewProfileSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.check(ctx, PermissionViewProfile)
}

func (m *middleware) CanUpdateProfileSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.check(ctx, PermissionUpdateProfile)
}

func (m *middleware) CanManageContactsSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.check(ctx, PermissionManageContacts)
}

func (m *middleware) CanManageRosterSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.check(ctx, PermissionManageRoster)
}

func (m *middleware) CanManageRelationshipsSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.check(ctx, PermissionManageRelationships)
}

// --- Non-self methods ---

func (m *middleware) CanViewProfile(ctx context.Context) error {
	return m.check(ctx, PermissionViewProfile)
}

func (m *middleware) CanCreateProfile(ctx context.Context) error {
	return m.check(ctx, PermissionCreateProfile)
}

func (m *middleware) CanUpdateProfile(ctx context.Context) error {
	return m.check(ctx, PermissionUpdateProfile)
}

func (m *middleware) CanMergeProfiles(ctx context.Context) error {
	return m.check(ctx, PermissionMergeProfiles)
}

func (m *middleware) CanManageContacts(ctx context.Context) error {
	return m.check(ctx, PermissionManageContacts)
}

func (m *middleware) CanManageRoster(ctx context.Context) error {
	return m.check(ctx, PermissionManageRoster)
}

func (m *middleware) CanManageRelationships(ctx context.Context) error {
	return m.check(ctx, PermissionManageRelationships)
}

// check performs the Keto permission check against the tenant namespace.
func (m *middleware) check(ctx context.Context, permission string) error {
	claims := security.ClaimsFromContext(ctx)
	if claims == nil {
		return authorizer.ErrInvalidSubject
	}

	subjectID, err := claims.GetSubject()
	if err != nil || subjectID == "" {
		return authorizer.ErrInvalidSubject
	}

	tenantID := claims.GetTenantID()
	if tenantID == "" {
		return authorizer.ErrInvalidObject
	}

	req := security.CheckRequest{
		Object:     security.ObjectRef{Namespace: NamespaceTenant, ID: tenantID},
		Permission: permission,
		Subject:    security.SubjectRef{Namespace: NamespaceProfile, ID: subjectID},
	}

	result, err := m.authorizer.Check(ctx, req)
	if err != nil {
		return err
	}
	if !result.Allowed {
		return authorizer.NewPermissionDeniedError(req.Object, permission, req.Subject, result.Reason)
	}

	return nil
}
