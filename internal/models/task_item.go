package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

const (
	TaskActiveState = 1
	TaskInactiveState = 0
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type TaskItem struct {
	Id                 uint64      `json:"id"`                    // 01
	Uuid               string      `json:"uuid"`                  // 02
	TenantId           uint64      `json:"tenant_id"`             // 03
	TypeOf             string      `json:"type_of"`               // 04
	Title              string      `json:"title"`                 // 05
	Description        string      `json:"description"`           // 06
	DueDate            time.Time   `json:"due_date"`              // 07
	IsClosed           bool        `json:"is_closed"`             // 08
	WasPostponed       string      `json:"was_postponed"`         // 09
	ClosingReason      int8        `json:"closing_reason"`        // 10
	ClosingReasonOther string      `json:"closing_reason_other"`  // 11
	OrderId            uint64      `json:"order_id"`              // 12
	OngoingOrderId     null.Int    `json:"ongoing_order_id"`      // 13
	CreatedTime        time.Time   `json:"created_time"`          // 14
	CreatedFromIP      null.String `json:"created_from_ip"`       // 15
	CreatedById        null.Int    `json:"created_by_id"`         // 16
	LastModifiedTime   time.Time   `json:"last_modified_time"`    // 17
	LastModifiedFromIP null.String `json:"last_modified_from_ip"` // 18
	LastModifiedById   null.Int    `json:"last_modified_by_id"`   // 19
	State              int8        `json:"state"`                 // 20
	OldId              uint64      `json:"old_id"`                // 21
}

type TaskItemRepository interface {
	Insert(ctx context.Context, u *TaskItem) error
	UpdateById(ctx context.Context, u *TaskItem) error
	GetById(ctx context.Context, id uint64) (*TaskItem, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *TaskItem) error
}
