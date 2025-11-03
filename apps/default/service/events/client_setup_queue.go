package events

import (
	"context"
	"errors"

	"github.com/pitabwire/frame/data"
	frevents "github.com/pitabwire/frame/events"
	"github.com/pitabwire/frame/queue"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/repository"
)

const ClientConnectedSetupQueueName = "client.connected.setup.queue"

type ClientConnectedSetupQueue struct {
	eventsMan        frevents.Manager
	queueMan         queue.Manager
	relationshipRepo repository.RelationshipRepository

	relationshipTopicName string
}

func NewClientConnectedSetupQueue(
	_ context.Context,
	cfg *config.ProfileConfig,
	queueMan queue.Manager,
	eventsMan frevents.Manager,
	relationshipRepo repository.RelationshipRepository,
) *ClientConnectedSetupQueue {
	return &ClientConnectedSetupQueue{
		eventsMan:             eventsMan,
		queueMan:              queueMan,
		relationshipRepo:      relationshipRepo,
		relationshipTopicName: cfg.QueueRelationshipConnectName,
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

	logger := util.Log(ctx).WithField("payload", relationshipID).WithField("type", csq.Name())
	logger.Debug("handling csq")

	relationship, err := csq.relationshipRepo.GetByID(ctx, relationshipID)
	if err != nil {
		if data.ErrorIsNoRows(err) {
			logger.WithError(err).Error("no such relationship exists")
			return nil
		}
		logger.WithError(err).Error("could not get relationship")
		return err
	}

	relationshipTopic, err := csq.queueMan.GetPublisher(csq.relationshipTopicName)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not get  publisher")
	}

	// Queue relationship for further processing by peripheral services
	err = relationshipTopic.Publish(ctx, relationship.ToAPI())
	if err != nil {
		logger.WithError(err).Error("could not publish relationship")
		return err
	}

	logger.WithField("relationship_id", relationship.GetID()).
		Debug(" We have successfully queued relationship connect")

	return nil
}
