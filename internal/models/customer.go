package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type Customer struct {
	Id                uint64    `json:"id"`
	Uuid              string    `json:"uuid"`
	TenantId          uint64    `json:"tenant_id"`
	Text              string    `json:"text"`
	Description       string    `json:"description"`
	State             int8      `json:"state"`
    OldId             uint64    `json:"old_id"`
}

type CustomerRepository interface {
	Insert(ctx context.Context, u *Customer) error
	UpdateById(ctx context.Context, u *Customer) error
	GetById(ctx context.Context, id uint64) (*Customer, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*Customer, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *Customer) error
}
