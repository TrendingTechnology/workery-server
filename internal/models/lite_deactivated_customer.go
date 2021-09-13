package models

import (
	"context"
	// "time"
	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteDeactivatedCustomer` model.
type LiteDeactivatedCustomerFilter struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type LiteDeactivatedCustomer struct {
	Id                      uint64    `json:"id"`
	TenantId                uint64    `json:"tenant_id"`
	Name                    string    `json:"name,omitempty"`
	LexicalName             string    `json:"lexical_name,omitempty"`
	DeactivationReason      int8      `json:"deactivation_reason"`
	DeactivationReasonOther string    `json:"deactivation_reason_other"`
	State                   int8      `json:"state"`
}

type LiteDeactivatedCustomerRepository interface {
	ListByFilter(ctx context.Context, filter *LiteDeactivatedCustomerFilter) ([]*LiteDeactivatedCustomer, error)
	CountByFilter(ctx context.Context, filter *LiteDeactivatedCustomerFilter) (uint64, error)
}
