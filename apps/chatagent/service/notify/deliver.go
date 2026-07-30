// Package notify delivers ChatAgent assistant replies through the Notification service.
//
// ChatAgent stays channel-agnostic: the Turn engine only produces text. This package
// maps a session ChannelBinding onto notificationv1.Send so the same conversation
// can run over SMS, email, push, in-app, WhatsApp, USSD, or web.
package notify

import (
	"context"
	"errors"
	"fmt"
	"strings"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/antinvestor/notification/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/util"
	"google.golang.org/protobuf/types/known/structpb"

	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"
)

// Channel type strings accepted by Notification.Notification.type.
const (
	TypeSMS      = "sms"
	TypeEmail    = "email"
	TypePush     = "push"
	TypeInApp    = "in-app"
	TypeWhatsApp = "whatsapp"
	TypeUSSD     = "ussd"
)

// Binding is the internal delivery target derived from proto ChannelBinding.
type Binding struct {
	Channel         chatagentv1.Channel
	ContactID       string
	ProfileID       string
	ProfileType     string
	Language        string
	SkipDelivery    bool
	Template        string
	SourceContactID string
	SourceProfileID string
	TemplatePayload map[string]any
	RouteID         string
}

// FromProto converts API channel binding to internal form.
func FromProto(p *chatagentv1.ChannelBinding) Binding {
	if p == nil {
		return Binding{Channel: chatagentv1.Channel_CHANNEL_UNSPECIFIED}
	}
	var payload map[string]any
	if tp := p.GetTemplatePayload(); tp != nil {
		payload = tp.AsMap()
	}
	return Binding{
		Channel:         p.GetChannel(),
		ContactID:       strings.TrimSpace(p.GetContactId()),
		ProfileID:       strings.TrimSpace(p.GetProfileId()),
		ProfileType:     strings.TrimSpace(p.GetProfileType()),
		Language:        strings.TrimSpace(p.GetLanguage()),
		SkipDelivery:    p.GetSkipDelivery(),
		Template:        strings.TrimSpace(p.GetTemplate()),
		SourceContactID: strings.TrimSpace(p.GetSourceContactId()),
		SourceProfileID: strings.TrimSpace(p.GetSourceProfileId()),
		TemplatePayload: payload,
		RouteID:         strings.TrimSpace(p.GetRouteId()),
	}
}

// ToProto converts internal binding back to API.
func (b Binding) ToProto() *chatagentv1.ChannelBinding {
	if b.Channel == chatagentv1.Channel_CHANNEL_UNSPECIFIED && b.ContactID == "" {
		return nil
	}
	var payload *structpb.Struct
	if len(b.TemplatePayload) > 0 {
		payload, _ = structpb.NewStruct(b.TemplatePayload)
	}
	profileType := b.ProfileType
	if profileType == "" && b.ProfileID != "" {
		profileType = "Profile"
	}
	return &chatagentv1.ChannelBinding{
		Channel:         b.Channel,
		ContactId:       b.ContactID,
		ProfileId:       b.ProfileID,
		ProfileType:     profileType,
		Language:        b.Language,
		SkipDelivery:    b.SkipDelivery,
		Template:        b.Template,
		SourceContactId: b.SourceContactID,
		SourceProfileId: b.SourceProfileID,
		TemplatePayload: payload,
		RouteId:         b.RouteID,
	}
}

// Name returns a stable short name for persistence/index (sms, email, web, …).
func (b Binding) Name() string {
	return ChannelName(b.Channel)
}

// ChannelName maps enum to storage/lookup key.
func ChannelName(ch chatagentv1.Channel) string {
	switch ch {
	case chatagentv1.Channel_CHANNEL_SMS:
		return TypeSMS
	case chatagentv1.Channel_CHANNEL_EMAIL:
		return TypeEmail
	case chatagentv1.Channel_CHANNEL_PUSH:
		return TypePush
	case chatagentv1.Channel_CHANNEL_IN_APP:
		return TypeInApp
	case chatagentv1.Channel_CHANNEL_WHATSAPP:
		return TypeWhatsApp
	case chatagentv1.Channel_CHANNEL_USSD:
		return TypeUSSD
	case chatagentv1.Channel_CHANNEL_WEB:
		return "web"
	default:
		return "web"
	}
}

// NotificationType is the string Notification service expects for routing.
func (b Binding) NotificationType() string {
	name := b.Name()
	if name == "web" {
		return ""
	}
	return name
}

