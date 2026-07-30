package business

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/data"
	"github.com/pitabwire/util"
	"gorm.io/gorm"

	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"

	"github.com/antinvestor/service-profile/apps/chatagent/service/engine"
	"github.com/antinvestor/service-profile/apps/chatagent/service/models"
	"github.com/antinvestor/service-profile/apps/chatagent/service/repository"
)

// ChatAgentBusiness is the domain API for intake sessions.
type ChatAgentBusiness interface {
	UpsertContext(ctx context.Context, req *chatagentv1.UpsertContextRequest) (*chatagentv1.UpsertContextResponse, error)
	GetContext(ctx context.Context, req *chatagentv1.GetContextRequest) (*chatagentv1.GetContextResponse, error)
	ListContexts(ctx context.Context, req *chatagentv1.ListContextsRequest) (*chatagentv1.ListContextsResponse, error)
	CreateSession(ctx context.Context, req *chatagentv1.CreateSessionRequest) (*chatagentv1.CreateSessionResponse, error)
	GetSession(ctx context.Context, req *chatagentv1.GetSessionRequest) (*chatagentv1.GetSessionResponse, error)
	Turn(ctx context.Context, req *chatagentv1.TurnRequest) (*chatagentv1.TurnResponse, error)
	EndSession(ctx context.Context, req *chatagentv1.EndSessionRequest) (*chatagentv1.EndSessionResponse, error)
	IngestMessage(ctx context.Context, req *chatagentv1.IngestMessageRequest) (*chatagentv1.IngestMessageResponse, error)
}

type chatAgentBusiness struct {
	contexts        repository.ContextRepository
	sessions        repository.SessionRepository
	messages        repository.MessageRepository
	agent           *engine.Agent
	notificationCli notificationv1connect.NotificationServiceClient // existing Notification service client (optional)
}

// NewChatAgentBusiness wires repositories, the turn engine, and the existing Notification service client.
// notificationCli may be nil (RPC-only replies; no Notification.Send).
func NewChatAgentBusiness(
	contexts repository.ContextRepository,
	sessions repository.SessionRepository,
	messages repository.MessageRepository,
	llm engine.Completer,
	notificationCli notificationv1connect.NotificationServiceClient,
) ChatAgentBusiness {
	return &chatAgentBusiness{
		contexts:        contexts,
		sessions:        sessions,
		messages:        messages,
		agent:           engine.NewAgent(llm),
		notificationCli: notificationCli,
	}
}

func (b *chatAgentBusiness) UpsertContext(
	ctx context.Context,
	req *chatagentv1.UpsertContextRequest,
) (*chatagentv1.UpsertContextResponse, error) {
	def := protoToContextDef(req.GetDefinition())
	if vErr := validateContextDef(def); vErr != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, vErr)
	}
	ver, err := b.contexts.NextVersion(ctx, def.Key)
	if err != nil {
		return nil, fmt.Errorf("next version: %w", err)
	}
	snap, err := contextDefToJSONMap(def)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	row := &models.IntakeContext{
		ContextKey:     def.Key,
		Version:        ver,
		Purpose:        def.Purpose,
		DefinitionJSON: snap,
		Active:         true,
	}
	if cErr := b.contexts.Create(ctx, row); cErr != nil {
		return nil, fmt.Errorf("create context: %w", cErr)
	}
	return &chatagentv1.UpsertContextResponse{
		Definition: contextDefToProto(def),
		Version:    int32(ver),
	}, nil
}

func (b *chatAgentBusiness) GetContext(
	ctx context.Context,
	req *chatagentv1.GetContextRequest,
) (*chatagentv1.GetContextResponse, error) {
	key := strings.TrimSpace(req.GetContextKey())
	if key == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("context_key required"))
	}
	var row *models.IntakeContext
	var err error
	if req.GetVersion() > 0 {
		row, err = b.contexts.GetVersion(ctx, key, int(req.GetVersion()))
	} else {
		row, err = b.contexts.GetLatest(ctx, key)
	}
	if err != nil {
		if data.ErrorIsNoRows(err) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("context not found"))
		}
		return nil, err
	}
	def, err := contextDefFromJSONMap(row.DefinitionJSON)
	if err != nil {
		return nil, fmt.Errorf("decode context: %w", err)
	}
	return &chatagentv1.GetContextResponse{
		Definition: contextDefToProto(def),
		Version:    int32(row.Version),
	}, nil
}

