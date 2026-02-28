package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the devices service.
type Middleware interface {
	CanDevicesManage(ctx context.Context) error
	CanDevicesManageSelf(ctx context.Context, targetProfileID string) error
	CanDevicesView(ctx context.Context) error
	CanDevicesViewSelf(ctx context.Context, targetProfileID string) error
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

func (m *middleware) CanDevicesManageSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionDevicesManage)
}

func (m *middleware) CanDevicesViewSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionDevicesView)
}

// --- Non-self methods ---

func (m *middleware) CanDevicesManage(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionDevicesManage)
}

func (m *middleware) CanDevicesView(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionDevicesView)
}
