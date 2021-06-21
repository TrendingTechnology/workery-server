package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type StaffTag struct {
	Id         uint64 `json:"id"`
	Uuid       string `json:"uuid"`
	TenantId   uint64 `json:"tenant_id"`
	StaffId    uint64 `json:"staff_id"`
	TagId      uint64 `json:"tag_id"`
	OldId      uint64 `json:"old_id"`
}

type StaffTagRepository interface {
	Insert(ctx context.Context, u *StaffTag) error
	UpdateById(ctx context.Context, u *StaffTag) error
	GetById(ctx context.Context, id uint64) (*StaffTag, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *StaffTag) error
}
