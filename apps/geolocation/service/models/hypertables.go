package models

import (
	"time"

	"github.com/antinvestor/common/timescale"
)

const (
	locationPointsChunkInterval = 24 * time.Hour
	locationPointsCompressAfter = 7 * 24 * time.Hour
	locationPointsRetain        = 90 * 24 * time.Hour
	geoEventsChunkInterval      = 24 * time.Hour
	geoEventsCompressAfter      = 14 * 24 * time.Hour
	geoEventsRetain             = 365 * 24 * time.Hour
)

// Hypertables returns the TimescaleDB configuration for this app's
// append-only tables. Applied idempotently by timescale.Ensure at
// service startup.
func Hypertables() []timescale.Hypertable {
	return []timescale.Hypertable{
		{
			Table:         "location_points",
			TimeColumn:    "true_created_at",
			ChunkInterval: locationPointsChunkInterval,
			SegmentBy:     []string{"partition_id", "subject_id"},
			CompressAfter: locationPointsCompressAfter,
			RetainFor:     locationPointsRetain,
		},
		{
			Table:         "geo_events",
			TimeColumn:    "true_created_at",
			ChunkInterval: geoEventsChunkInterval,
			SegmentBy:     []string{"partition_id", "subject_id"},
			CompressAfter: geoEventsCompressAfter,
			RetainFor:     geoEventsRetain,
		},
	}
}
