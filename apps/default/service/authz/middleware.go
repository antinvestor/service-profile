package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the profile service.
type Middleware interface {
	CanProfileView(ctx context.Context) error
	CanProfileViewSelf(ctx context.Context, targetProfileID string) error
	CanProfileCreate(ctx context.Context) error
	CanProfileUpdate(ctx context.Context) error
	CanProfileUpdateSelf(ctx context.Context, targetProfileID string) error
	CanProfilesMerge(ctx context.Context) error
	CanContactsManage(ctx context.Context) error
	CanContactsManageSelf(ctx context.Context, targetProfileID string) error
	CanRosterManage(ctx context.Context) error
	CanRosterManageSelf(ctx context.Context, targetProfileID string) error
	CanRelationshipsManage(ctx context.Context) error
	CanRelationshipsManageSelf(ctx context.Context, targetProfileID string) error
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

func (m *middleware) CanProfileViewSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionProfileView)
}

func (m *middleware) CanProfileUpdateSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionProfileUpdate)
}

func (m *middleware) CanContactsManageSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionContactsManage)
}

func (m *middleware) CanRosterManageSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionRosterManage)
}

func (m *middleware) CanRelationshipsManageSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionRelationshipsManage)
}

// --- Non-self methods ---

func (m *middleware) CanProfileView(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionProfileView)
}

func (m *middleware) CanProfileCreate(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionProfileCreate)
}

func (m *middleware) CanProfileUpdate(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionProfileUpdate)
}

func (m *middleware) CanProfilesMerge(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionProfilesMerge)
}

func (m *middleware) CanContactsManage(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionContactsManage)
}

func (m *middleware) CanRosterManage(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionRosterManage)
}

func (m *middleware) CanRelationshipsManage(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionRelationshipsManage)
}
