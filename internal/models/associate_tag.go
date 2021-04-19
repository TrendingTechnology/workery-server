package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type AssociateTag struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`
	AssociateId uint64 `json:"associate_id"`
	TagId       uint64 `json:"tag_id"`
	OldId       uint64 `json:"old_id"`
}

type AssociateTagRepository interface {
	Insert(ctx context.Context, u *AssociateTag) error
	UpdateById(ctx context.Context, u *AssociateTag) error
	GetById(ctx context.Context, id uint64) (*AssociateTag, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*AssociateTag, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *AssociateTag) error
}
