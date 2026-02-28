package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the settings service.
type Middleware interface {
	CanSettingsManage(ctx context.Context) error
	CanSettingsView(ctx context.Context) error
}

type middleware struct {
	checker *authorizer.FunctionChecker
}

// NewMiddleware creates a new authorisation middleware backed by the given authorizer.
func NewMiddleware(service security.Authorizer) Middleware {
	return &middleware{checker: authorizer.NewFunctionChecker(service, NamespaceProfile)}
}

func (m *middleware) CanSettingsManage(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionSettingsManage)
}

func (m *middleware) CanSettingsView(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionSettingsView)
}
