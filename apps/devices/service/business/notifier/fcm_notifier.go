package notifier

import (
	"context"
	"errors"
	"strings"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/pitabwire/util"
)

const (
	fcmProvider            = "fcm"
	defaultFCMMaxBatchSize = 500
)

type fcmNotifier struct {
	cfg *config.DevicesConfig

	client *messaging.Client
}

func NewFCMNotifier(ctx context.Context, cfg *config.DevicesConfig) (Notifier, error) {

	firebaseApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	messagingCli, err := firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &fcmNotifier{
		cfg:    cfg,
		client: messagingCli,
	}, nil
}

func (f *fcmNotifier) Register(ctx context.Context, req *devicev1.RegisterKeyRequest) (*devicev1.KeyObject, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	return nil, nil
}

func (f *fcmNotifier) DeRegister(_ context.Context, _ *devicev1.KeyObject) error {
	return nil
}

func (f *fcmNotifier) Notify(ctx context.Context, req *devicev1.NotifyRequest, keys ...*devicev1.KeyObject) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	for _, key := range keys {
		if key == nil || key.KeyType != devicev1.KeyType_FCM_TOKEN {
			continue
		}

		// Create a list containing up to 500 messages.
		messageBatches := f.getMessageBatches(ctx, key, req)

		for _, messages := range messageBatches {
			br, err := f.client.SendEach(ctx, messages)
			if err != nil {
				return err
			}

			util.Log(ctx).
				WithField("response", br).
				Debug("notification batch sent via FCM")
		}
	}

	return nil
}

func (f *fcmNotifier) batchMaxSize() int {
	if f.cfg != nil && f.cfg.FCMMaxBatchSize > 0 {
		return f.cfg.FCMMaxBatchSize
	}
	return defaultFCMMaxBatchSize
}

func (f *fcmNotifier) getMessageBatches(_ context.Context, key *devicev1.KeyObject, req *devicev1.NotifyRequest) [][]*messaging.Message {

	var responseSlice [][]*messaging.Message

	registrationToken := strings.TrimSpace(string(key.GetKey()))

	data := make(map[string]string)

	for k, v := range req.GetData().GetFields() {
		data[k] = v.GetStringValue()
	}

	responseSlice = append(responseSlice, []*messaging.Message{

		{
			Data: data,
			Notification: &messaging.Notification{
				Title:    req.GetTitle(),
				Body:     req.GetBody(),
				ImageURL: "",
			},
			Token: registrationToken,
		},
	})

	return responseSlice

}
