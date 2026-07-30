package handlers

import (
	"context"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2"
	"github.com/pitabwire/frame/v2/datastore"

	"github.com/antinvestor/service-profile/apps/chatagent/service/business"
	"github.com/antinvestor/service-profile/apps/chatagent/service/engine"
	"github.com/antinvestor/service-profile/apps/chatagent/service/llm"
	"github.com/antinvestor/service-profile/apps/chatagent/service/repository"
	chatagentv1 "github.com/antinvestor/service-profile/gen/go/chatagent/v1"
	"github.com/antinvestor/service-profile/gen/go/chatagent/v1/chatagentv1connect"
	"github.com/antinvestor/service-profile/pkg/errorutil"
)

// ChatAgentServer implements Connect ChatAgentService.
type ChatAgentServer struct {
	biz business.ChatAgentBusiness
	chatagentv1connect.UnimplementedChatAgentServiceHandler
}

// LLMConfig configures optional inference for the handler stack.
type LLMConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

// ServerDeps optional dependencies for ChatAgentServer.
// NotificationClient is the existing Notification service client (same pattern as profile app).
type ServerDeps struct {
	LLM                LLMConfig
	NotificationClient notificationv1connect.NotificationServiceClient
}

// NewChatAgentServer builds the handler with repositories, optional LLM, and Notification client.
func NewChatAgentServer(ctx context.Context, svc *frame.Service, deps ServerDeps) *ChatAgentServer {
	workMan := svc.WorkManager()
	dbPool := svc.DatastoreManager().GetPool(ctx, datastore.DefaultPoolName)

	ctxRepo := repository.NewContextRepository(ctx, dbPool, workMan)
	sessRepo := repository.NewSessionRepository(ctx, dbPool, workMan)
	msgRepo := repository.NewMessageRepository(ctx, dbPool, workMan)

	var completer engine.Completer
	if deps.LLM.BaseURL != "" {
		httpClient := svc.HTTPClientManager().Client(ctx)
		completer = llm.New(deps.LLM.BaseURL, deps.LLM.APIKey, deps.LLM.Model, httpClient)
	}

	return &ChatAgentServer{
		biz: business.NewChatAgentBusiness(ctxRepo, sessRepo, msgRepo, completer, deps.NotificationClient),
	}
}

func (s *ChatAgentServer) UpsertContext(
	ctx context.Context,
	req *connect.Request[chatagentv1.UpsertContextRequest],
) (*connect.Response[chatagentv1.UpsertContextResponse], error) {
	resp, err := s.biz.UpsertContext(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

func (s *ChatAgentServer) GetContext(
	ctx context.Context,
	req *connect.Request[chatagentv1.GetContextRequest],
) (*connect.Response[chatagentv1.GetContextResponse], error) {
	resp, err := s.biz.GetContext(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

func (s *ChatAgentServer) ListContexts(
	ctx context.Context,
	req *connect.Request[chatagentv1.ListContextsRequest],
) (*connect.Response[chatagentv1.ListContextsResponse], error) {
	resp, err := s.biz.ListContexts(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

func (s *ChatAgentServer) CreateSession(
	ctx context.Context,
	req *connect.Request[chatagentv1.CreateSessionRequest],
) (*connect.Response[chatagentv1.CreateSessionResponse], error) {
	resp, err := s.biz.CreateSession(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

func (s *ChatAgentServer) GetSession(
	ctx context.Context,
	req *connect.Request[chatagentv1.GetSessionRequest],
) (*connect.Response[chatagentv1.GetSessionResponse], error) {
	resp, err := s.biz.GetSession(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

func (s *ChatAgentServer) Turn(
	ctx context.Context,
	req *connect.Request[chatagentv1.TurnRequest],
) (*connect.Response[chatagentv1.TurnResponse], error) {
	resp, err := s.biz.Turn(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

func (s *ChatAgentServer) EndSession(
	ctx context.Context,
	req *connect.Request[chatagentv1.EndSessionRequest],
) (*connect.Response[chatagentv1.EndSessionResponse], error) {
	resp, err := s.biz.EndSession(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}

func (s *ChatAgentServer) IngestMessage(
	ctx context.Context,
	req *connect.Request[chatagentv1.IngestMessageRequest],
) (*connect.Response[chatagentv1.IngestMessageResponse], error) {
	resp, err := s.biz.IngestMessage(ctx, req.Msg)
	if err != nil {
		return nil, errorutil.CleanErr(err)
	}
	return connect.NewResponse(resp), nil
}
