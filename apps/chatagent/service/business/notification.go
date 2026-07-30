package business

import (
	"context"
	"strings"

	commonv1 "buf.build/gen/go/antinvestor/common/protocolbuffers/go/common/v1"
	notificationv1 "buf.build/gen/go/antinvestor/notification/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/data"
	"github.com/pitabwire/util"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/antinvestor/service-profile/apps/chatagent/service/models"
	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"
)

// shouldSendNotification reports whether a NotificationTarget should trigger Notification.Send.
// Empty type = RPC-only (web). Matches how Notification service is used elsewhere: type + recipient.
func shouldSendNotification(t *chatagentv1.NotificationTarget) bool {
	if t == nil || t.GetSkip() {
		return false
	}
	typ := strings.TrimSpace(strings.ToLower(t.GetType()))
	if typ == "" || typ == "web" {
		return false
	}
	r := t.GetRecipient()
	if r == nil {
		return false
	}
	return strings.TrimSpace(r.GetContactId()) != "" || strings.TrimSpace(r.GetProfileId()) != ""
}

// deliverReply queues the assistant reply via the existing Notification service client.
// Pattern matches apps/default contact verification: build notificationv1.Notification → Send → drain stream.
// Soft-fails (logs, returns false) so channel issues never break the turn.
func (b *chatAgentBusiness) deliverReply(
	ctx context.Context,
	target *chatagentv1.NotificationTarget,
	subjectID, sessionID, reply string,
) bool {
	reply = strings.TrimSpace(reply)
	if reply == "" || !shouldSendNotification(target) {
		return false
	}
	if b.notificationCli == nil {
		util.Log(ctx).WithField("session_id", sessionID).
			Debug("chatagent: Notification client not configured; skipping Send")
		return false
	}

	n := buildOutboundNotification(target, subjectID, sessionID, reply)
	stream, err := b.notificationCli.Send(ctx, connect.NewRequest(&notificationv1.SendRequest{
		Data: []*notificationv1.Notification{n},
	}))
	if err != nil {
		util.Log(ctx).WithError(err).WithFields(map[string]any{
			"session_id": sessionID,
			"type":       n.GetType(),
		}).Warn("chatagent: Notification.Send failed (turn still succeeded)")
		return false
	}
	if stream == nil {
		// Test stubs often return nil stream (same as profile verification tests).
		return true
	}
	for stream.Receive() {
		if rerr := stream.Err(); rerr != nil {
			util.Log(ctx).WithError(rerr).WithField("session_id", sessionID).
				Warn("chatagent: Notification.Send stream error (turn still succeeded)")
			return false
		}
	}
	if rerr := stream.Err(); rerr != nil {
		util.Log(ctx).WithError(rerr).WithField("session_id", sessionID).
			Warn("chatagent: Notification.Send stream error (turn still succeeded)")
		return false
	}
	util.Log(ctx).WithFields(map[string]any{
		"session_id": sessionID,
		"type":       n.GetType(),
		"contact_id": n.GetRecipient().GetContactId(),
	}).Debug("chatagent: assistant reply queued via Notification.Send")
	return true
}

// buildOutboundNotification maps NotificationTarget → notificationv1.Notification
// using the same fields the Notification service already understands.
func buildOutboundNotification(
	target *chatagentv1.NotificationTarget,
	subjectID, sessionID, reply string,
) *notificationv1.Notification {
	typ := strings.TrimSpace(strings.ToLower(target.GetType()))
	recipient := cloneContactLink(target.GetRecipient())
	if recipient == nil {
		recipient = &commonv1.ContactLink{}
	}
	if strings.TrimSpace(recipient.GetProfileId()) == "" {
		recipient.ProfileId = strings.TrimSpace(subjectID)
	}
	if strings.TrimSpace(recipient.GetProfileType()) == "" && recipient.GetProfileId() != "" {
		recipient.ProfileType = "Profile"
	}

	vars := map[string]any{
		"reply":      reply,
		"session_id": sessionID,
		"subject_id": subjectID,
		"type":       typ,
	}
	if p := target.GetPayload(); p != nil {
		for k, v := range p.AsMap() {
			if _, exists := vars[k]; !exists {
				vars[k] = v
			}
		}
	}
	payload, _ := structpb.NewStruct(vars)

	n := &notificationv1.Notification{
		Recipient:   recipient,
		Source:      cloneContactLink(target.GetSource()),
		Type:        typ,
		Template:    strings.TrimSpace(target.GetTemplate()),
		Payload:     payload,
		Language:    strings.TrimSpace(target.GetLanguage()),
		OutBound:    true,
		AutoRelease: true,
		RouteId:     strings.TrimSpace(target.GetRouteId()),
		Priority:    notificationv1.PRIORITY_HIGH,
	}
	// Same convention as other Send callers: raw body when no template.
	if n.GetTemplate() == "" {
		n.Data = reply
	}
	return n
}

