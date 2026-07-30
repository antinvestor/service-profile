package business //nolint:testpackage // tests unexported NotificationTarget helpers

import (
	"context"
	"errors"
	"testing"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/antinvestor/notification/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/data"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"
)

func TestShouldSendNotification(t *testing.T) {
	t.Parallel()
	require.False(t, shouldSendNotification(nil))
	require.False(t, shouldSendNotification(&chatagentv1.NotificationTarget{Type: "sms"}))
	require.False(t, shouldSendNotification(&chatagentv1.NotificationTarget{
		Type:      "web",
		Recipient: &commonv1.ContactLink{ContactId: "c"},
	}))
	require.False(t, shouldSendNotification(&chatagentv1.NotificationTarget{
		Type:      "sms",
		Recipient: &commonv1.ContactLink{ContactId: "c"},
		Skip:      true,
	}))
	require.True(t, shouldSendNotification(&chatagentv1.NotificationTarget{
		Type:      "sms",
		Recipient: &commonv1.ContactLink{ContactId: "c"},
	}))
	require.True(t, shouldSendNotification(&chatagentv1.NotificationTarget{
		Type:      "email",
		Recipient: &commonv1.ContactLink{ProfileId: "p1"},
	}))
}

func TestBuildOutboundNotification_MatchesNotificationServiceShape(t *testing.T) {
	t.Parallel()
	payload, err := structpb.NewStruct(map[string]any{"brand": "stawi"})
	require.NoError(t, err)
	n := buildOutboundNotification(&chatagentv1.NotificationTarget{
		Type: "sms",
		Recipient: &commonv1.ContactLink{
			ContactId:   "contact-1",
			ProfileId:   "profile-1",
			ProfileType: "Profile",
		},
		Language: "en",
		Payload:  payload,
	}, "profile-1", "sess-9", "What role are you targeting?")

	// Same fields contact_verification_queue / Notification.Send expect.
	require.Equal(t, "sms", n.GetType())
	require.True(t, n.GetOutBound())
	require.True(t, n.GetAutoRelease())
	require.Equal(t, "What role are you targeting?", n.GetData())
	require.Equal(t, "contact-1", n.GetRecipient().GetContactId())
	require.Equal(t, "profile-1", n.GetRecipient().GetProfileId())
	require.Equal(t, "sess-9", n.GetPayload().AsMap()["session_id"])
	require.Equal(t, "What role are you targeting?", n.GetPayload().AsMap()["reply"])
	require.Equal(t, "stawi", n.GetPayload().AsMap()["brand"])
}

func TestNotificationTargetJSONRoundTrip(t *testing.T) {
	t.Parallel()
	in := &chatagentv1.NotificationTarget{
		Type: "whatsapp",
		Recipient: &commonv1.ContactLink{
			ContactId: "ct-1",
			ProfileId: "pf-1",
		},
		Language: "sw",
		Template: "chat.reply",
		RouteId:  "r1",
	}
	m := notificationTargetToJSONMap(in)
	require.NotEmpty(t, m)
	out := notificationTargetFromJSONMap(m)
	require.Equal(t, "whatsapp", out.GetType())
	require.Equal(t, "ct-1", out.GetRecipient().GetContactId())
	require.Equal(t, "pf-1", out.GetRecipient().GetProfileId())
	require.Equal(t, "sw", out.GetLanguage())
	require.Equal(t, "chat.reply", out.GetTemplate())
}

func TestNotificationTargetFromLegacyChannelBindingJSON(t *testing.T) {
	t.Parallel()
	// Older sessions stored ChannelBinding flat fields.
	m := data.JSONMap{
		"channel_name": "sms",
		"contact_id":   "legacy-c",
		"profile_id":   "legacy-p",
		"language":     "en",
	}
	out := notificationTargetFromJSONMap(m)
	require.Equal(t, "sms", out.GetType())
	require.Equal(t, "legacy-c", out.GetRecipient().GetContactId())
	require.Equal(t, "legacy-p", out.GetRecipient().GetProfileId())
}

func TestMergeNotificationTargets(t *testing.T) {
	t.Parallel()
	stored := &chatagentv1.NotificationTarget{
		Type: "sms",
		Recipient: &commonv1.ContactLink{
			ContactId: "old",
			ProfileId: "p1",
		},
		Language: "en",
	}
	inbound := &chatagentv1.NotificationTarget{
		Type: "sms",
		Recipient: &commonv1.ContactLink{
			ContactId: "new",
		},
		Language: "sw",
		Skip:     true,
	}
	out := mergeNotificationTargets(stored, inbound)
	require.Equal(t, "new", out.GetRecipient().GetContactId())
	require.Equal(t, "p1", out.GetRecipient().GetProfileId())
	require.Equal(t, "sw", out.GetLanguage())
	require.True(t, out.GetSkip())
}

// capturingNotificationClient records Send requests — same stub style as profile tests.
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
	return nil, nil //nolint:nilnil // test stub: deliverReply accepts nil stream
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

func TestDeliverReply_UsesNotificationClientSend(t *testing.T) {
	t.Parallel()
	stub := &capturingNotificationClient{}
	b := &chatAgentBusiness{notificationCli: stub}
	ok := b.deliverReply(context.Background(), &chatagentv1.NotificationTarget{
		Type: "sms",
		Recipient: &commonv1.ContactLink{
			ContactId: "c1",
			ProfileId: "p1",
		},
	}, "p1", "sess-1", "Hello from agent")
	require.True(t, ok)
	require.NotNil(t, stub.last)
	require.Len(t, stub.last.GetData(), 1)
	require.Equal(t, "sms", stub.last.GetData()[0].GetType())
	require.Equal(t, "Hello from agent", stub.last.GetData()[0].GetData())
}

func TestDeliverReply_SkipsWhenNoClient(t *testing.T) {
	t.Parallel()
	b := &chatAgentBusiness{}
	ok := b.deliverReply(context.Background(), &chatagentv1.NotificationTarget{
		Type:      "sms",
		Recipient: &commonv1.ContactLink{ContactId: "c1"},
	}, "p1", "sess-1", "hi")
	require.False(t, ok)
}