func (b *chatAgentBusiness) ListContexts(
	ctx context.Context,
	_ *chatagentv1.ListContextsRequest,
) (*chatagentv1.ListContextsResponse, error) {
	rows, err := b.contexts.ListLatest(ctx)
	if err != nil {
		return nil, err
	}
	out := &chatagentv1.ListContextsResponse{}
	for _, row := range rows {
		out.Contexts = append(out.Contexts, &chatagentv1.ContextSummary{
			ContextKey: row.ContextKey,
			Version:    int32(row.Version),
			Purpose:    row.Purpose,
			Active:     row.Active,
		})
	}
	return out, nil
}

func (b *chatAgentBusiness) CreateSession(
	ctx context.Context,
	req *chatagentv1.CreateSessionRequest,
) (*chatagentv1.CreateSessionResponse, error) {
	subject := strings.TrimSpace(req.GetSubjectId())
	if subject == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("subject_id required"))
	}

	def, version, err := b.resolveConfig(ctx, req.GetContextKey(), int(req.GetContextVersion()), req.GetInlineConfig())
	if err != nil {
		return nil, err
	}
	if vErr := validateContextDef(def); vErr != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, vErr)
	}

	snap, err := contextDefToJSONMap(def)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	seed := structToFields(req.GetSeedFields())
	docs := protoDocs(req.GetDocuments())
	seedMsgs := protoMessages(req.GetSeedMessages())
	runtime := structToFields(req.GetRuntime())
	notifTarget := normalizeNotificationTarget(req.GetNotification(), subject)

	fields := engine.Sanitize(def, engine.MergeFields(nil, seed))
	fields = engine.ApplyDocuments(def, fields, docs)

	sess := &models.Session{
		SubjectID:      subject,
		ContextKey:     def.Key,
		ContextVersion: version,
		ConfigSnapshot: snap,
		Fields:         fieldsToJSONMap(fields),
		Runtime:        fieldsToJSONMap(runtime),
		Documents:      docsToJSONMap(docs),
		// Persist NotificationTarget (column names kept for schema compatibility).
		Channel:     notificationTargetToJSONMap(notifTarget),
		ChannelName: notificationType(notifTarget),
		ContactID:   notificationContactID(notifTarget),
		Status:      models.SessionStatusActive,
		Ready:       false,
	}

	// Optional immediate evaluation of seed evidence (no chat message yet).
	var reply string
	var source string
	var status map[string]engine.FieldStatus
	var missing []string
	if req.GetEvaluateEvidence() || len(docs) > 0 || len(seed) > 0 {
		res, terr := b.agent.Evaluate(ctx, def, fields, engine.Evidence{
			SeedFields: seed,
			Documents:  docs,
			Messages:   seedMsgs,
		})
		if terr != nil {
			return nil, terr
		}
		fields = res.Fields
		status = res.FieldStatus
		missing = res.Missing
		reply = res.Reply
		source = res.Source
		sess.Fields = fieldsToJSONMap(fields)
		sess.Ready = res.Ready
		if res.Ready {
			sess.Status = models.SessionStatusReady
		}
		// Append evaluation assistant message if any.
		if reply != "" {
			seedMsgs = append(seedMsgs, engine.Message{Role: "assistant", Content: reply})
		}
	} else {
		status, missing, sess.Ready = engine.Assess(def, fields)
		if sess.Ready {
			sess.Status = models.SessionStatusReady
		}
	}

	if sErr := b.sessions.Create(ctx, sess); sErr != nil {
		return nil, fmt.Errorf("create session: %w", sErr)
	}

	// Persist seed transcript.
	seq := 1
	for _, m := range seedMsgs {
		msg := &models.Message{
			SessionID: sess.ID,
			Seq:       seq,
			Role:      m.Role,
			Content:   m.Content,
		}
		if mErr := b.messages.Create(ctx, msg); mErr != nil {
			util.Log(ctx).WithError(mErr).Warn("chatagent: seed message persist failed")
		}
		seq++
	}
	sess.MessageCount = len(seedMsgs)
	if _, uErr := b.sessions.Update(ctx, sess, "message_count", "fields", "ready", "status"); uErr != nil {
		util.Log(ctx).WithError(uErr).Warn("chatagent: session message_count update failed")
	}

	api, err := b.toAPISession(ctx, sess, def, fields, status, missing, seedMsgs)
	if err != nil {
		return nil, err
	}
	delivered := false
	if reply != "" {
		delivered = b.deliverReply(ctx, notifTarget, sess.SubjectID, sess.ID, reply)
	}
	_ = source
	return &chatagentv1.CreateSessionResponse{
		Session:   api,
		Reply:     reply,
		Delivered: delivered,
	}, nil
}