func cloneContactLink(in *commonv1.ContactLink) *commonv1.ContactLink {
	if in == nil {
		return nil
	}
	return &commonv1.ContactLink{
		ProfileName:    in.GetProfileName(),
		ProfileType:    in.GetProfileType(),
		ProfileId:      in.GetProfileId(),
		ProfileImageId: in.GetProfileImageId(),
		ContactId:      in.GetContactId(),
	}
}

// normalizeNotificationTarget fills defaults (profile_id from subject) without inventing channels.
func normalizeNotificationTarget(t *chatagentv1.NotificationTarget, subjectID string) *chatagentv1.NotificationTarget {
	if t == nil {
		return nil
	}
	out := &chatagentv1.NotificationTarget{
		Type:     strings.TrimSpace(t.GetType()),
		Language: strings.TrimSpace(t.GetLanguage()),
		Template: strings.TrimSpace(t.GetTemplate()),
		Payload:  t.GetPayload(),
		RouteId:  strings.TrimSpace(t.GetRouteId()),
		Skip:     t.GetSkip(),
		Source:   cloneContactLink(t.GetSource()),
	}
	r := cloneContactLink(t.GetRecipient())
	if r == nil {
		r = &commonv1.ContactLink{}
	}
	if strings.TrimSpace(r.GetProfileId()) == "" {
		r.ProfileId = strings.TrimSpace(subjectID)
	}
	if strings.TrimSpace(r.GetProfileType()) == "" && r.GetProfileId() != "" {
		r.ProfileType = "Profile"
	}
	// Drop empty recipient entirely when no identifiers.
	if strings.TrimSpace(r.GetContactId()) == "" && strings.TrimSpace(r.GetProfileId()) == "" {
		r = nil
	}
	out.Recipient = r
	if out.GetType() == "" && r == nil && out.GetTemplate() == "" && !out.GetSkip() {
		return nil
	}
	return out
}

// mergeNotificationTargets prefers inbound non-empty fields over stored (for IngestMessage).
func mergeNotificationTargets(stored, inbound *chatagentv1.NotificationTarget) *chatagentv1.NotificationTarget {
	if stored == nil {
		return normalizeNotificationTarget(inbound, "")
	}
	if inbound == nil {
		return stored
	}
	out := &chatagentv1.NotificationTarget{
		Type:      stored.GetType(),
		Language:  stored.GetLanguage(),
		Template:  stored.GetTemplate(),
		Payload:   stored.GetPayload(),
		RouteId:   stored.GetRouteId(),
		Skip:      stored.GetSkip(),
		Recipient: cloneContactLink(stored.GetRecipient()),
		Source:    cloneContactLink(stored.GetSource()),
	}
	if t := strings.TrimSpace(inbound.GetType()); t != "" {
		out.Type = t
	}
	if t := strings.TrimSpace(inbound.GetLanguage()); t != "" {
		out.Language = t
	}
	if t := strings.TrimSpace(inbound.GetTemplate()); t != "" {
		out.Template = t
	}
	if t := strings.TrimSpace(inbound.GetRouteId()); t != "" {
		out.RouteId = t
	}
	if inbound.GetPayload() != nil {
		out.Payload = inbound.GetPayload()
	}
	if inbound.GetSkip() {
		out.Skip = true
	}
	if ir := inbound.GetRecipient(); ir != nil {
		if out.GetRecipient() == nil {
			out.Recipient = &commonv1.ContactLink{}
		}
		if id := strings.TrimSpace(ir.GetContactId()); id != "" {
			out.Recipient.ContactId = id
		}
		if id := strings.TrimSpace(ir.GetProfileId()); id != "" {
			out.Recipient.ProfileId = id
		}
		if pt := strings.TrimSpace(ir.GetProfileType()); pt != "" {
			out.Recipient.ProfileType = pt
		}
		if n := strings.TrimSpace(ir.GetProfileName()); n != "" {
			out.Recipient.ProfileName = n
		}
	}
	if is := inbound.GetSource(); is != nil {
		out.Source = cloneContactLink(is)
	}
	return out
}

