package business

import (
	"context"
	"errors"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	settingsv1 "buf.build/gen/go/antinvestor/settingz/protocolbuffers/go/settings/v1"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/frame/workerpool"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

type SettingsBusiness interface {
	Get(context.Context, *settingsv1.GetRequest) (*settingsv1.GetResponse, error)
	List(context.Context, *settingsv1.ListRequest) ([]*settingsv1.SettingObject, error)
	Set(context.Context, *settingsv1.SetRequest) (*settingsv1.SetResponse, error)
	Search(
		ctx context.Context,
		msg *commonv1.SearchRequest,
	) (workerpool.JobResultPipe[[]*settingsv1.SettingObject], error)
}

type settingsBusiness struct {
	refRepo repository.ReferenceRepository
	valRepo repository.SettingValRepository
}

func NewSettingsBusiness(
	refRepo repository.ReferenceRepository,
	valRepo repository.SettingValRepository,
) SettingsBusiness {
	return &settingsBusiness{
		refRepo: refRepo,
		valRepo: valRepo,
	}
}

func (nb *settingsBusiness) Get(ctx context.Context, req *settingsv1.GetRequest) (*settingsv1.GetResponse, error) {
	logger := util.Log(ctx).WithField("request", req)
	logger.Debug("handling get request")

	ref := req.GetKey()
	sRef, err := nb.refRepo.GetByNameAndObjectAndLanguage(ctx, ref.GetModule(),
		ref.GetName(), ref.GetObject(), ref.GetObjectId(), ref.GetLang())
	if err != nil && !data.ErrorIsNoRows(err) {
		logger.WithError(err).Error("could not get settingRef")
		return nil, err
	}
	if sRef != nil {
		sValList, valErr := nb.valRepo.GetByRef(ctx, sRef.GetID())
		if valErr != nil {
			logger.WithError(valErr).Error("could not get settingRef")
			return nil, valErr
		}
		return &settingsv1.GetResponse{
			Data: sValList[0].ToAPI(sRef),
		}, nil
	}
	return &settingsv1.GetResponse{

		Data: &settingsv1.SettingObject{
			Id:    "",
			Key:   req.GetKey(),
			Value: "",
		},
	}, nil
}

func (nb *settingsBusiness) Set(ctx context.Context, req *settingsv1.SetRequest) (*settingsv1.SetResponse, error) {
	logger := util.Log(ctx).WithField("request", req)
	logger.Debug("handling set/update setting")

	ref := req.GetKey()

	sRef, err := nb.refRepo.GetByNameAndObjectAndLanguage(ctx, ref.GetModule(),
		ref.GetName(), ref.GetObject(), ref.GetObjectId(), ref.GetLang())
	if err != nil {
		if !data.ErrorIsNoRows(err) {
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
		err = nb.refRepo.Create(ctx, sRef)
		if err != nil {
			logger.WithError(err).Error("error saving setting ref")
			return nil, err
		}
	}

	var sVal *models.SettingVal
	sValList, err := nb.valRepo.GetByRef(ctx, sRef.GetID())
	if err != nil || len(sValList) == 0 {
		if err != nil && !data.ErrorIsNoRows(err) {
			logger.WithError(err).Error("error querying for setting value")
			return nil, err
		}

		sVal = &models.SettingVal{
			Ref:     sRef.GetID(),
			Version: 0,
		}
	} else {
		sVal = sValList[0]
	}

	sVal.Detail = req.GetValue()
	if sVal.Version == 0 {
		err = nb.valRepo.Create(ctx, sVal)
	} else {
		_, err = nb.valRepo.Update(ctx, sVal, "detail")
	}
	if err != nil {
		logger.WithError(err).Error("error saving setting value")
		return nil, err
	}
	return &settingsv1.SetResponse{
		Data: sVal.ToAPI(sRef),
	}, nil
}

func (nb *settingsBusiness) List(
	ctx context.Context,
	req *settingsv1.ListRequest,
) ([]*settingsv1.SettingObject, error) {
	logger := util.Log(ctx).WithField("request", req)

	logger.Info("handling setting list request")

	setting := req.GetKey()
	settingsList, err := nb.refRepo.SearchRef(ctx, setting.GetModule(),
		setting.GetName(), setting.GetObject(), setting.GetObjectId(), setting.GetLang())
	if err != nil {
		logger.WithError(err).Warn("failed to search for settings")
		return nil, err
	}

	var results []*settingsv1.SettingObject
	var referenceList []string
	sRefMap := map[string]*models.SettingRef{}
	for _, sRef := range settingsList {
		referenceList = append(referenceList, sRef.GetID())
		sRefMap[sRef.GetID()] = sRef
	}

	sValList, getErr := nb.valRepo.GetByRef(ctx, referenceList...)
	if getErr != nil {
		if !data.ErrorIsNoRows(getErr) {
			logger.WithError(getErr).Error("error querying for setting value")
			return results, getErr
		}
	}

	for _, val := range sValList {
		sRef, ok := sRefMap[val.Ref]
		if ok {
			results = append(results, val.ToAPI(sRef))
		}
	}

	return results, nil
}

func (nb *settingsBusiness) Search(
	_ context.Context,
	_ *commonv1.SearchRequest,
) (workerpool.JobResultPipe[[]*settingsv1.SettingObject], error) {
	return nil, errors.New("not implemented")
}
