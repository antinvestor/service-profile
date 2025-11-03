package business

import (
	"context"
	"errors"
	"fmt"
	"strings"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/frame/workerpool"

	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/business/notifier"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
)

type NotifyBusiness interface {
	RegisterKey(ctx context.Context, req *devicev1.RegisterKeyRequest) (*devicev1.KeyObject, error)
	DeRegisterKey(ctx context.Context, req *devicev1.DeRegisterKeyRequest) error
	Notify(ctx context.Context, req *devicev1.NotifyRequest) ([]*devicev1.NotifyResult, error)
}

type notifyBusiness struct {
	cfg *config.DevicesConfig

	qMan    queue.Manager
	workMan workerpool.Manager

	keysBusiness KeysBusiness
	deviceRepo   repository.DeviceRepository
	notifiers    map[devicev1.KeyType]notifier.Notifier
}

// NewNotifyBusiness creates a new instance of NotificationBusiness.
func NewNotifyBusiness(
	ctx context.Context,
	cfg *config.DevicesConfig,
	qMan queue.Manager,
	workMan workerpool.Manager,
	keyBusiness KeysBusiness,
	deviceRepo repository.DeviceRepository,
) (NotifyBusiness, error) {
	n := &notifyBusiness{
		cfg:     cfg,
		qMan:    qMan,
		workMan: workMan,

		keysBusiness: keyBusiness,
		deviceRepo:   deviceRepo,
	}

	fcmNotifier, err := notifier.NewFCMNotifier(ctx, cfg)
	if err != nil {
		return nil, err
	}

	n.notifiers = map[devicev1.KeyType]notifier.Notifier{
		devicev1.KeyType_FCM_TOKEN: fcmNotifier,
	}

	return n, nil
}

func (n notifyBusiness) RegisterKey(
	ctx context.Context,
	req *devicev1.RegisterKeyRequest,
) (*devicev1.KeyObject, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	deviceID := strings.TrimSpace(req.GetDeviceId())
	if deviceID == "" {
		return nil, errors.New("device id is required")
	}

	if _, err := n.deviceRepo.GetByID(ctx, deviceID); err != nil {
		return nil, err
	}

	notifyHandler, err := n.notifierFor(req.GetKeyType())
	if err != nil {
		return nil, err
	}

	return notifyHandler.Register(ctx, req)
}

func (n notifyBusiness) DeRegisterKey(ctx context.Context, req *devicev1.DeRegisterKeyRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	keyID := strings.TrimSpace(req.GetId())
	if keyID == "" {
		return errors.New("key id is required")
	}

	resultCh, err := n.keysBusiness.RemoveKeys(ctx, keyID)
	if err != nil {
		return err
	}

	var removedKeys []*devicev1.KeyObject
	for res := range resultCh {
		if res.IsError() {
			return res.Error()
		}
		removedKeys = append(removedKeys, res.Item()...)
	}

	for _, key := range removedKeys {
		if key == nil {
			continue
		}

		notifyHandler, notifyErr := n.notifierFor(key.GetKeyType())
		if notifyErr != nil {
			return notifyErr
		}

		if notifyErr = notifyHandler.DeRegister(ctx, key); notifyErr != nil {
			return notifyErr
		}
	}

	return nil
}

func (n notifyBusiness) getActiveDeviceKey(
	ctx context.Context,
	deviceID string,
	requestedType devicev1.KeyType,
	keyID string,
) (map[devicev1.KeyType][]*devicev1.KeyObject, error) {
	keysCh, err := n.keysBusiness.GetKeys(ctx, deviceID, requestedType)
	if err != nil {
		return nil, err
	}

	targetKeyID := strings.TrimSpace(keyID)
	var matchedKeys []*devicev1.KeyObject
	keyGroups := map[devicev1.KeyType][]*devicev1.KeyObject{}

	for res := range keysCh {
		if res.IsError() {
			return nil, res.Error()
		}

		for _, key := range res.Item() {
			if targetKeyID != "" && key.GetId() != targetKeyID {
				continue
			}

			if requestedType != devicev1.KeyType(0) && key.GetKeyType() != requestedType {
				continue
			}

			matchedKeys = append(matchedKeys, key)
			keyType := key.GetKeyType()
			keyGroups[keyType] = append(keyGroups[keyType], key)
		}
	}

	if len(matchedKeys) == 0 {
		if targetKeyID != "" {
			return nil, fmt.Errorf("no notification token found for key id %s", targetKeyID)
		}
		return nil, errors.New("no notification tokens registered for device")
	}

	if requestedType == devicev1.KeyType(0) && len(keyGroups) > 1 {
		return nil, errors.New("multiple key types matched notification request; specify key_type")
	}

	return keyGroups, nil
}

func (n notifyBusiness) Notify(ctx context.Context, req *devicev1.NotifyRequest) ([]*devicev1.NotifyResult, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	deviceID := strings.TrimSpace(req.GetDeviceId())
	if deviceID == "" {
		return nil, errors.New("device id is required")
	}

	if _, err := n.deviceRepo.GetByID(ctx, deviceID); err != nil {
		return nil, err
	}

	requestedType := req.GetKeyType()

	keyGroups, err := n.getActiveDeviceKey(ctx, deviceID, requestedType, req.GetKeyId())
	if err != nil {
		return nil, err
	}

	var allResponses []*devicev1.NotifyResult

	for keyType, keys := range keyGroups {
		notifyHandler, notifyErr := n.notifierFor(keyType)
		if notifyErr != nil {
			return nil, notifyErr
		}

		response, notifyErr := notifyHandler.Notify(ctx, req, keys...)
		if notifyErr != nil {
			return nil, notifyErr
		}

		allResponses = append(allResponses, response...)
	}

	return allResponses, nil
}

func (n notifyBusiness) notifierFor(keyType devicev1.KeyType) (notifier.Notifier, error) {
	notifyHandler, ok := n.notifiers[keyType]
	if !ok {
		return nil, fmt.Errorf("no Notifier configured for key type %s", keyType.String())
	}
	return notifyHandler, nil
}
