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

type WorkOrderDeposit struct {
	Id                 uint64      `json:"id"`
	Uuid               string      `json:"uuid"`
	TenantId           uint64      `json:"tenant_id"`
	OrderId            uint64      `json:"order_id"`
	PaidAt             null.Time   `json:"paid_at"`
	DepositMethod      int8        `json:"deposit_method"`
	PaidTo             null.Int    `json:"paid_to"`
	Currency           string      `json:"currency"`
	Amount             float64     `json:"amount"`
	PaidFor            int8        `json:"paid_for"`
	CreatedTime        time.Time   `json:"created_time"`
	CreatedById        null.Int    `json:"created_by_id"`
	CreatedFromIP      null.String `json:"created_from_ip"`
	LastModifiedTime   time.Time   `json:"last_modified_time"`
	LastModifiedById   null.Int    `json:"last_modified_by_id"`
	LastModifiedFromIP null.String `json:"last_modified_from_ip"`
	State              int8        `json:"state"`
	OldId              uint64      `json:"old_id"`
}

type WorkOrderDepositRepository interface {
	Insert(ctx context.Context, u *WorkOrderDeposit) error
	UpdateById(ctx context.Context, u *WorkOrderDeposit) error
	GetById(ctx context.Context, id uint64) (*WorkOrderDeposit, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrderDeposit) error
}
