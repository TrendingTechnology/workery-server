package models

import (
	"context"
	// "time"

	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteTask` model.
type LiteTaskFilter struct {
	TenantId   uint64    `json:"tenant_id"`
	States     []int8    `json:"states"`
	IsClosed   null.Bool `json:"is_closed"`
	LastSeenId uint64    `json:"last_seen_id"`
	Limit      uint64    `json:"limit"`
}

type LiteTask struct {
	Id uint64 `json:"id"`
	// Uuid                              string      `json:"uuid"`
	TenantId uint64 `json:"tenant_id"`
	State    int8   `json:"state"`
}

type LiteTaskRepository interface {
	ListByFilter(ctx context.Context, filter *LiteTaskFilter) ([]*LiteTask, error)
	CountByFilter(ctx context.Context, filter *LiteTaskFilter) (uint64, error)
}
