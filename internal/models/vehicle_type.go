package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type VehicleType struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`
	Text        string `json:"text"`
	Description string `json:"description"`
	State       int8   `json:"state"`
	OldId       uint64 `json:"old_id"`
}

type VehicleTypeRepository interface {
	Insert(ctx context.Context, u *VehicleType) error
	UpdateById(ctx context.Context, u *VehicleType) error
	GetById(ctx context.Context, id uint64) (*VehicleType, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *VehicleType) error
}
