package queue

import (
	"context"
	"time"

	"github.com/antinvestor/service-profile/grpc/notification"
	"github.com/antinvestor/service-profile/utils"
)

func Notification(env *utils.Env, ctx context.Context,
	profileId string, contactId string, language string,
	template string, variables map[string]string) error {

	_, err := notificationSendOut(env, ctx, profileId, contactId, language, template, variables)
	if err != nil{
		return err
	}

	return nil
}

func notificationSendOut(env *utils.Env, ctx context.Context,
	profileId string, contactId string, language string,
	template string, variables map[string]string) (*notification.StatusResponse, error) {

	notificationCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	notificationService := notification.NewNotificationServiceClient(env.GetNotificationServiceConn())

	messageOut := notification.MessageOut{
		Autosend:         "true",
		MessageTemplete:  template,
		Language:         language,
		ProfileID:        profileId,
		ContactID:        contactId,
		MessageVariables: variables,
	}

	return notificationService.Out(notificationCtx, &messageOut)
}
