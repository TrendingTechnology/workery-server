package models

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteFinancial` model.
type LiteFinancialFilter struct {
	TenantId             uint64      `json:"tenant_id"`
	States               []int8      `json:"states"`
	LastModifiedById     null.Int    `json:"last_modified_by_id"`
	AssociateName        null.String `json:"associate_name"`
	AssociateLexicalName null.String `json:"associate_lexical_name"`
	CustomerName         null.String `json:"customer_name"`
	CustomerLexicalName  null.String `json:"customer_lexical_name"`
	SortOrder            string      `json:"sort_order"`
	SortField            string      `json:"sort_field"`
	Search               null.String `json:"search"`
	Offset               uint64      `json:"offset"`
	Limit                uint64      `json:"limit"`
}

type LiteFinancial struct {
	Id                           uint64      `json:"id"`
	TenantId                     uint64      `json:"tenant_id"`
	CustomerId                   uint64      `json:"customer_id"`
	CustomerName                 string      `json:"customer_name"`
	AssociateId                  null.Int    `json:"associate_id"`
	AssociateName                null.String `json:"associate_name"`
	InvoiceServiceFeePaymentDate null.Time   `json:"invoice_service_fee_payment_date"`
	TypeOf                       int8        `json:"type_of"`
	State                        int8        `json:"state"`
}

type LiteFinancialRepository interface {
	ListByFilter(ctx context.Context, filter *LiteFinancialFilter) ([]*LiteFinancial, error)
	CountByFilter(ctx context.Context, filter *LiteFinancialFilter) (uint64, error)
}
