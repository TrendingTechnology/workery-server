package models

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteAssociate` model.
type LiteAssociateFilter struct {
	TenantId   uint64      `json:"tenant_id"`
	States     []int8      `json:"states"`
	SortOrder  string      `json:"sort_order"`
	SortField  string      `json:"sort_field"`
	Search     null.String `json:"search"`
	Offset     uint64      `json:"last_seen_id"`
	Limit      uint64      `json:"limit"`
}

type LiteAssociate struct {
	Id                 uint64 `json:"id"`
	TenantId           uint64 `json:"tenant_id"`
	State              int8 `json:"state"`
	GivenName          string       `json:"given_name"`
	LastName           string       `json:"last_name"`
	Name               string       `json:"name,omitempty"`
	LexicalName        string       `json:"lexical_name,omitempty"`
	Telephone          string `json:"telephone"`
	TelephoneTypeOf    int8   `json:"telephone_type_of"`
	TelephoneExtension string `json:"telephone_extension"`
	Email              string `json:"email"`
	JoinDate           null.Time `json:"join_date"`
}

type LiteAssociateRepository interface {
	ListByFilter(ctx context.Context, filter *LiteAssociateFilter) ([]*LiteAssociate, error)
	CountByFilter(ctx context.Context, filter *LiteAssociateFilter) (uint64, error)
}
