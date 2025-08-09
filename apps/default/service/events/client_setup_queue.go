package events

import (
	"context"
	"errors"

	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/proto"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

const ClientConnectedSetupQueueName = "client.connected.setup.queue"

type ClientConnectedSetupQueue struct {
	Service *frame.Service
}

func NewClientConnectedSetupQueue(service *frame.Service) *ClientConnectedSetupQueue {
	return &ClientConnectedSetupQueue{
		Service: service,
	}
}

func (csq *ClientConnectedSetupQueue) Name() string {
	return ClientConnectedSetupQueueName
}

func (csq *ClientConnectedSetupQueue) PayloadType() any {
	pType := ""
	return &pType
}

func (csq *ClientConnectedSetupQueue) Validate(_ context.Context, payload any) error {
	_, ok := payload.(*string)
	if !ok {
		return errors.New("invalid payload type, expected *string")
	}

	return nil
}

func (csq *ClientConnectedSetupQueue) Execute(ctx context.Context, payload any) error {
	relationshipIDPtr, ok := payload.(*string)
	if !ok {
		return errors.New("invalid payload type, expected *string")
	}
	relationshipID := *relationshipIDPtr

	logger := csq.Service.Log(ctx).WithField("payload", relationshipID).WithField("type", csq.Name())
	logger.Debug("handling csq")

	relationshipRepo := repository.NewRelationshipRepository(csq.Service)
	relationship, err := relationshipRepo.GetByID(ctx, relationshipID)
	if err != nil {
		if frame.ErrorIsNoRows(err) {
			logger.WithError(err).Error("no such relationship exists")
			return nil
		}
		logger.WithError(err).Error("could not get relationship")
		return err
	}

	binaryProto, err := proto.Marshal(relationship.ToAPI())
	if err != nil {
		logger.WithError(err).Error("could not encode api object")
		return err
	}

	profileConfig, _ := csq.Service.Config().(*config.ProfileConfig)

	// Queue relationship for further processing by peripheral services
	err = csq.Service.Publish(ctx, profileConfig.QueueRelationshipConnectName, binaryProto)
	if err != nil {
		logger.WithError(err).Error("could not publish relationship")
		return err
	}

	logger.WithField("relationship_id", relationship.GetID()).
		Debug(" We have successfully queued relationship connect")

	return nil
}
