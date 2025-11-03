package notifier

import (
	"context"
	"errors"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/devices/config"
)

const (
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

func (f *fcmNotifier) Notify(
	ctx context.Context,
	req *devicev1.NotifyRequest,
	keys ...*devicev1.KeyObject,
) ([]*devicev1.NotifyResult, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	var responses []*devicev1.NotifyResult

	for _, key := range keys {
		if key == nil || key.GetKeyType() != devicev1.KeyType_FCM_TOKEN {
			continue
		}

		// Create a list containing up to 500 messages.
		var messages []*messaging.Message

		for _, message := range req.GetNotifications() {
			messages = append(messages, f.toFCMMessage(ctx, key, message))

			if len(messages) >= f.batchMaxSize() {
				br, err := f.client.SendEach(ctx, messages)
				if err != nil {
					return nil, err
				}
				util.Log(ctx).
					WithField("response", br).
					Debug("notification batch sent via FCM")
				messages = []*messaging.Message{}

				responses = append(responses, f.toNotifyResult(br)...)
			}
		}

		br, err := f.client.SendEach(ctx, messages)
		if err != nil {
			return nil, err
		}
		responses = append(responses, f.toNotifyResult(br)...)
		util.Log(ctx).
			WithField("response", br).
			Debug("notification batch sent via FCM")
	}

	return responses, nil
}

func (f *fcmNotifier) batchMaxSize() int {
	if f.cfg != nil && f.cfg.FCMMaxBatchSize > 0 {
		return f.cfg.FCMMaxBatchSize
	}
	return defaultFCMMaxBatchSize
}

func (f *fcmNotifier) toFCMMessage(
	_ context.Context,
	key *devicev1.KeyObject,
	req *devicev1.NotifyMessage,
) *messaging.Message {
	registrationToken := strings.TrimSpace(string(key.GetKey()))

	data := make(map[string]string)

	for k, v := range req.GetData().GetFields() {
		data[k] = v.GetStringValue()
	}

	return &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{

			Title:    req.GetTitle(),
			Body:     req.GetBody(),
			ImageURL: "",
		},
		Token: registrationToken,
	}
}

func (f *fcmNotifier) toNotifyResult(br *messaging.BatchResponse) []*devicev1.NotifyResult {
	var response []*devicev1.NotifyResult

	for _, resp := range br.Responses {
		message := "ok"
		if resp.Error != nil {
			message = resp.Error.Error()
		}

		response = append(response, &devicev1.NotifyResult{
			Success:        resp.Success,
			Message:        message,
			NotificationId: resp.MessageID,
		})
	}

	return response
}
