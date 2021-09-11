package models

import (
	"context"
	// "time"
	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteStaff` model.
type LiteStaffFilter struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type LiteStaff struct {
	Id        uint64    `json:"id"`
	TenantId  uint64    `json:"tenant_id"`
	State     int8      `json:"state"`
	GivenName string    `json:"given_name"`
	LastName  string    `json:"last_name"`
	Telephone string    `json:"telephone"`
	Email     string    `json:"email"`
	JoinDate  null.Time `json:"join_date"`
	TypeOf    int8      `json:"type_of"`
}

type LiteStaffRepository interface {
	ListByFilter(ctx context.Context, filter *LiteStaffFilter) ([]*LiteStaff, error)
	CountByFilter(ctx context.Context, filter *LiteStaffFilter) (uint64, error)
}
