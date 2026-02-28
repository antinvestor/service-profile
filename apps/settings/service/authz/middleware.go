package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the settings service.
type Middleware interface {
	CanManageSettings(ctx context.Context) error
	CanViewSettings(ctx context.Context) error
}

type middleware struct {
	checker *authorizer.FunctionChecker
}

// NewMiddleware creates a new authorisation middleware backed by the given authorizer.
func NewMiddleware(service security.Authorizer) Middleware {
	return &middleware{checker: authorizer.NewFunctionChecker(service, NamespaceProfile)}
}

func (m *middleware) CanManageSettings(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionManageSettings)
}

func (m *middleware) CanViewSettings(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionViewSettings)
}
