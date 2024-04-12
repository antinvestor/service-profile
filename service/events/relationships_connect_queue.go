package events

import (
	"context"
	"errors"
	"github.com/antinvestor/service-profile/config"
	"github.com/antinvestor/service-profile/service/repository"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/proto"
)

type RelationshipConnectQueue struct {
	Service *frame.Service
}

func (rcq *RelationshipConnectQueue) Name() string {
	return "relationship.connect.queue"
}

func (rcq *RelationshipConnectQueue) PayloadType() any {
	pType := ""
	return &pType
}

func (rcq *RelationshipConnectQueue) Validate(_ context.Context, payload any) error {
	if _, ok := payload.(*string); !ok {
		return errors.New(" payload is not of type string")
	}

	return nil
}

func (rcq *RelationshipConnectQueue) Execute(ctx context.Context, payload any) error {
	relationshipID := *payload.(*string)

	logger := rcq.Service.L().WithField("payload", relationshipID).WithField("type", rcq.Name())
	logger.Debug("handling rcq")

	relationshipRepo := repository.NewRelationshipRepository(rcq.Service)
	relationship, err := relationshipRepo.GetByID(ctx, relationshipID)
	if err != nil {
		return err
	}

	binaryProto, err := proto.Marshal(relationship.ToAPI())
	if err != nil {
		return err
	}

	profileConfig, _ := rcq.Service.Config().(*config.ProfileConfig)

	// Queue relationship for further processing by peripheral services
	err = rcq.Service.Publish(ctx, profileConfig.QueueRelationshipConnectName, binaryProto)
	if err != nil {
		return err
	}

	logger.WithField("relationship_id", relationship.GetID()).
		Debug(" We have successfully queued relationship connect")

	return nil
}
