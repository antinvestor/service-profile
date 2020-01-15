package handlers

import (
	"antinvestor.com/service/profile/grpc/profile"
	"antinvestor.com/service/profile/models"
	"antinvestor.com/service/profile/queue"
	"antinvestor.com/service/profile/utils"
	"context"
)

func (server *ProfileServer) GetByContact(ctx context.Context,
	request *profile.ProfileContactRequest, ) (*profile.ProfileObject, error) {

	contact := models.Contact{Detail: request.GetContact()}
	err := contact.GetByDetail(server.Env.GetRDb(ctx))
	if err != nil {
		return nil, err
	}
	return server.getProfileByID(ctx, contact.ProfileID)
}

func (server *ProfileServer) AddContact(ctx context.Context, request *profile.ProfileAddContactRequest,
) (*profile.ProfileObject, error) {

	p := models.Profile{}
	p.ProfileID = request.GetID()
	if err := server.Env.GetRDb(ctx).Find(&p).Error; err != nil {
		return nil, err
	}

	err := createContact(server.Env, ctx, p.ProfileID, request.GetContact());
	if err != nil {
		return nil, err
	}

	return p.ToObject(server.Env.GetRDb(ctx))
}

func createContact(env *utils.Env, ctx context.Context, profileID string, contactDetail string) error {

	contact := models.Contact{}
	if err := contact.Create(env.GeWtDb(ctx), profileID, contactDetail); err != nil {
		return err
	}

	return verifyContact(env, ctx, contact)

}

func verifyContact(env *utils.Env, ctx context.Context, contact models.Contact) error{
	verification := models.Verification{}

	var productID = utils.GetAuthSourceProductID(ctx)
	err := verification.Create( env.GeWtDb(ctx), productID,contact, 24*60*60)
	if err != nil {
		return err
	}

	variables := make(map[string]string)
	variables["pin"] = verification.Pin
	variables["linkHash"] = verification.LinkHash
	variables["expiryDate"] = verification.ExpiresAt.String()


	return queue.Notification(env, ctx, contact.ProfileID, contact.ContactID,
		"", utils.MessageTemplateContactVerification, variables)

}
