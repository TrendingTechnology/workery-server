package models

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteAssociateAwayLog` model.
type LiteAssociateAwayLogFilter struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type LiteAssociateAwayLog struct {
	Id                   uint64    `json:"id"`
	TenantId             uint64    `json:"tenant_id"`
	AssociateId          uint64    `json:"associate_id"`
	AssociateName        string    `json:"associate_name"`
	AssociateLexicalName string    `json:"associate_lexical_name"`
	StartDate            null.Time `json:"start_date"`
	UntilFurtherNotice   bool      `json:"until_further_notice"`
	UntilDate            null.Time `json:"until_date"`
	State                int8      `json:"state"`
}

type LiteAssociateAwayLogRepository interface {
	ListByFilter(ctx context.Context, filter *LiteAssociateAwayLogFilter) ([]*LiteAssociateAwayLog, error)
	CountByFilter(ctx context.Context, filter *LiteAssociateAwayLogFilter) (uint64, error)
}
