package events

import (
	"context"
	"errors"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/repository"

	"google.golang.org/protobuf/proto"

	"github.com/pitabwire/frame"
)

type ClientConnectedSetupQueue struct {
	Service *frame.Service
}

func (csq *ClientConnectedSetupQueue) Name() string {
	return "client.connected.setup.queue"
}

func (csq *ClientConnectedSetupQueue) PayloadType() any {
	pType := ""
	return &pType
}

func (csq *ClientConnectedSetupQueue) Validate(_ context.Context, payload any) error {
	if _, ok := payload.(*string); !ok {
		return errors.New(" payload is not of type string")
	}

	return nil
}

func (csq *ClientConnectedSetupQueue) Execute(ctx context.Context, payload any) error {
	relationshipID := *payload.(*string)

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
