package models

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

type CustomerTagFilter struct {
	TenantId   uint64      `json:"tenant_id"`
	States     []int8      `json:"states"`
	SortOrder  string      `json:"sort_order"`
	SortField  string      `json:"sort_field"`
	CustomerId null.Int    `json:"customer_id"`
	Search     null.String `json:"search"`
	Offset     uint64      `json:"offset"`
	Limit      uint64      `json:"limit"`
}

// State
//---------------------
// 1 = Active
// 0 = Inactive

type CustomerTag struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid,omitempty"`
	TenantId    uint64 `json:"tenant_id"`
	CustomerId  uint64 `json:"customer_id"`
	TagId       uint64 `json:"tag_id"`
	OldId       uint64 `json:"old_id"`
	Text        string `json:"text,omitempty"`        // Referenced value from 'tags'.
	Description string `json:"description,omitempty"` // Referenced value from 'tags'.
	State       int8   `json:"state,omitempty"`       // Referenced value from 'tags'.
}

type CustomerTagRepository interface {
	Insert(ctx context.Context, u *CustomerTag) error
	UpdateById(ctx context.Context, u *CustomerTag) error
	GetById(ctx context.Context, id uint64) (*CustomerTag, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *CustomerTag) error
	ListByFilter(ctx context.Context, filter *CustomerTagFilter) ([]*CustomerTag, error)
	CountByFilter(ctx context.Context, filter *CustomerTagFilter) (uint64, error)
}
