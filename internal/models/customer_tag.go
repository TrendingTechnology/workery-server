package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type CustomerTag struct {
	Id         uint64 `json:"id"`
	Uuid       string `json:"uuid"`
	TenantId   uint64 `json:"tenant_id"`
	CustomerId uint64 `json:"customer_id"`
	TagId      uint64 `json:"tag_id"`
	OldId      uint64 `json:"old_id"`
}

type CustomerTagRepository interface {
	Insert(ctx context.Context, u *CustomerTag) error
	UpdateById(ctx context.Context, u *CustomerTag) error
	GetById(ctx context.Context, id uint64) (*CustomerTag, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *CustomerTag) error
}
