package notify_test

import (
	"context"
	"errors"
	"testing"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/antinvestor/notification/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/antinvestor/service-profile/apps/chatagent/service/notify"
	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"
)

func TestBinding_ShouldDeliver(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		b    notify.Binding
		want bool
	}{
		{"web", notify.Binding{Channel: chatagentv1.Channel_CHANNEL_WEB, ContactID: "c1"}, false},
		{"sms no contact", notify.Binding{Channel: chatagentv1.Channel_CHANNEL_SMS}, false},
		{"sms contact", notify.Binding{Channel: chatagentv1.Channel_CHANNEL_SMS, ContactID: "c1"}, true},
		{"whatsapp profile only", notify.Binding{Channel: chatagentv1.Channel_CHANNEL_WHATSAPP, ProfileID: "p1"}, true},
		{"sms skip", notify.Binding{Channel: chatagentv1.Channel_CHANNEL_SMS, ContactID: "c1", SkipDelivery: true}, false},
		{"email", notify.Binding{Channel: chatagentv1.Channel_CHANNEL_EMAIL, ContactID: "c1"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, tc.b.ShouldDeliver())
		})
	}
}

func TestChannelName(t *testing.T) {
	t.Parallel()
	require.Equal(t, "sms", notify.ChannelName(chatagentv1.Channel_CHANNEL_SMS))
	require.Equal(t, "whatsapp", notify.ChannelName(chatagentv1.Channel_CHANNEL_WHATSAPP))
	require.Equal(t, "web", notify.ChannelName(chatagentv1.Channel_CHANNEL_WEB))
	require.Equal(t, "in-app", notify.ChannelName(chatagentv1.Channel_CHANNEL_IN_APP))
}

func TestFromProtoRoundTrip(t *testing.T) {
	t.Parallel()
	payload, err := structpb.NewStruct(map[string]any{"brand": "stawi"})
	require.NoError(t, err)
	p := &chatagentv1.ChannelBinding{
		Channel:         chatagentv1.Channel_CHANNEL_WHATSAPP,
		ContactId:       "contact-1",
		ProfileId:       "profile-1",
		ProfileType:     "Profile",
		Language:        "en",
		Template:        "chat.agent.reply",
		SourceContactId: "bot",
		TemplatePayload: payload,
		RouteId:         "route-wa",
	}
	b := notify.FromProto(p)
	require.True(t, b.ShouldDeliver())
	out := b.ToProto()
	require.Equal(t, p.GetChannel(), out.GetChannel())
	require.Equal(t, p.GetContactId(), out.GetContactId())
	require.Equal(t, p.GetTemplate(), out.GetTemplate())
	require.Equal(t, "stawi", out.GetTemplatePayload().AsMap()["brand"])
}

// capturingNotificationClient records Send requests for tests.
type capturingNotificationClient struct {
	last *notificationv1.SendRequest
	err  error
}

var _ notificationv1connect.NotificationServiceClient = (*capturingNotificationClient)(nil)

func (c *capturingNotificationClient) Send(
	_ context.Context,
	req *connect.Request[notificationv1.SendRequest],
) (*connect.ServerStreamForClient[notificationv1.SendResponse], error) {
	if c.err != nil {
		return nil, c.err
	}
	c.last = req.Msg
	// nil stream is treated as success by deliverer (test mode).
	return nil, nil //nolint:nilnil // test stub: deliverer accepts nil stream
}

func (c *capturingNotificationClient) Release(context.Context, *connect.Request[notificationv1.ReleaseRequest]) (*connect.ServerStreamForClient[notificationv1.ReleaseResponse], error) {
	return nil, errUnused
}
func (c *capturingNotificationClient) Receive(context.Context, *connect.Request[notificationv1.ReceiveRequest]) (*connect.ServerStreamForClient[notificationv1.ReceiveResponse], error) {
	return nil, errUnused
}
func (c *capturingNotificationClient) Search(context.Context, *connect.Request[commonv1.SearchRequest]) (*connect.ServerStreamForClient[notificationv1.SearchResponse], error) {
	return nil, errUnused
}
func (c *capturingNotificationClient) Status(context.Context, *connect.Request[commonv1.StatusRequest]) (*connect.Response[commonv1.StatusResponse], error) {
	return nil, errUnused
}
func (c *capturingNotificationClient) StatusUpdate(context.Context, *connect.Request[commonv1.StatusUpdateRequest]) (*connect.Response[commonv1.StatusUpdateResponse], error) {
	return nil, errUnused
}
func (c *capturingNotificationClient) TemplateSearch(context.Context, *connect.Request[notificationv1.TemplateSearchRequest]) (*connect.ServerStreamForClient[notificationv1.TemplateSearchResponse], error) {
	return nil, errUnused
}
func (c *capturingNotificationClient) TemplateSave(context.Context, *connect.Request[notificationv1.TemplateSaveRequest]) (*connect.Response[notificationv1.TemplateSaveResponse], error) {
	return nil, errUnused
}

var errUnused = errors.New("unused in this test")

func TestNotificationDeliverer_SendsRawBody(t *testing.T) {
	t.Parallel()
	stub := &capturingNotificationClient{}
	d := notify.NewNotificationDeliverer(stub)
	delivered, err := d.Deliver(context.Background(), notify.Binding{
		Channel:   chatagentv1.Channel_CHANNEL_SMS,
		ContactID: "contact-sms",
		ProfileID: "prof-1",
		Language:  "en",
	}, "prof-1", "sess-9", "What role are you targeting?")
	require.NoError(t, err)
	require.True(t, delivered)
	require.NotNil(t, stub.last)
	data := stub.last.GetData()
	require.Len(t, data, 1)
	n := data[0]
	require.Equal(t, "sms", n.GetType())
	require.Equal(t, "What role are you targeting?", n.GetData())
	require.True(t, n.GetOutBound())
	require.True(t, n.GetAutoRelease())
	require.Equal(t, "contact-sms", n.GetRecipient().GetContactId())
	require.Equal(t, "sess-9", n.GetPayload().AsMap()["session_id"])
	require.Equal(t, "What role are you targeting?", n.GetPayload().AsMap()["reply"])
}

func TestNotificationDeliverer_SkipsWeb(t *testing.T) {
	t.Parallel()
	stub := &capturingNotificationClient{}
	d := notify.NewNotificationDeliverer(stub)
	delivered, err := d.Deliver(context.Background(), notify.Binding{
		Channel:   chatagentv1.Channel_CHANNEL_WEB,
		ContactID: "c",
	}, "s", "sess", "hi")
	require.NoError(t, err)
	require.False(t, delivered)
	require.Nil(t, stub.last)
}

func TestNoopDeliverer(t *testing.T) {
	t.Parallel()
	delivered, err := notify.NoopDeliverer{}.Deliver(context.Background(), notify.Binding{
		Channel: chatagentv1.Channel_CHANNEL_SMS, ContactID: "c",
	}, "s", "sess", "hi")
	require.NoError(t, err)
	require.False(t, delivered)
}
