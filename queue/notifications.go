package queue

import (
	"context"
	napi "github.com/antinvestor/service-notification-api"
)

func Notification(ctx context.Context, ncli *napi.NotificationClient,
	profileId string, contactId string, language string,
	template string, variables map[string]string) error {

	_, err := notificationSendOut(ctx, ncli, profileId, contactId, language, template, variables)
	if err != nil {
		return err
	}

	return nil
}

func notificationSendOut(ctx context.Context, ncli *napi.NotificationClient,
	profileId string, contactId string, language string,
	template string, variables map[string]string) (*napi.StatusResponse, error) {

	//notificationCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	//defer cancel()
	//
	//
	//
	//messageOut := napi.MessageOut{
	//	Autosend:         true,
	//	MessageTemplete:  template,
	//	Language:         language,
	//	ProfileID:        profileId,
	//	ContactID:        contactId,
	//	MessageVariables: variables,
	//}

	//return ncli.Out(notificationCtx, &messageOut)
	return nil, nil
}