// ShouldDeliver reports whether a reply should go through Notification.
func (b Binding) ShouldDeliver() bool {
	if b.SkipDelivery {
		return false
	}
	if b.NotificationType() == "" {
		return false
	}
	return b.ContactID != "" || b.ProfileID != ""
}

// Deliverer sends assistant replies to users via Notification.
type Deliverer interface {
	// Deliver queues the reply. Returns delivered=false when delivery is skipped (web / no client).
	Deliver(ctx context.Context, binding Binding, subjectID, sessionID, reply string) (delivered bool, err error)
}

// NotificationDeliverer implements Deliverer with notificationv1.NotificationServiceClient.
type NotificationDeliverer struct {
	cli notificationv1connect.NotificationServiceClient
}

// NewNotificationDeliverer wraps a Notification service client. cli may be nil (no-op deliverer).
func NewNotificationDeliverer(cli notificationv1connect.NotificationServiceClient) *NotificationDeliverer {
	return &NotificationDeliverer{cli: cli}
}

// Deliver implements Deliverer.
func (d *NotificationDeliverer) Deliver(
	ctx context.Context,
	binding Binding,
	subjectID, sessionID, reply string,
) (bool, error) {
	reply = strings.TrimSpace(reply)
	if reply == "" || !binding.ShouldDeliver() {
		return false, nil
	}
	if d == nil || d.cli == nil {
		util.Log(ctx).WithField("session_id", sessionID).
			Debug("chatagent: notification client not configured; skipping channel delivery")
		return false, nil
	}

	nType := binding.NotificationType()
	profileID := binding.ProfileID
	if profileID == "" {
		profileID = strings.TrimSpace(subjectID)
	}
	profileType := binding.ProfileType
	if profileType == "" {
		profileType = "Profile"
	}

	vars := map[string]any{
		"reply":      reply,
		"session_id": sessionID,
		"subject_id": subjectID,
		"channel":    nType,
	}
	for k, v := range binding.TemplatePayload {
		if _, exists := vars[k]; !exists {
			vars[k] = v
		}
	}
	payload, err := structpb.NewStruct(vars)
	if err != nil {
		return false, fmt.Errorf("notification payload: %w", err)
	}

	n := &notificationv1.Notification{
		Recipient: &commonv1.ContactLink{
			ProfileType: profileType,
			ProfileId:   profileID,
			ContactId:   binding.ContactID,
		},
		Type:        nType,
		Template:    binding.Template,
		Payload:     payload,
		Language:    binding.Language,
		OutBound:    true,
		AutoRelease: true,
		RouteId:     binding.RouteID,
		Priority:    notificationv1.PRIORITY_HIGH,
	}
	// Raw body when no template — channel integrations use Notification.data.
	if binding.Template == "" {
		n.Data = reply
	}
	if binding.SourceContactID != "" || binding.SourceProfileID != "" {
		n.Source = &commonv1.ContactLink{
			ContactId: binding.SourceContactID,
			ProfileId: binding.SourceProfileID,
			ProfileType: func() string {
				if binding.SourceProfileID != "" {
					return "Service"
				}
				return ""
			}(),
		}
	}

	stream, err := d.cli.Send(ctx, connect.NewRequest(&notificationv1.SendRequest{
		Data: []*notificationv1.Notification{n},
	}))
	if err != nil {
		return false, fmt.Errorf("notification send: %w", err)
	}
	if stream == nil {
		// Test stubs often return nil stream.
		return true, nil
	}
	for stream.Receive() {
		if rerr := stream.Err(); rerr != nil {
			return false, fmt.Errorf("notification send stream: %w", rerr)
		}
	}
	if rerr := stream.Err(); rerr != nil {
		return false, fmt.Errorf("notification send stream: %w", rerr)
	}
	util.Log(ctx).WithFields(map[string]any{
		"session_id": sessionID,
		"channel":    nType,
		"contact_id": binding.ContactID,
	}).Debug("chatagent: assistant reply delivered via notification")
	return true, nil
}

// NoopDeliverer never delivers (tests / web-only deploys).
type NoopDeliverer struct{}

func (NoopDeliverer) Deliver(context.Context, Binding, string, string, string) (bool, error) {
	return false, nil
}

// ErrNoDeliverer is returned when delivery is required but no client is wired.
var ErrNoDeliverer = errors.New("notification deliverer not configured")