func (b *chatAgentBusiness) GetSession(
	ctx context.Context,
	req *chatagentv1.GetSessionRequest,
) (*chatagentv1.GetSessionResponse, error) {
	sess, def, fields, msgs, err := b.loadSessionBundle(ctx, req.GetSessionId())
	if err != nil {
		return nil, err
	}
	status, missing, _ := engine.Assess(def, fields)
	api, err := b.toAPISession(ctx, sess, def, fields, status, missing, msgs)
	if err != nil {
		return nil, err
	}
	return &chatagentv1.GetSessionResponse{Session: api}, nil
}

func (b *chatAgentBusiness) Turn(
	ctx context.Context,
	req *chatagentv1.TurnRequest,
) (*chatagentv1.TurnResponse, error) {
	sess, def, fields, msgs, err := b.loadSessionBundle(ctx, req.GetSessionId())
	if err != nil {
		return nil, err
	}
	if sess.Status == models.SessionStatusEnded {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("session ended"))
	}

	newDocs := protoDocs(req.GetDocuments())
	// Merge new docs into stored document bag.
	allDocs := append(docsFromJSONMap(sess.Documents), newDocs...)
	sess.Documents = docsToJSONMap(allDocs)

	evidence := engine.Evidence{
		SeedFields: fieldsFromJSONMap(sess.Fields),
		Documents:  allDocs,
		Messages:   msgs,
		Structured: structToFields(req.GetStructured()),
	}

	res, err := b.agent.Turn(ctx, def, engine.TurnInput{
		Message:     req.GetMessage(),
		Evidence:    evidence,
		PriorFields: fields,
	})
	if err != nil {
		return nil, err
	}

	// Append transcript.
	seq, err := b.messages.NextSeq(ctx, sess.ID)
	if err != nil {
		return nil, fmt.Errorf("next seq: %w", err)
	}
	for _, m := range res.Messages {
		row := &models.Message{
			SessionID: sess.ID,
			Seq:       seq,
			Role:      m.Role,
			Content:   m.Content,
		}
		if mErr := b.messages.Create(ctx, row); mErr != nil {
			return nil, fmt.Errorf("persist message: %w", mErr)
		}
		msgs = append(msgs, m)
		seq++
	}

	sess.Fields = fieldsToJSONMap(res.Fields)
	sess.Ready = res.Ready
	sess.MessageCount = len(msgs)
	if res.Ready {
		sess.Status = models.SessionStatusReady
	} else if sess.Status != models.SessionStatusEnded {
		sess.Status = models.SessionStatusActive
	}
	if _, uErr := b.sessions.Update(ctx, sess,
		"fields", "ready", "status", "message_count", "documents"); uErr != nil {
		return nil, fmt.Errorf("update session: %w", uErr)
	}

	api, err := b.toAPISession(ctx, sess, def, res.Fields, res.FieldStatus, res.Missing, msgs)
	if err != nil {
		return nil, err
	}
	notifTarget := notificationTargetFromJSONMap(sess.Channel)
	delivered := b.deliverReply(ctx, notifTarget, sess.SubjectID, sess.ID, res.Reply)
	return &chatagentv1.TurnResponse{
		Session:   api,
		Reply:     res.Reply,
		Source:    res.Source,
		Delivered: delivered,
	}, nil
}

