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

type PublicImageUpload struct {
	Id                 uint64      `json:"id"`
	Uuid               string      `json:"uuid"`
	TenantId           uint64      `json:"tenant_id"`
	ImageFile          string      `json:"image_file"`
	CreatedTime        time.Time   `json:"created_time"`
	CreatedFromIP      string      `json:"created_from_ip"`
	LastModifiedTime   time.Time   `json:"last_modified_time"`
	LastModifiedFromIP string      `json:"last_modified_from_ip"`
	CreatedById        null.Int    `json:"created_by_id,omitempty"`
	CreatedByString    null.String `json:"created_by_name,omitempty"`
	LastModifiedById   null.Int    `json:"last_modified_by_id,omitempty"`
	LastModifiedByName null.String `json:"last_modified_by_name,omitempty"`
	State              int8        `json:"state"`
	OldId              uint64      `json:"old_id"`
}

type PublicImageUploadRepository interface {
	Insert(ctx context.Context, u *PublicImageUpload) error
	UpdateById(ctx context.Context, u *PublicImageUpload) error
	GetById(ctx context.Context, id uint64) (*PublicImageUpload, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *PublicImageUpload) error
}
