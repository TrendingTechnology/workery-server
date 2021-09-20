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

type ActivitySheetItem struct {
	Id                   uint64      `json:"id"`
	Uuid                 string      `json:"uuid"`
	TenantId             uint64      `json:"tenant_id"`
	Comment              string      `json:"comment"`
	CreatedTime          time.Time   `json:"created_time"`
	CreatedById          null.Int    `json:"created_by_id"`
	CreatedByName        null.String `json:"created_by_name"`
	CreatedFromIP        null.String `json:"created_from_ip"`
	AssociateId          uint64      `json:"associate_id"`
	AssociateName        string      `json:"associate_name"`
	AssociateLexicalName string      `json:"associate_lexical_name"`
	OrderId              null.Int    `json:"order_id"`
	State                int8        `json:"state"`
	OngoingOrderId       null.Int    `json:"ongoing_order_id"`
	OldId                uint64      `json:"old_id"`
}

type ActivitySheetItemRepository interface {
	Insert(ctx context.Context, u *ActivitySheetItem) error
	UpdateById(ctx context.Context, u *ActivitySheetItem) error
	GetById(ctx context.Context, id uint64) (*ActivitySheetItem, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *ActivitySheetItem) error
}
