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
	checker *authorizer.FunctionChecker
}

// NewMiddleware creates a new authorisation middleware backed by the given authorizer.
func NewMiddleware(service security.Authorizer) Middleware {
	return &middleware{checker: authorizer.NewFunctionChecker(service, NamespaceProfile)}
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
	return m.checker.Check(ctx, PermissionViewProfile)
}

func (m *middleware) CanUpdateProfileSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionUpdateProfile)
}

func (m *middleware) CanManageContactsSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionManageContacts)
}

func (m *middleware) CanManageRosterSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionManageRoster)
}

func (m *middleware) CanManageRelationshipsSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionManageRelationships)
}

// --- Non-self methods ---

func (m *middleware) CanViewProfile(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionViewProfile)
}

func (m *middleware) CanCreateProfile(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionCreateProfile)
}

func (m *middleware) CanUpdateProfile(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionUpdateProfile)
}

func (m *middleware) CanMergeProfiles(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionMergeProfiles)
}

func (m *middleware) CanManageContacts(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionManageContacts)
}

func (m *middleware) CanManageRoster(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionManageRoster)
}

func (m *middleware) CanManageRelationships(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionManageRelationships)
}
