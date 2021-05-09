package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type WorkOrderSkillSet struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`
	OrderId     uint64 `json:"order_id"`
	SkillSetId uint64 `json:"skill_set_id"`
	OldId      uint64 `json:"old_id"`
}

type WorkOrderSkillSetRepository interface {
	Insert(ctx context.Context, u *WorkOrderSkillSet) error
	UpdateById(ctx context.Context, u *WorkOrderSkillSet) error
	GetById(ctx context.Context, id uint64) (*WorkOrderSkillSet, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*WorkOrderSkillSet, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrderSkillSet) error
}
