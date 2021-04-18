package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type AssociateVehicleType struct {
	Id                     uint64 `json:"id"`
	Uuid                   string `json:"uuid"`
	TenantId               uint64 `json:"tenant_id"`
	AssociateId            uint64 `json:"associate_id"`
	VehicleTypeId uint64 `json:"vehicle_type_id"`
	OldId                  uint64 `json:"old_id"`
}

type AssociateVehicleTypeRepository interface {
	Insert(ctx context.Context, u *AssociateVehicleType) error
	UpdateById(ctx context.Context, u *AssociateVehicleType) error
	GetById(ctx context.Context, id uint64) (*AssociateVehicleType, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*AssociateVehicleType, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *AssociateVehicleType) error
}
