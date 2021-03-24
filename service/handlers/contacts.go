package handlers

import (
	"context"
	napi "github.com/antinvestor/service-notification-api"
	papi "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/models"
	"github.com/go-errors/errors"
	"github.com/pitabwire/frame"
)

func (ps *ProfileServer) GetByContact(ctx context.Context,
	request *papi.ProfileContactRequest, ) (*papi.ProfileObject, error) {

	contact := models.Contact{Detail: request.GetContact()}
	err := contact.GetByDetail(ps.Service.DB(ctx, true))
	if err != nil {
		return nil, err
	}
	return ps.getProfileByID(ctx, contact.Profile.ID)
}

func (ps *ProfileServer) AddContact(ctx context.Context, request *papi.ProfileAddContactRequest,
) (*papi.ProfileObject, error) {

	p := models.Profile{}
	p.ID = request.GetID()
	if err := ps.Service.DB(ctx, true).Find(&p).Error; err != nil {
		return nil, err
	}

	_, err := createContact(ctx, ps.Service, ps.NotificationCli, p.ID, request.GetContact())
	if err != nil {
		return nil, err
	}

	return p.ToObject(ps.Service.DB(ctx, true))
}

func createContact(ctx context.Context, service *frame.Service, ncli *napi.NotificationClient, profileID string, contactDetail string) (*models.Contact, error) {

	contact := models.Contact{Detail: contactDetail}
	if err := contact.Create(service.DB(ctx, false), profileID, contactDetail); err != nil {
		return nil, errors.Wrap(err, 1)
	}
	err := verifyContact(ctx, service, ncli, contact)
	return &contact, err
}

func GetAuthSourceProductID(ctx context.Context) string {
	contextProductId := ctx.Value(config.ContextKeyProductID)
	if contextProductId == nil {
		return ""
	} else {
		return contextProductId.(string)
	}
}

func verifyContact(ctx context.Context, service *frame.Service, ncli *napi.NotificationClient, contact models.Contact) error {
	verification := models.Verification{}

	var productID = GetAuthSourceProductID(ctx)
	err := verification.Create(service.DB(ctx, false), productID, contact, 24*60*60)
	if err != nil {
		return err
	}

	variables := make(map[string]string)
	variables["pin"] = verification.Pin
	variables["linkHash"] = verification.LinkHash
	variables["expiryDate"] = verification.ExpiresAt.String()

	_, err = ncli.Send(ctx, contact.Profile.ID, contact.ID, contact.Language, config.MessageTemplateContactVerification, variables)
		return err


}
