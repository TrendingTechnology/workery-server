package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type AssociateInsuranceRequirement struct {
	Id                     uint64 `json:"id"`
	Uuid                   string `json:"uuid"`
	TenantId               uint64 `json:"tenant_id"`
	AssociateId            uint64 `json:"associate_id"`
	InsuranceRequirementId uint64 `json:"insurance_requirement_id"`
	OldId                  uint64 `json:"old_id"`
}

type AssociateInsuranceRequirementRepository interface {
	Insert(ctx context.Context, u *AssociateInsuranceRequirement) error
	UpdateById(ctx context.Context, u *AssociateInsuranceRequirement) error
	GetById(ctx context.Context, id uint64) (*AssociateInsuranceRequirement, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *AssociateInsuranceRequirement) error
}
