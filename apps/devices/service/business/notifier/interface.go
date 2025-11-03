package notifier

import (
	"context"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
)

type Notifier interface {
	Register(ctx context.Context, req *devicev1.RegisterKeyRequest) (*devicev1.KeyObject, error)
	DeRegister(ctx context.Context, key *devicev1.KeyObject) error
	Notify(
		ctx context.Context,
		req *devicev1.NotifyRequest,
		keys ...*devicev1.KeyObject,
	) ([]*devicev1.NotifyResult, error)
}
