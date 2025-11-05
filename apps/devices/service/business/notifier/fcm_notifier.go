package notifier

import (
	"context"
	"errors"
	"strings"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
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

func (f *fcmNotifier) Register(_ context.Context, req *devicev1.RegisterKeyRequest) (*devicev1.KeyObject, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// FCM tokens are registered directly by the client and validated on notify.
	// The actual token storage is handled by the keys business layer.
	// This method simply validates the request and returns a KeyObject.
	return &devicev1.KeyObject{
		KeyType: devicev1.KeyType_FCM_TOKEN,
	}, nil
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

	if len(keys) == 0 {
		return []*devicev1.NotifyResult{}, nil
	}

	// Pre-allocate responses slice with estimated capacity for better performance
	responses := make([]*devicev1.NotifyResult, 0, len(keys)*len(req.GetNotifications()))
	notifications := req.GetNotifications()
	batchSize := f.batchMaxSize()

	for _, key := range keys {
		if key == nil || key.GetKeyType() != devicev1.KeyType_FCM_TOKEN {
			continue
		}

		// Process notifications in batches
		messages := make([]*messaging.Message, 0, batchSize)

		for _, message := range notifications {
			messages = append(messages, f.toFCMMessage(ctx, key, message))

			// Send batch when it reaches max size
			if len(messages) >= batchSize {
				batchResponses, err := f.sendBatch(ctx, messages)
				if err != nil {
					util.Log(ctx).WithError(err).Error("failed to send FCM batch")
					return nil, err
				}
				responses = append(responses, batchResponses...)
				messages = make([]*messaging.Message, 0, batchSize)
			}
		}

		// Send remaining messages
		if len(messages) > 0 {
			batchResponses, err := f.sendBatch(ctx, messages)
			if err != nil {
				util.Log(ctx).WithError(err).Error("failed to send FCM batch")
				return nil, err
			}
			responses = append(responses, batchResponses...)
		}
	}

	return responses, nil
}

func (f *fcmNotifier) sendBatch(ctx context.Context, messages []*messaging.Message) ([]*devicev1.NotifyResult, error) {
	if len(messages) == 0 {
		return []*devicev1.NotifyResult{}, nil
	}

	br, err := f.client.SendEach(ctx, messages)
	if err != nil {
		return nil, err
	}

	util.Log(ctx).
		WithField("batch_size", len(messages)).
		WithField("success_count", br.SuccessCount).
		WithField("failure_count", br.FailureCount).
		Debug("FCM notification batch sent")

	return f.toNotifyResult(br), nil
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
	// Pre-allocate with exact capacity for better performance
	response := make([]*devicev1.NotifyResult, 0, len(br.Responses))

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
