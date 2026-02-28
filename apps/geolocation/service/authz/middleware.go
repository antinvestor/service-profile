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
	checker *authorizer.FunctionChecker
}

// NewMiddleware creates a new authorisation middleware backed by the given authorizer.
func NewMiddleware(service security.Authorizer) Middleware {
	return &middleware{checker: authorizer.NewFunctionChecker(service, NamespaceProfile)}
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
	return m.checker.Check(ctx, PermissionViewGeolocation)
}

func (m *middleware) CanIngestLocationSelf(ctx context.Context, targetSubjectID string) error {
	if isSelf(ctx, targetSubjectID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionIngestLocation)
}

// --- Non-self methods ---

func (m *middleware) CanManageGeolocation(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionManageGeolocation)
}

func (m *middleware) CanViewGeolocation(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionViewGeolocation)
}

func (m *middleware) CanIngestLocation(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionIngestLocation)
}
