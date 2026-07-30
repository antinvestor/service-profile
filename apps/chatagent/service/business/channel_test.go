package business

import (
	"testing"

	"github.com/pitabwire/frame/v2/data"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/chatagent/service/notify"
	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"
)

func TestChannelBindingJSONRoundTrip(t *testing.T) {
	t.Parallel()
	in := notify.Binding{
		Channel:         chatagentv1.Channel_CHANNEL_SMS,
		ContactID:       "ct-1",
		ProfileID:       "pf-1",
		ProfileType:     "Profile",
		Language:        "sw",
		Template:        "chat.reply",
		SourceContactID: "bot",
		TemplatePayload: map[string]any{"k": "v"},
		RouteID:         "r1",
	}
	m := channelBindingToJSONMap(in)
	require.NotEmpty(t, m)
	out := channelBindingFromJSONMap(m)
	require.Equal(t, in.Channel, out.Channel)
	require.Equal(t, in.ContactID, out.ContactID)
	require.Equal(t, in.ProfileID, out.ProfileID)
	require.Equal(t, in.Language, out.Language)
	require.Equal(t, in.Template, out.Template)
	require.Equal(t, in.RouteID, out.RouteID)
	require.Equal(t, "v", out.TemplatePayload["k"])
	require.Equal(t, "sms", out.Name())
}

func TestChannelBindingFromEmpty(t *testing.T) {
	t.Parallel()
	b := channelBindingFromJSONMap(data.JSONMap{})
	require.False(t, b.ShouldDeliver())
	require.Nil(t, b.ToProto())
}

func TestMergeBindings(t *testing.T) {
	t.Parallel()
	stored := notify.Binding{
		Channel:   chatagentv1.Channel_CHANNEL_SMS,
		ContactID: "old",
		ProfileID: "p1",
		Language:  "en",
	}
	inbound := notify.Binding{
		Channel:      chatagentv1.Channel_CHANNEL_SMS,
		ContactID:    "new",
		Language:     "sw",
		SkipDelivery: true,
	}
	out := mergeBindings(stored, inbound)
	require.Equal(t, "new", out.ContactID)
	require.Equal(t, "p1", out.ProfileID)
	require.Equal(t, "sw", out.Language)
	require.True(t, out.SkipDelivery)
}

func TestChannelFromName(t *testing.T) {
	t.Parallel()
	require.Equal(t, chatagentv1.Channel_CHANNEL_WHATSAPP, channelFromName("whatsapp"))
	require.Equal(t, chatagentv1.Channel_CHANNEL_IN_APP, channelFromName("in-app"))
	require.Equal(t, chatagentv1.Channel_CHANNEL_UNSPECIFIED, channelFromName(""))
}
