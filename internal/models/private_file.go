package models

import (
	"context"
	"database/sql"
	"time"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type PrivateFile struct {
	Id                 uint64        `json:"id"`
	Uuid               string        `json:"uuid"`
	TenantId           uint64        `json:"tenant_id"`
	ImageFile          string        `json:"image_file"`
	CreatedTime        time.Time     `json:"created_time"`
	CreatedFromIP      string        `json:"created_from_ip"`
	LastModifiedTime   time.Time     `json:"last_modified_time"`
	LastModifiedFromIP string        `json:"last_modified_from_ip"`
	CreatedById        sql.NullInt64 `json:"created_by_id,omitempty"`
	LastModifiedById   sql.NullInt64 `json:"last_modified_by_id,omitempty"`
	State              int8          `json:"state"`
	OldId              uint64        `json:"old_id"`
}

type PrivateFileRepository interface {
	Insert(ctx context.Context, u *PrivateFile) error
	UpdateById(ctx context.Context, u *PrivateFile) error
	GetById(ctx context.Context, id uint64) (*PrivateFile, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *PrivateFile) error
}
