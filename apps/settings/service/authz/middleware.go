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
	authorizer security.Authorizer
}

// NewMiddleware creates a new authorisation middleware backed by the given authorizer.
func NewMiddleware(authorizer security.Authorizer) Middleware {
	return &middleware{authorizer: authorizer}
}

func (m *middleware) CanManageSettings(ctx context.Context) error {
	return m.check(ctx, PermissionManageSettings)
}

func (m *middleware) CanViewSettings(ctx context.Context) error {
	return m.check(ctx, PermissionViewSettings)
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
