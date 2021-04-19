package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type WorkOrder struct {
	Id                 uint64      `json:"id"`
	Uuid               string      `json:"uuid"`
	TenantId           uint64      `json:"tenant_id"`
	CustomerId         uint64      `json:"customer_id"`
	AssociateId        null.Int    `json:"associate_id"`
	State              int8        `json:"state"`
	CreatedTime        time.Time   `json:"created_time"`
	CreatedById        null.Int    `json:"created_by_id"`
	CreatedFromIP      null.String `json:"created_from_ip"`
	LastModifiedTime   time.Time   `json:"last_modified_time"`
	LastModifiedById   null.Int    `json:"last_modified_by_id"`
	LastModifiedFromIP null.String `json:"last_modified_from_ip"`
	OldId              uint64      `json:"old_id"`
}

type WorkOrderRepository interface {
	Insert(ctx context.Context, u *WorkOrder) error
	UpdateById(ctx context.Context, u *WorkOrder) error
	GetById(ctx context.Context, id uint64) (*WorkOrder, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrder) error
}
