package dbutil

import (
	"context"
	"time"
)

const defaultBatchSize = 30

type SearchQuery struct {
	ProfileID            string
	Query                string
	PropertiesToSearchOn []string
	StartAt              *time.Time
	EndAt                *time.Time

	Pagination *Paginator
}

func NewSearchQuery(
	_ context.Context,
	profileID, query string,
	props []string,
	startAt, endAt string,
	resultPage, resultCount int,
) (*SearchQuery, error) {
	if resultCount == 0 {
		resultCount = defaultBatchSize
	}

	sq := &SearchQuery{
		ProfileID:            profileID,
		Query:                query,
		PropertiesToSearchOn: props,
		Pagination: &Paginator{
			Offset:    resultPage * resultCount,
			Limit:     resultCount,
			BatchSize: defaultBatchSize,
		},
	}

	if startAt != "" {
		parsedTime, err := time.Parse(time.DateTime, startAt)
		if err != nil {
			return nil, err
		}
		sq.StartAt = &parsedTime
	}

	if endAt != "" {
		parsedTime, err := time.Parse(time.DateTime, endAt)
		if err != nil {
			return nil, err
		}
		sq.EndAt = &parsedTime
	}

	return sq, nil
}

type Paginator struct {
	Offset int
	Limit  int

	BatchSize int
}

func (sq *Paginator) CanLoad() bool {
	return sq.Offset < sq.Limit
}

func (sq *Paginator) Stop(loadedCount int) bool {
	sq.Offset += loadedCount
	if sq.Offset+sq.BatchSize > sq.Limit {
		sq.BatchSize = sq.Limit - sq.Offset
	}

	return loadedCount < sq.BatchSize
}
