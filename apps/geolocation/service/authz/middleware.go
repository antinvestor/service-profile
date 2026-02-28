package authz

import (
	"context"

	"github.com/pitabwire/frame/security"
	"github.com/pitabwire/frame/security/authorizer"
)

// Middleware defines authorisation checks for the geolocation service.
type Middleware interface {
	CanGeolocationManage(ctx context.Context) error
	CanGeolocationView(ctx context.Context) error
	CanGeolocationViewSelf(ctx context.Context, targetSubjectID string) error
	CanLocationIngest(ctx context.Context) error
	CanLocationIngestSelf(ctx context.Context, targetSubjectID string) error
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

func (m *middleware) CanGeolocationViewSelf(ctx context.Context, targetSubjectID string) error {
	if isSelf(ctx, targetSubjectID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionGeolocationView)
}

func (m *middleware) CanLocationIngestSelf(ctx context.Context, targetSubjectID string) error {
	if isSelf(ctx, targetSubjectID) {
		return nil
	}
	return m.checker.Check(ctx, PermissionLocationIngest)
}

// --- Non-self methods ---

func (m *middleware) CanGeolocationManage(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionGeolocationManage)
}

func (m *middleware) CanGeolocationView(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionGeolocationView)
}

func (m *middleware) CanLocationIngest(ctx context.Context) error {
	return m.checker.Check(ctx, PermissionLocationIngest)
}
