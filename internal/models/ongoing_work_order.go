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

type OngoingWorkOrder struct {
	Id                   uint64      `json:"id"`
	Uuid                 string      `json:"uuid"`
	TenantId             uint64      `json:"tenant_id"`
	CustomerId           uint64      `json:"customer_id"`
	CustomerName         null.String `json:"customer_name"`
	CustomerLexicalName  null.String `json:"customer_lexical_name"`
	AssociateId          null.Int    `json:"associate_id"`
	AssociateName        null.String `json:"associate_name"`
	AssociateLexicalName null.String `json:"associate_lexical_name"`
	State                int8        `json:"state"`
	CreatedTime          time.Time   `json:"created_time"`
	CreatedById          null.Int    `json:"created_by_id"`
	CreatedFromIP        null.String `json:"created_from_ip"`
	LastModifiedTime     time.Time   `json:"last_modified_time"`
	LastModifiedById     null.Int    `json:"last_modified_by_id"`
	LastModifiedFromIP   null.String `json:"last_modified_from_ip"`
	OldId                uint64      `json:"old_id"`
}

type OngoingWorkOrderRepository interface {
	Insert(ctx context.Context, u *OngoingWorkOrder) error
	UpdateById(ctx context.Context, u *OngoingWorkOrder) error
	GetById(ctx context.Context, id uint64) (*OngoingWorkOrder, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *OngoingWorkOrder) error
}
