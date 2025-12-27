package events_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/antinvestor/service-profile/apps/default/service/events"
	"github.com/antinvestor/service-profile/apps/default/service/models"
)

func TestContactKeyRotationQueue_Name(t *testing.T) {
	queue := events.NewContactKeyRotationQueue(&config.ProfileConfig{}, &config.DEK{}, nil)
	require.Equal(t, events.ContactKeyRotationEventHandlerName, queue.Name())
}

func TestContactKeyRotationQueue_PayloadType(t *testing.T) {
	queue := events.NewContactKeyRotationQueue(&config.ProfileConfig{}, &config.DEK{}, nil)
	payload := queue.PayloadType()
	require.NotNil(t, payload)
	_, ok := payload.(*models.Contact)
	require.True(t, ok)
}

func TestContactKeyRotationQueue_Validate(t *testing.T) {
	queue := events.NewContactKeyRotationQueue(&config.ProfileConfig{}, &config.DEK{}, nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		payload any
		wantErr bool
	}{
		{
			name:    "valid string pointer",
			payload: func() *string { s := "test"; return &s }(),
			wantErr: false,
		},
		{
			name:    "invalid payload type - string",
			payload: "test",
			wantErr: true,
		},
		{
			name:    "invalid payload type - nil",
			payload: nil,
			wantErr: true,
		},
		{
			name:    "invalid payload type - contact",
			payload: &models.Contact{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := queue.Validate(ctx, tt.payload)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestContactKeyRotationQueue_Execute_InvalidPayload(t *testing.T) {
	queue := events.NewContactKeyRotationQueue(&config.ProfileConfig{}, &config.DEK{}, nil)
	ctx := context.Background()

	// Test with invalid payload type
	err := queue.Execute(ctx, "not a string pointer")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid payload type")
}
