package events

import (
	"context"
	"errors"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/proto"
)

type RelationshipDisConnectQueue struct {
	Service *frame.Service
}

func (rdq *RelationshipDisConnectQueue) Name() string {
	return "relationship.disconnect.queue"
}

func (rdq *RelationshipDisConnectQueue) PayloadType() any {
	pType := ""
	return &pType
}

func (rdq *RelationshipDisConnectQueue) Validate(_ context.Context, payload any) error {
	if _, ok := payload.(*string); !ok {
		return errors.New(" payload is not of type string")
	}

	return nil
}

func (rdq *RelationshipDisConnectQueue) Execute(ctx context.Context, payload any) error {
	relationshipID := *payload.(*string)

	logger := rdq.Service.L().WithField("payload", relationshipID).WithField("type", rdq.Name())
	logger.Debug("handling relationship disconnect")

	relationshipRepo := repository.NewRelationshipRepository(rdq.Service)
	relationship, err := relationshipRepo.GetByID(ctx, relationshipID)
	if err != nil {
		if frame.DBErrorIsRecordNotFound(err) {
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

	profileConfig, _ := rdq.Service.Config().(*config.ProfileConfig)

	// Queue relationship for further processing by peripheral services
	err = rdq.Service.Publish(ctx, profileConfig.QueueRelationshipDisConnectName, binaryProto)
	if err != nil {
		logger.WithError(err).Error("could not publish relationship")
		return err
	}

	logger.WithField("relationship_id", relationship.GetID()).
		Debug(" We have successfully queued relationship disconnect")

	return nil
}
