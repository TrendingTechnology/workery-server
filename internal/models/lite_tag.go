package models

import (
	"context"
	// "time"
	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteTag` model.
type LiteTagFilter struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type LiteTag struct {
	Id          uint64 `json:"id"`
	TenantId    uint64 `json:"tenant_id"`
	Text        string `json:"text"`
	Description string `json:"description"`
	State       int8   `json:"state"`
}

type LiteTagRepository interface {
	ListByFilter(ctx context.Context, filter *LiteTagFilter) ([]*LiteTag, error)
	CountByFilter(ctx context.Context, filter *LiteTagFilter) (uint64, error)
}
