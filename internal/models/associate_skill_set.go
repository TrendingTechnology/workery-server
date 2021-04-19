package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type AssociateSkillSet struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`
	AssociateId uint64 `json:"associate_id"`
	SkillSetId  uint64 `json:"skill_set_id"`
	OldId       uint64 `json:"old_id"`
}

type AssociateSkillSetRepository interface {
	Insert(ctx context.Context, u *AssociateSkillSet) error
	UpdateById(ctx context.Context, u *AssociateSkillSet) error
	GetById(ctx context.Context, id uint64) (*AssociateSkillSet, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*AssociateSkillSet, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *AssociateSkillSet) error
}