// notificationTargetToJSONMap persists NotificationTarget on the session row.
func notificationTargetToJSONMap(t *chatagentv1.NotificationTarget) data.JSONMap {
	if t == nil {
		return data.JSONMap{}
	}
	m := map[string]any{
		"type":     strings.TrimSpace(t.GetType()),
		"language": strings.TrimSpace(t.GetLanguage()),
		"template": strings.TrimSpace(t.GetTemplate()),
		"route_id": strings.TrimSpace(t.GetRouteId()),
		"skip":     t.GetSkip(),
	}
	if r := t.GetRecipient(); r != nil {
		m["recipient"] = map[string]any{
			"profile_name":     r.GetProfileName(),
			"profile_type":     r.GetProfileType(),
			"profile_id":       r.GetProfileId(),
			"profile_image_id": r.GetProfileImageId(),
			"contact_id":       r.GetContactId(),
		}
	}
	if s := t.GetSource(); s != nil {
		m["source"] = map[string]any{
			"profile_name":     s.GetProfileName(),
			"profile_type":     s.GetProfileType(),
			"profile_id":       s.GetProfileId(),
			"profile_image_id": s.GetProfileImageId(),
			"contact_id":       s.GetContactId(),
		}
	}
	if p := t.GetPayload(); p != nil {
		m["payload"] = p.AsMap()
	}
	jm, err := models.JSONMapFromStruct(m)
	if err != nil {
		return data.JSONMap{}
	}
	return jm
}

func notificationTargetFromJSONMap(m data.JSONMap) *chatagentv1.NotificationTarget {
	if len(m) == 0 {
		return nil
	}
	t := &chatagentv1.NotificationTarget{
		Type:     stringField(m, "type"),
		Language: stringField(m, "language"),
		Template: stringField(m, "template"),
		RouteId:  stringField(m, "route_id"),
		Skip:     boolField(m, "skip"),
	}
	if raw, ok := m["recipient"]; ok {
		t.Recipient = contactLinkFromAny(raw)
	}
	// Backward-compat for earlier ChannelBinding JSON shape (contact_id / profile_id at top level).
	if t.GetRecipient() == nil {
		cid := stringField(m, "contact_id")
		pid := stringField(m, "profile_id")
		if cid != "" || pid != "" {
			t.Recipient = &commonv1.ContactLink{
				ContactId:   cid,
				ProfileId:   pid,
				ProfileType: stringField(m, "profile_type"),
			}
		}
	}
	if raw, ok := m["source"]; ok {
		t.Source = contactLinkFromAny(raw)
	}
	if t.GetSource() == nil {
		sc := stringField(m, "source_contact_id")
		sp := stringField(m, "source_profile_id")
		if sc != "" || sp != "" {
			t.Source = &commonv1.ContactLink{ContactId: sc, ProfileId: sp}
		}
	}
	// channel_name / channel enum from earlier shape → type string
	if t.GetType() == "" {
		if name := stringField(m, "channel_name"); name != "" {
			t.Type = name
		}
	}
	if raw, ok := m["payload"]; ok {
		if mp, isMap := raw.(map[string]any); isMap {
			t.Payload, _ = structpb.NewStruct(mp)
		}
	}
	if raw, ok := m["template_payload"]; ok && t.GetPayload() == nil {
		if mp, isMap := raw.(map[string]any); isMap {
			t.Payload, _ = structpb.NewStruct(mp)
		}
	}
	if !t.GetSkip() && boolField(m, "skip_delivery") {
		t.Skip = true
	}
	if t.GetType() == "" && t.GetRecipient() == nil && t.GetTemplate() == "" {
		return nil
	}
	return t
}

func contactLinkFromAny(raw any) *commonv1.ContactLink {
	mp, ok := raw.(map[string]any)
	if !ok || mp == nil {
		return nil
	}
	get := func(k string) string {
		v, exists := mp[k]
		if !exists || v == nil {
			return ""
		}
		s, _ := v.(string)
		return strings.TrimSpace(s)
	}
	c := &commonv1.ContactLink{
		ProfileName:    get("profile_name"),
		ProfileType:    get("profile_type"),
		ProfileId:      get("profile_id"),
		ProfileImageId: get("profile_image_id"),
		ContactId:      get("contact_id"),
	}
	if c.GetContactId() == "" && c.GetProfileId() == "" {
		return nil
	}
	return c
}

func notificationType(t *chatagentv1.NotificationTarget) string {
	if t == nil {
		return ""
	}
	return strings.TrimSpace(strings.ToLower(t.GetType()))
}

func notificationContactID(t *chatagentv1.NotificationTarget) string {
	if t == nil || t.GetRecipient() == nil {
		return ""
	}
	return strings.TrimSpace(t.GetRecipient().GetContactId())
}

func notificationProfileID(t *chatagentv1.NotificationTarget) string {
	if t == nil || t.GetRecipient() == nil {
		return ""
	}
	return strings.TrimSpace(t.GetRecipient().GetProfileId())
}
