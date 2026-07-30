package repository

import (
	"context"
	"fmt"

	"github.com/pitabwire/frame/v2/datastore"
	"github.com/pitabwire/frame/v2/datastore/pool"
	"github.com/pitabwire/frame/v2/workerpool"

	"github.com/antinvestor/service-profile/apps/chatagent/service/models"
)

// SessionRepository stores intake sessions.
type SessionRepository interface {
	datastore.BaseRepository[*models.Session]
	GetByID(ctx context.Context, id string) (*models.Session, error)
	GetActiveBySubjectContext(ctx context.Context, subjectID, contextKey string) (*models.Session, error)
	// GetActiveByChannel finds the latest active/ready session for a subject+context on a channel contact.
	GetActiveByChannel(ctx context.Context, subjectID, contextKey, channelName, contactID string) (*models.Session, error)
}

type sessionRepository struct {
	datastore.BaseRepository[*models.Session]
}

func NewSessionRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) SessionRepository {
	return &sessionRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Session](
			ctx, dbPool, workMan, func() *models.Session { return &models.Session{} },
		),
	}
}

func (r *sessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	var row models.Session
	if err := r.Pool().DB(ctx, true).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *sessionRepository) GetActiveBySubjectContext(
	ctx context.Context,
	subjectID, contextKey string,
) (*models.Session, error) {
	var row models.Session
	err := r.Pool().DB(ctx, true).
		Where("subject_id = ? AND context_key = ? AND status IN ?",
			subjectID, contextKey, []string{models.SessionStatusActive, models.SessionStatusReady}).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *sessionRepository) GetActiveByChannel(
	ctx context.Context,
	subjectID, contextKey, channelName, contactID string,
) (*models.Session, error) {
	var row models.Session
	q := r.Pool().DB(ctx, true).
		Where("status IN ?", []string{models.SessionStatusActive, models.SessionStatusReady})
	if subjectID != "" {
		q = q.Where("subject_id = ?", subjectID)
	}
	if contextKey != "" {
		q = q.Where("context_key = ?", contextKey)
	}
	if channelName != "" {
		q = q.Where("channel_name = ?", channelName)
	}
	if contactID != "" {
		q = q.Where("contact_id = ?", contactID)
	}
	err := q.Order("created_at DESC").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// MessageRepository stores transcript rows.
type MessageRepository interface {
	datastore.BaseRepository[*models.Message]
	ListBySession(ctx context.Context, sessionID string) ([]*models.Message, error)
	NextSeq(ctx context.Context, sessionID string) (int, error)
}

type messageRepository struct {
	datastore.BaseRepository[*models.Message]
}

func NewMessageRepository(ctx context.Context, dbPool pool.Pool, workMan workerpool.Manager) MessageRepository {
	return &messageRepository{
		BaseRepository: datastore.NewBaseRepository[*models.Message](
			ctx, dbPool, workMan, func() *models.Message { return &models.Message{} },
		),
	}
}

func (r *messageRepository) ListBySession(ctx context.Context, sessionID string) ([]*models.Message, error) {
	var rows []*models.Message
	err := r.Pool().DB(ctx, true).
		Where("session_id = ?", sessionID).
		Order("seq ASC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	return rows, nil
}

func (r *messageRepository) NextSeq(ctx context.Context, sessionID string) (int, error) {
	var maxSeq *int
	err := r.Pool().DB(ctx, true).Model(&models.Message{}).
		Select("MAX(seq)").
		Where("session_id = ?", sessionID).
		Scan(&maxSeq).Error
	if err != nil {
		return 0, err
	}
	if maxSeq == nil {
		return 1, nil
	}
	return *maxSeq + 1, nil
}
