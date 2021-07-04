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

type PrivateFile struct {
	Id                 uint64      `json:"id"`                    // 1
	Uuid               string      `json:"uuid"`                  // 2
	TenantId           uint64      `json:"tenant_id"`             // 3
	S3Key              string      `json:"s3_key"`                // 4
	Title              string      `json:"title"`                 // 5
	Description        string      `json:"description"`           // 6
	IndexedText        string      `json:"indexed_text"`          // 7
	CreatedTime        time.Time   `json:"created_time"`          // 8
	CreatedFromIP      null.String `json:"created_from_ip"`       // 9
	CreatedById        null.Int    `json:"created_by_id"`         // 10
	LastModifiedTime   time.Time   `json:"last_modified_time"`    // 11
	LastModifiedById   null.Int    `json:"last_modified_by_id"`   // 12
	LastModifiedFromIP null.String `json:"last_modified_from_ip"` // 13
	AssociateId        null.Int    `json:"associate_id"`          // 14
	CustomerId         null.Int    `json:"customer_id"`           //15
	PartnerId          null.Int    `json:"partner_id"`            // 16
	StaffId            null.Int    `json:"staff_id"`              // 17
	WorkOrderId        null.Int    `json:"work_order_id"`         // 18
	State              int8        `json:"state"`                 // 19
	OldId              uint64      `json:"old_id"`                // 20
}

type PrivateFileRepository interface {
	Insert(ctx context.Context, u *PrivateFile) error
	UpdateById(ctx context.Context, u *PrivateFile) error
	GetById(ctx context.Context, id uint64) (*PrivateFile, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *PrivateFile) error
}
