package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteTaskItem` model.
type LiteTaskItemFilter struct {
	TenantId   uint64      `json:"tenant_id"`
	States     []int8      `json:"states"`
	SortOrder  string      `json:"sort_order"`
	SortField  string      `json:"sort_field"`
	IsClosed   null.Bool   `json:"is_closed"`
	Search     null.String `json:"search"`
	Offset     uint64      `json:"offset"`
	Limit      uint64      `json:"limit"`
}

type LiteTaskItem struct {
	Id                   uint64      `json:"id"`
	TenantId             uint64      `json:"tenant_id"`
	State                int8        `json:"state"`
	DueDate              time.Time   `json:"due_date"`
	TypeOf               int8        `json:"type_of"`
	CustomerId           null.Int    `json:"customer_id"`
	CustomerName         null.String `json:"customer_name,omitempty"`
	CustomerLexicalName  null.String `json:"customer_lexical_name,omitempty"`
	AssociateId          null.Int    `json:"associate_id"`
	AssociateName        null.String `json:"associate_name,omitempty"`
	AssociateLexicalName null.String `json:"associate_lexical_name,omitempty"`
	OrderTypeOf          int8        `json:"order_type_of"`
}

type LiteTaskItemRepository interface {
	ListByFilter(ctx context.Context, filter *LiteTaskItemFilter) ([]*LiteTaskItem, error)
	CountByFilter(ctx context.Context, filter *LiteTaskItemFilter) (uint64, error)
}
