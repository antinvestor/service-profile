package handlers

import (
	"context"

	settingsV1 "github.com/antinvestor/apis/go/settings/v1"
	"github.com/antinvestor/service-profile/apps/settings/service/business"
	"github.com/pitabwire/frame"
)

type SettingsServer struct {
	Service *frame.Service
	settingsV1.UnimplementedSettingsServiceServer
}

func (server *SettingsServer) newSettingsBusiness(ctx context.Context) (business.SettingsBusiness, error) {
	return business.NewSettingsBusiness(ctx, server.Service)
}

// Get a single setting and its stored value.
func (server *SettingsServer) Get(ctx context.Context, req *settingsV1.GetRequest) (*settingsV1.GetResponse, error) {
	notificationBusiness, err := server.newSettingsBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return notificationBusiness.Get(ctx, req)
}

// Set save the setting value appropriately.
func (server *SettingsServer) Set(ctx context.Context, req *settingsV1.SetRequest) (*settingsV1.SetResponse, error) {
	notificationBusiness, err := server.newSettingsBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return notificationBusiness.Set(ctx, req)
}

// List Pulls all setting values that match some criteria in the name & any other setting properties.
func (server *SettingsServer) List(req *settingsV1.ListRequest, stream settingsV1.SettingsService_ListServer) error {
	settingsBusiness, err := server.newSettingsBusiness(stream.Context())
	if err != nil {
		return err
	}
	return settingsBusiness.List(req, stream)
}