func (b *chatAgentBusiness) IngestMessage(
	ctx context.Context,
	req *chatagentv1.IngestMessageRequest,
) (*chatagentv1.IngestMessageResponse, error) {
	notifTarget := normalizeNotificationTarget(req.GetNotification(), strings.TrimSpace(req.GetSubjectId()))
	if notifTarget == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("notification target required"))
	}
	msg := strings.TrimSpace(req.GetMessage())
	if msg == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("message required"))
	}

	sessionCreated := false
	var sess *models.Session
	var err error

	if id := strings.TrimSpace(req.GetSessionId()); id != "" {
		sess, err = b.sessions.GetByID(ctx, id)
		if err != nil {
			if data.ErrorIsNoRows(err) || errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, connect.NewError(connect.CodeNotFound, errors.New("session not found"))
			}
			return nil, err
		}
	} else {
		subject := strings.TrimSpace(req.GetSubjectId())
		if subject == "" {
			subject = notificationProfileID(notifTarget)
		}
		contextKey := strings.TrimSpace(req.GetContextKey())
		contactID := notificationContactID(notifTarget)
		if subject == "" && contactID == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				errors.New("session_id or subject_id/notification.recipient required"))
		}
		sess, err = b.sessions.GetActiveByChannel(ctx, subject, contextKey, notificationType(notifTarget), contactID)
		if err != nil {
			if !data.ErrorIsNoRows(err) && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			sess = nil
		}
		if sess == nil {
			if !req.GetCreateIfMissing() {
				return nil, connect.NewError(connect.CodeNotFound,
					errors.New("no active session for notification target; set create_if_missing or session_id"))
			}
			if subject == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument,
					errors.New("subject_id required to create session"))
			}
			createResp, cErr := b.CreateSession(ctx, &chatagentv1.CreateSessionRequest{
				SubjectId:        subject,
				ContextKey:       contextKey,
				InlineConfig:     req.GetInlineConfig(),
				SeedFields:       req.GetSeedFields(),
				Runtime:          req.GetRuntime(),
				EvaluateEvidence: req.GetEvaluateEvidence(),
				Notification:     notifTarget,
			})
			if cErr != nil {
				return nil, cErr
			}
			sessionCreated = true
			sess, err = b.sessions.GetByID(ctx, createResp.GetSession().GetId())
			if err != nil {
				return nil, err
			}
		}
	}

	stored := notificationTargetFromJSONMap(sess.Channel)
	merged := mergeNotificationTargets(stored, notifTarget)
	// Persist if inbound filled contact gaps.
	if notificationContactID(stored) == "" && notificationContactID(merged) != "" {
		sess.Channel = notificationTargetToJSONMap(merged)
		sess.ChannelName = notificationType(merged)
		sess.ContactID = notificationContactID(merged)
		if _, uErr := b.sessions.Update(ctx, sess, "channel", "channel_name", "contact_id"); uErr != nil {
			util.Log(ctx).WithError(uErr).Warn("chatagent: notification target update failed")
		}
	}

	turnResp, err := b.Turn(ctx, &chatagentv1.TurnRequest{
		SessionId: sess.ID,
		Message:   msg,
	})
	if err != nil {
		return nil, err
	}
	return &chatagentv1.IngestMessageResponse{
		Session:        turnResp.GetSession(),
		Reply:          turnResp.GetReply(),
		Source:         turnResp.GetSource(),
		Delivered:      turnResp.GetDelivered(),
		SessionCreated: sessionCreated,
	}, nil
}

func (b *chatAgentBusiness) EndSession(
	ctx context.Context,
	req *chatagentv1.EndSessionRequest,
) (*chatagentv1.EndSessionResponse, error) {
	sess, def, fields, msgs, err := b.loadSessionBundle(ctx, req.GetSessionId())
	if err != nil {
		return nil, err
	}
	sess.Status = models.SessionStatusEnded
	if _, uErr := b.sessions.Update(ctx, sess, "status"); uErr != nil {
		return nil, fmt.Errorf("end session: %w", uErr)
	}
	status, missing, _ := engine.Assess(def, fields)
	api, err := b.toAPISession(ctx, sess, def, fields, status, missing, msgs)
	if err != nil {
		return nil, err
	}
	return &chatagentv1.EndSessionResponse{Session: api}, nil
}

func (b *chatAgentBusiness) resolveConfig(
	ctx context.Context,
	contextKey string,
	version int,
	inline *chatagentv1.ContextDefinition,
) (engine.ContextDef, int, error) {
	var base engine.ContextDef
	ver := version
	key := strings.TrimSpace(contextKey)
	if key != "" {
		var row *models.IntakeContext
		var err error
		if version > 0 {
			row, err = b.contexts.GetVersion(ctx, key, version)
		} else {
			row, err = b.contexts.GetLatest(ctx, key)
		}
		if err != nil {
			if inline == nil {
				if data.ErrorIsNoRows(err) || errors.Is(err, gorm.ErrRecordNotFound) {
					return engine.ContextDef{}, 0, connect.NewError(connect.CodeNotFound, errors.New("context not found"))
				}
				return engine.ContextDef{}, 0, err
			}
			// Registry miss with inline fallback.
		} else {
			base, err = contextDefFromJSONMap(row.DefinitionJSON)
			if err != nil {
				return engine.ContextDef{}, 0, fmt.Errorf("decode registry context: %w", err)
			}
			ver = row.Version
		}
	}
	if inline != nil {
		over := protoToContextDef(inline)
		base = mergeContextDefs(base, over)
		if ver == 0 {
			ver = 1
		}
	}
	if base.Key == "" {
		return engine.ContextDef{}, 0, connect.NewError(connect.CodeInvalidArgument,
			errors.New("context_key or inline_config required"))
	}
	return base, ver, nil
}

