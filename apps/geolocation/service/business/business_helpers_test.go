package business //nolint:testpackage // tests access unexported helper functions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/antinvestor/service-profile/apps/geolocation/service/models"
)

func TestBusinessHelpers(t *testing.T) {
	t.Parallel()

	t.Run("compute confidence", func(t *testing.T) {
		t.Parallel()
		require.InDelta(t, 1.0, computeConfidence(0), 0.001)
		require.Greater(t, computeConfidence(5), computeConfidence(50))
		require.InDelta(t, minConfidence, computeConfidence(10_000), 0.001)
	})

	t.Run("apply hysteresis", func(t *testing.T) {
		cases := []struct {
			name           string
			rawInside      bool
			currentInside  bool
			accuracyMeters float64
			bufferMeters   float64
			want           bool
		}{
			{
				name:           "same state",
				rawInside:      true,
				currentInside:  true,
				accuracyMeters: 20,
				bufferMeters:   30,
				want:           true,
			},
			{
				name:           "good accuracy transition",
				rawInside:      true,
				currentInside:  false,
				accuracyMeters: 5,
				bufferMeters:   30,
				want:           true,
			},
			{
				name:           "poor accuracy hold",
				rawInside:      true,
				currentInside:  false,
				accuracyMeters: 50,
				bufferMeters:   30,
				want:           false,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(
					t,
					tc.want,
					applyHysteresis(tc.rawInside, tc.currentInside, tc.accuracyMeters, tc.bufferMeters),
				)
			})
		}
	})

	t.Run("thresholds and conversions", func(t *testing.T) {
		threshold := 15.0
		consecutive := 4
		cooldown := 90
		route := &models.Route{
			DeviationThresholdM:       &threshold,
			DeviationConsecutiveCount: &consecutive,
			DeviationCooldownSec:      &cooldown,
		}
		resolved := resolveThresholds(route)
		require.InDelta(t, threshold, resolved.thresholdM, 0.001)
		require.Equal(t, consecutive, resolved.consecutiveCount)
		require.Equal(t, cooldown, resolved.cooldownSec)

		value := int32(7)
		require.Equal(t, 7, *int32PtrToIntPtr(&value))
		require.Nil(t, int32PtrToIntPtr(nil))
	})

	t.Run("route deviation trigger logic", func(t *testing.T) {
		now := time.Now().UTC()
		biz := &routeDeviationBusiness{}
		thresholds := routeThresholds{consecutiveCount: 2, cooldownSec: 30}

		cases := []struct {
			name     string
			state    *models.RouteDeviationState
			at       time.Time
			expected bool
		}{
			{
				name:     "not enough consecutive points",
				state:    &models.RouteDeviationState{ConsecutiveOffRoute: 1},
				at:       now,
				expected: false,
			},
			{
				name:     "first deviation triggers",
				state:    &models.RouteDeviationState{ConsecutiveOffRoute: 2, Deviated: false},
				at:       now,
				expected: true,
			},
			{
				name:     "already deviated without last event still triggers",
				state:    &models.RouteDeviationState{ConsecutiveOffRoute: 2, Deviated: true},
				at:       now,
				expected: true,
			},
			{
				name: "cooldown not met",
				state: &models.RouteDeviationState{
					ConsecutiveOffRoute: 2,
					Deviated:            true,
					LastDeviationEventAt: func() *time.Time {
						ts := now.Add(-10 * time.Second)
						return &ts
					}(),
				},
				at:       now,
				expected: false,
			},
			{
				name: "cooldown met",
				state: &models.RouteDeviationState{
					ConsecutiveOffRoute: 2,
					Deviated:            true,
					LastDeviationEventAt: func() *time.Time {
						ts := now.Add(-31 * time.Second)
						return &ts
					}(),
				},
				at:       now,
				expected: true,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.expected, biz.shouldTriggerDeviation(tc.state, tc.at, thresholds))
			})
		}
	})

	t.Run("time helpers", func(t *testing.T) {
		now := time.Now().UTC().Truncate(time.Second)
		from, to := resolveTimeRange(timestamppb.New(now.Add(-time.Hour)), timestamppb.New(now))
		require.Equal(t, now.Add(-time.Hour), from)
		require.Equal(t, now, to)

		defaultFrom, defaultTo := resolveTimeRange(nil, nil)
		require.WithinDuration(t, time.Now().UTC(), defaultTo, 2*time.Second)
		require.WithinDuration(t, defaultTo.Add(-24*time.Hour), defaultFrom, 2*time.Second)

		ts := timestampFromTime(now)
		require.True(t, ts.AsTime().Equal(now))
		require.Equal(t, 10, clampLimit(50, 10, 10))
		require.Equal(t, 20, clampLimit(0, 20, 50))
	})
}
