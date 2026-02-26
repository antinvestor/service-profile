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

func (m *middleware) CanManageDevicesSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.check(ctx, PermissionManageDevices)
}

func (m *middleware) CanViewDevicesSelf(ctx context.Context, targetProfileID string) error {
	if isSelf(ctx, targetProfileID) {
		return nil
	}
	return m.check(ctx, PermissionViewDevices)
}

// --- Non-self methods ---

func (m *middleware) CanManageDevices(ctx context.Context) error {
	return m.check(ctx, PermissionManageDevices)
}

func (m *middleware) CanViewDevices(ctx context.Context) error {
	return m.check(ctx, PermissionViewDevices)
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
