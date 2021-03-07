package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type SkillSetInsuranceRequirement struct {
	Id                uint64    `json:"id"`
	Uuid              string    `json:"uuid"`
	TenantId          uint64    `json:"tenant_id"`
	SkillSetId             uint64    `json:"skill_set_id"`
	InsuranceRequirementId uint64 `json:"insurance_requirement_id"`
    OldId                  uint64    `json:"old_id"`
}

type SkillSetInsuranceRequirementRepository interface {
	Insert(ctx context.Context, u *SkillSetInsuranceRequirement) error
	UpdateById(ctx context.Context, u *SkillSetInsuranceRequirement) error
	GetById(ctx context.Context, id uint64) (*SkillSetInsuranceRequirement, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*SkillSetInsuranceRequirement, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *SkillSetInsuranceRequirement) error
}
