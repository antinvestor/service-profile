package events

import (
	"context"
	"errors"

	"github.com/pitabwire/frame"

	"github.com/antinvestor/service-profile/apps/settings/service/models"
	"github.com/antinvestor/service-profile/apps/settings/service/repository"
)

type SettingsAuditor struct {
	Service *frame.Service
}

func (e *SettingsAuditor) Name() string {
	return "setting.audit"
}

func (e *SettingsAuditor) PayloadType() interface{} {
	return &models.SettingAudit{}
}

func (e *SettingsAuditor) Validate(_ context.Context, payload interface{}) error {
	audit, ok := payload.(*models.SettingAudit)
	if !ok {
		return errors.New(" payload is not of type models.SettingAudit")
	}

	if audit.GetID() == "" {
		return errors.New(" audit Id should already have been set ")
	}

	return nil
}

func (e *SettingsAuditor) Execute(ctx context.Context, payload interface{}) error {
	audit, ok := payload.(*models.SettingAudit)
	if !ok {
		return errors.New("payload is not of type models.SettingAudit")
	}

	log := e.Service.Log(ctx).WithField("type", e.Name())
	log.WithField("payload", audit).Info("handling event")

	auditRepo := repository.NewSettingAuditRepository(ctx, e.Service)
	err := auditRepo.Save(ctx, audit)
	if err != nil {
		log.WithError(err).Warn("could not save audit to db")
		return err
	}
	return nil
}
