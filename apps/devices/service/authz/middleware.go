package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the devices service.
type Middleware interface {
	CanManageDevices(ctx context.Context) error
	CanManageDevicesSelf(ctx context.Context, targetProfileID string) error
	CanViewDevices(ctx context.Context) error
	CanViewDevicesSelf(ctx context.Context, targetProfileID string) error
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

func (m *middleware) CanManageDevicesSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionManageDevices)
}

func (m *middleware) CanViewDevicesSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionViewDevices)
}

// --- Non-self methods ---

func (m *middleware) CanManageDevices(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionManageDevices)
}

func (m *middleware) CanViewDevices(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionViewDevices)
}
