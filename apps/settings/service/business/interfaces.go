package business

import (
	"context"

	settingsV1 "github.com/antinvestor/apis/go/settings/v1"
	"github.com/pitabwire/frame"
)

type SettingsBusiness interface {
	Get(context.Context, *settingsV1.GetRequest) (*settingsV1.GetResponse, error)
	List(*settingsV1.ListRequest, settingsV1.SettingsService_ListServer) error
	Set(context.Context, *settingsV1.SetRequest) (*settingsV1.SetResponse, error)
}

func NewSettingsBusiness(_ context.Context, service *frame.Service) (SettingsBusiness, error) {
	if service == nil {
		return nil, ErrInitializationFail
	}

	return &settingsBusiness{
		service: service,
	}, nil
}
