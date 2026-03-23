package business

import (
	"context"
	"fmt"
	"time"

	"github.com/pitabwire/frame/events"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
	"github.com/antinvestor/service-profile/apps/geolocation/service/repository"
)

type CatchupConfig struct {
	BatchSize int
	Interval  time.Duration
}

const (
	defaultCatchupBatchSize = 500
	defaultCatchupInterval  = 5 * time.Minute
)

type CatchupBusiness interface {
	RunCatchup(ctx context.Context) error
	StartScheduler(ctx context.Context)
}

type catchupBusiness struct {
	pointRepo repository.LocationPointRepository
	eventsMan events.Manager
	cfg       CatchupConfig
}

func NewCatchupBusiness(
	pointRepo repository.LocationPointRepository,
	eventsMan events.Manager,
	cfg CatchupConfig,
) CatchupBusiness {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = defaultCatchupBatchSize
	}
	if cfg.Interval <= 0 {
		cfg.Interval = defaultCatchupInterval
	}

	return &catchupBusiness{
		pointRepo: pointRepo,
		eventsMan: eventsMan,
		cfg:       cfg,
	}
}

func (b *catchupBusiness) StartScheduler(ctx context.Context) {
	log := util.Log(ctx)

	if err := b.RunCatchup(ctx); err != nil {
		log.WithError(err).Error("initial catchup run failed")
	}

	ticker := time.NewTicker(b.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("catchup scheduler stopped")
			return
		case <-ticker.C:
			if err := b.RunCatchup(ctx); err != nil {
				log.WithError(err).Error("scheduled catchup run failed")
			}
		}
	}
}

func (b *catchupBusiness) RunCatchup(ctx context.Context) error {
	log := util.Log(ctx)

	points, err := b.pointRepo.GetPendingForProcessing(ctx, b.cfg.BatchSize)
	if err != nil {
		return fmt.Errorf("query pending location points: %w", err)
	}
	if len(points) == 0 {
		return nil
	}

	log.Info("catchup: found pending location points", "count", len(points))

	var emitted int
	for _, point := range points {
		event := &models.LocationPointIngestedEvent{
			EventTenancy: models.EventTenancy{
				TenantID:    point.TenantID,
				PartitionID: point.PartitionID,
				AccessID:    point.AccessID,
			},
			PointID:   point.GetID(),
			SubjectID: point.SubjectID,
			DeviceID:  point.DeviceID,
			Latitude:  point.Latitude,
			Longitude: point.Longitude,
			Accuracy:  point.Accuracy,
			Timestamp: point.TS.UnixMilli(),
		}

		if emitErr := b.eventsMan.Emit(ctx, LocationPointIngestedEventName, event); emitErr != nil {
			log.WithError(emitErr).Error("catchup: failed to re-emit event",
				"point_id", point.GetID(),
				"subject_id", point.SubjectID,
			)
			continue
		}
		emitted++
	}

	log.Info("catchup complete", "found", len(points), "emitted", emitted)
	return nil
}
