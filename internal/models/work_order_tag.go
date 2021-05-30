package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type WorkOrderTag struct {
	Id       uint64 `json:"id"`
	Uuid     string `json:"uuid"`
	TenantId uint64 `json:"tenant_id"`
	OrderId  uint64 `json:"order_id"`
	TagId    uint64 `json:"tag_id"`
	OldId    uint64 `json:"old_id"`
}

type WorkOrderTagRepository interface {
	Insert(ctx context.Context, u *WorkOrderTag) error
	UpdateById(ctx context.Context, u *WorkOrderTag) error
	GetById(ctx context.Context, id uint64) (*WorkOrderTag, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*WorkOrderTag, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrderTag) error
}
