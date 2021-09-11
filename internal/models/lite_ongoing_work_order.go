package models

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteOngoingWorkOrder` model.
type LiteOngoingWorkOrderFilter struct {
	TenantId             uint64      `json:"tenant_id"`
	States               []int8      `json:"states"`
	LastModifiedById     null.Int    `json:"last_modified_by_id"`
	CustomerName         null.String `json:"customer_name"`
	CustomerLexicalName  null.String `json:"customer_lexical_name"`
	AssociateName        null.String `json:"associate_name"`
	AssociateLexicalName null.String `json:"associate_lexical_name"`
	SortOrder            string      `json:"sort_order"`
	SortField            string      `json:"sort_field"`
	Search               null.String `json:"search"`
	Offset               uint64      `json:"offset"`
	Limit                uint64      `json:"limit"`
}

type LiteOngoingWorkOrder struct {
	Id                   uint64      `json:"id"`
	TenantId             uint64      `json:"tenant_id"`
	CustomerId           uint64      `json:"customer_id"`
	CustomerName         null.String `json:"customer_name"`
	CustomerLexicalName  null.String `json:"customer_lexical_name"`
	AssociateId          null.Int    `json:"associate_id"`
	AssociateName        null.String `json:"associate_name"`
	AssociateLexicalName null.String `json:"associate_lexical_name"`
	State                int8        `json:"state"`
}

type LiteOngoingWorkOrderRepository interface {
	ListByFilter(ctx context.Context, filter *LiteOngoingWorkOrderFilter) ([]*LiteOngoingWorkOrder, error)
	CountByFilter(ctx context.Context, filter *LiteOngoingWorkOrderFilter) (uint64, error)
}
