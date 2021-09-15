package models

import (
	"context"
	// "time"
	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteSkillSet` model.
type LiteSkillSetFilter struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type LiteSkillSet struct {
	Id          uint64 `json:"id"`
	TenantId    uint64 `json:"tenant_id"`
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	State       int8   `json:"state"`
}

type LiteSkillSetRepository interface {
	ListByFilter(ctx context.Context, filter *LiteSkillSetFilter) ([]*LiteSkillSet, error)
	CountByFilter(ctx context.Context, filter *LiteSkillSetFilter) (uint64, error)
}
