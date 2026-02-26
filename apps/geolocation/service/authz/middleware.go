package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the geolocation service.
type Middleware interface {
	CanManageGeolocation(ctx context.Context) error
	CanViewGeolocation(ctx context.Context) error
	CanViewGeolocationSelf(ctx context.Context, targetSubjectID string) error
	CanIngestLocation(ctx context.Context) error
	CanIngestLocationSelf(ctx context.Context, targetSubjectID string) error
}

type middleware struct {
	authorizer security.Authorizer
}

// NewMiddleware creates a new authorisation middleware backed by the given authorizer.
func NewMiddleware(authorizer security.Authorizer) Middleware {
	return &middleware{authorizer: authorizer}
}

// isSelf checks if the target subject is the caller's own profile.
func isSelf(ctx context.Context, targetSubjectID string) bool {
	claims := security.ClaimsFromContext(ctx)
	if claims == nil {
		return false
	}
	subjectID, err := claims.GetSubject()
	if err != nil {
		return false
	}
	return subjectID == targetSubjectID
}

// --- Self-bypass methods ---

func (m *middleware) CanViewGeolocationSelf(ctx context.Context, targetSubjectID string) error {
	if isSelf(ctx, targetSubjectID) {
		return nil
	}
	return m.check(ctx, PermissionViewGeolocation)
}

func (m *middleware) CanIngestLocationSelf(ctx context.Context, targetSubjectID string) error {
	if isSelf(ctx, targetSubjectID) {
		return nil
	}
	return m.check(ctx, PermissionIngestLocation)
}

// --- Non-self methods ---

func (m *middleware) CanManageGeolocation(ctx context.Context) error {
	return m.check(ctx, PermissionManageGeolocation)
}

func (m *middleware) CanViewGeolocation(ctx context.Context) error {
	return m.check(ctx, PermissionViewGeolocation)
}

func (m *middleware) CanIngestLocation(ctx context.Context) error {
	return m.check(ctx, PermissionIngestLocation)
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
