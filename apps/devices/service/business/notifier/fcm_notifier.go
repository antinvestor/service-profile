package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	devicev1 "github.com/antinvestor/apis/go/device/v1"
	"github.com/antinvestor/service-profile/apps/devices/config"
	"github.com/antinvestor/service-profile/apps/devices/service/business"
	"github.com/antinvestor/service-profile/apps/devices/service/repository"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"
)

const (
	fcmProvider           = "fcm"
	providerField         = "provider"
	tokenField            = "token"
	defaultFCMBatchSize   = 500
	defaultFCMEndpoint    = "https://fcm.googleapis.com/fcm/send"
	authorizationTemplate = "key=%s"
)

type fcmNotifier struct {
	cfg  *config.DevicesConfig
	keys business.KeysBusiness
}

func NewFCMNotifier(cfg *config.DevicesConfig, keys business.KeysBusiness, _ repository.DeviceRepository) Notifier {
	return &fcmNotifier{
		cfg:  cfg,
		keys: keys,
	}
}

func (f *fcmNotifier) Register(ctx context.Context, req *devicev1.RegisterKeyRequest) (*devicev1.KeyObject, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.GetKeyType() != devicev1.KeyType_FCM_TOKEN {
		return nil, fmt.Errorf("unsupported key type: %s", req.GetKeyType().String())
	}

	extras := map[string]any{
		providerField: fcmProvider,
	}

	if req.GetExtras() != nil {
		for k, v := range req.GetExtras().AsMap() {
			extras[k] = v
		}
	}

	rawToken, ok := extras[tokenField]
	if !ok {
		return nil, errors.New("extras.token must be provided for FCM registration")
	}

	token, ok := rawToken.(string)
	if !ok || strings.TrimSpace(token) == "" {
		return nil, errors.New("extras.token must be a non-empty string")
	}

	token = strings.TrimSpace(token)
	extras[tokenField] = token

	return f.keys.AddKey(ctx, strings.TrimSpace(req.GetDeviceId()), req.GetKeyType(), []byte(token), data.JSONMap(extras))
}

func (f *fcmNotifier) DeRegister(_ context.Context, _ *devicev1.KeyObject) error {
	return nil
}

func (f *fcmNotifier) Notify(ctx context.Context, req *devicev1.NotifyRequest, keys []*devicev1.KeyObject) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	serverKey := strings.TrimSpace(f.cfg.FCMServerKey)
	if serverKey == "" {
		return errors.New("fcm server key is not configured")
	}

	var tokens []string
	for _, key := range keys {
		if key == nil {
			continue
		}

		token := strings.TrimSpace(string(key.GetKey()))
		extraMap := map[string]any{}
		if key.GetExtra() != nil {
			extraMap = key.GetExtra().AsMap()
		}

		if provider, ok := extraMap[providerField].(string); ok && provider != "" && !strings.EqualFold(provider, fcmProvider) {
			continue
		}

		if extraToken, ok := extraMap[tokenField].(string); ok && strings.TrimSpace(extraToken) != "" {
			token = strings.TrimSpace(extraToken)
		}

		if token != "" {
			tokens = append(tokens, token)
		}
	}

	if len(tokens) == 0 {
		return errors.New("no notification tokens registered for device")
	}

	dataPayload := map[string]any{}
	if req.GetData() != nil {
		dataPayload = req.GetData().AsMap()
	}

	extraPayload := map[string]any{}
	if req.GetExtras() != nil {
		extraPayload = req.GetExtras().AsMap()
	}

	batchSize := f.batchSize()
	for start := 0; start < len(tokens); start += batchSize {
		end := start + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		if err := f.sendBatch(ctx, serverKey, tokens[start:end], req.GetTitle(), req.GetBody(), dataPayload, extraPayload); err != nil {
			return err
		}
	}

	return nil
}

func (f *fcmNotifier) batchSize() int {
	if f.cfg != nil && f.cfg.NotificationBatchSize > 0 {
		return f.cfg.NotificationBatchSize
	}
	return defaultFCMBatchSize
}

func (f *fcmNotifier) endpoint() string {
	if f.cfg != nil && strings.TrimSpace(f.cfg.FCMEndpoint) != "" {
		return strings.TrimSpace(f.cfg.FCMEndpoint)
	}
	return defaultFCMEndpoint
}

func (f *fcmNotifier) sendBatch(
	ctx context.Context,
	serverKey string,
	tokens []string,
	title, body string,
	dataPayload map[string]any,
	extraPayload map[string]any,
) error {
	payload := map[string]any{
		"registration_ids": tokens,
	}

	if title != "" || body != "" {
		payload["notification"] = map[string]string{
			"title": title,
			"body":  body,
		}
	}

	if len(dataPayload) > 0 {
		payload["data"] = dataPayload
	}

	for k, v := range extraPayload {
		if _, exists := payload[k]; exists {
			continue
		}
		payload[k] = v
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, f.endpoint(), bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf(authorizationTemplate, serverKey))

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		util.Log(ctx).
			WithField("token_count", len(tokens)).
			WithField("endpoint", f.endpoint()).
			Debug("notification batch sent via FCM")
		return nil
	}

	respBody, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf("fcm notification failed with status %d", response.StatusCode)
	}

	return fmt.Errorf("fcm notification failed: status=%d body=%s", response.StatusCode, string(respBody))
}
