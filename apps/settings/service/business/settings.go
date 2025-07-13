package business

import (
	"context"

	settingsV1 "github.com/antinvestor/apis/go/settings/v1"
	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

type settingsBusiness struct {
	service *frame.Service
}

func (nb *settingsBusiness) Get(ctx context.Context, req *settingsV1.GetRequest) (*settingsV1.GetResponse, error) {
	logger := nb.service.Log(ctx).WithField("request", req)
	logger.Debug("handling get request")

	referenceRepo := repository.NewReferenceRepository(ctx, nb.service)
	ref := req.GetKey()
	sRef, err := referenceRepo.GetByNameAndObjectAndLanguage(ctx, ref.GetModule(),
		ref.GetName(), ref.GetObject(), ref.GetObjectId(), ref.GetLang())
	if err != nil && !frame.ErrorIsNoRows(err) {
		logger.WithError(err).Error("could not get settingRef")
		return nil, err
	}
	if sRef != nil {
		valRepo := repository.NewSettingValRepository(ctx, nb.service)
		sVal, valErr := valRepo.GetByRef(ctx, sRef.GetID())
		if valErr != nil {
			logger.WithError(valErr).Error("could not get settingRef")
			return nil, valErr
		}
		return &settingsV1.GetResponse{
			Data: sVal.ToAPI(sRef),
		}, nil
	}
	return &settingsV1.GetResponse{

		Data: &settingsV1.SettingObject{
			Id:    "",
			Key:   req.GetKey(),
			Value: "",
		},
	}, nil
}

func (nb *settingsBusiness) Set(ctx context.Context, req *settingsV1.SetRequest) (*settingsV1.SetResponse, error) {
	logger := nb.service.Log(ctx).WithField("request", req)
	logger.Debug("handling set/update setting")

	referenceRepo := repository.NewReferenceRepository(ctx, nb.service)
	ref := req.GetKey()

	sRef, err := referenceRepo.GetByNameAndObjectAndLanguage(ctx, ref.GetModule(),
		ref.GetName(), ref.GetObject(), ref.GetObjectId(), ref.GetLang())
	if err != nil {
		if !frame.ErrorIsNoRows(err) {
			logger.WithError(err).Error("error querying for setting ref")
			return nil, err
		}

		sRef = &models.SettingRef{
			Name:     ref.GetName(),
			Object:   ref.GetObject(),
			ObjectID: ref.GetObjectId(),
			Language: ref.GetLang(),
			Module:   ref.GetModule(),
		}
		err = referenceRepo.Save(ctx, sRef)
		if err != nil {
			logger.WithError(err).Error("error saving setting ref")
			return nil, err
		}
	}

	valRepo := repository.NewSettingValRepository(ctx, nb.service)
	sVal, err := valRepo.GetByRef(ctx, sRef.GetID())
	if err != nil {
		if !frame.ErrorIsNoRows(err) {
			logger.WithError(err).Error("error querying for setting value")
			return nil, err
		}

		sVal = &models.SettingVal{
			Ref:     sRef.GetID(),
			Version: 0,
		}
	}

	sVal.Detail = req.GetValue()
	sVal.Version++
	err = valRepo.Save(ctx, sVal)
	if err != nil {
		logger.WithError(err).Error("error saving setting value")
		return nil, err
	}
	return &settingsV1.SetResponse{
		Data: sVal.ToAPI(sRef),
	}, nil
}

func (nb *settingsBusiness) List(req *settingsV1.ListRequest, stream settingsV1.SettingsService_ListServer) error {
	ctx := stream.Context()

	logger := nb.service.Log(ctx).WithField("request", req)

	logger.Info("handling setting list request")

	referenceRepo := repository.NewReferenceRepository(ctx, nb.service)

	setting := req.GetKey()
	settingsList, err := referenceRepo.Search(ctx, setting.GetModule(),
		setting.GetName(), setting.GetObject(), setting.GetObjectId(), setting.GetLang())
	if err != nil {
		logger.WithError(err).Warn("failed to search for settings")
		return err
	}

	valRepo := repository.NewSettingValRepository(ctx, nb.service)

	var results []*settingsV1.SettingObject
	for _, sRef := range settingsList {
		sVal, getErr := valRepo.GetByRef(ctx, sRef.GetID())
		if getErr != nil {
			if !frame.ErrorIsNoRows(getErr) {
				logger.WithError(getErr).Error("error querying for setting value")
				return getErr
			}

			sVal = &models.SettingVal{
				Ref:     sRef.GetID(),
				Version: 0,
			}
		}

		results = append(results, sVal.ToAPI(sRef))
	}

	err = stream.Send(&settingsV1.ListResponse{
		Data: results})
	if err != nil {
		logger.WithError(err).Warn(" unable to send a result")
		return err
	}

	return nil
}