// mergeContextDefs applies non-empty overlay fields onto base (hybrid config).
func mergeContextDefs(base, over engine.ContextDef) engine.ContextDef {
	if over.Key != "" {
		base.Key = over.Key
	}
	if over.Purpose != "" {
		base.Purpose = over.Purpose
	}
	if over.SystemPrompt != "" {
		base.SystemPrompt = over.SystemPrompt
	}
	if over.ExtractRules != "" {
		base.ExtractRules = over.ExtractRules
	}
	if len(over.Fields) > 0 {
		base.Fields = over.Fields
	}
	if over.ReplyPolicy.MaxSentences > 0 || over.ReplyPolicy.CompleteMessage != "" {
		base.ReplyPolicy = over.ReplyPolicy
	}
	return base
}

func validateContextDef(def engine.ContextDef) error {
	if strings.TrimSpace(def.Key) == "" {
		return errors.New("context_key required")
	}
	if strings.TrimSpace(def.Purpose) == "" && strings.TrimSpace(def.SystemPrompt) == "" {
		return errors.New("purpose or system_prompt required")
	}
	if len(def.Fields) == 0 {
		return errors.New("at least one field required")
	}
	seen := map[string]struct{}{}
	for _, f := range def.Fields {
		if strings.TrimSpace(f.Name) == "" {
			return errors.New("field name required")
		}
		if _, ok := seen[f.Name]; ok {
			return fmt.Errorf("duplicate field %q", f.Name)
		}
		seen[f.Name] = struct{}{}
	}
	return nil
}

func (b *chatAgentBusiness) loadSessionBundle(
	ctx context.Context,
	sessionID string,
) (*models.Session, engine.ContextDef, engine.Fields, []engine.Message, error) {
	id := strings.TrimSpace(sessionID)
	if id == "" {
		return nil, engine.ContextDef{}, nil, nil, connect.NewError(connect.CodeInvalidArgument, errors.New("session_id required"))
	}
	sess, err := b.sessions.GetByID(ctx, id)
	if err != nil {
		if data.ErrorIsNoRows(err) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, engine.ContextDef{}, nil, nil, connect.NewError(connect.CodeNotFound, errors.New("session not found"))
		}
		return nil, engine.ContextDef{}, nil, nil, err
	}
	def, err := contextDefFromJSONMap(sess.ConfigSnapshot)
	if err != nil {
		return nil, engine.ContextDef{}, nil, nil, fmt.Errorf("decode config snapshot: %w", err)
	}
	fields := fieldsFromJSONMap(sess.Fields)
	rows, err := b.messages.ListBySession(ctx, sess.ID)
	if err != nil {
		return nil, engine.ContextDef{}, nil, nil, err
	}
	msgs := make([]engine.Message, 0, len(rows))
	for _, r := range rows {
		msgs = append(msgs, engine.Message{Role: r.Role, Content: r.Content})
	}
	return sess, def, fields, msgs, nil
}

func (b *chatAgentBusiness) toAPISession(
	_ context.Context,
	sess *models.Session,
	def engine.ContextDef,
	fields engine.Fields,
	status map[string]engine.FieldStatus,
	missing []string,
	msgs []engine.Message,
) (*chatagentv1.ChatSession, error) {
	return &chatagentv1.ChatSession{
		Id:             sess.ID,
		SubjectId:      sess.SubjectID,
		ContextKey:     sess.ContextKey,
		ContextVersion: int32(sess.ContextVersion),
		ConfigSnapshot: contextDefToProto(def),
		Fields:         fieldsToStruct(fields),
		Messages:       engineMessagesToProto(msgs),
		Ready:          sess.Ready,
		Status:         sessionStatusProto(sess.Status, sess.Ready),
		Missing:        missing,
		FieldStatus:    fieldStatusToProto(status),
		Runtime:        fieldsToStruct(fieldsFromJSONMap(sess.Runtime)),
		Notification:   notificationTargetFromJSONMap(sess.Channel),
	}, nil
}
