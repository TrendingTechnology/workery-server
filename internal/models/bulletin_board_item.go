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

type BulletinBoardItem struct {
	Id                 uint64        `json:"id"`
	Uuid               string        `json:"uuid"`
	TenantId           uint64        `json:"tenant_id"`
	Text               string        `json:"text"`
	CreatedTime        time.Time     `json:"created_time"`
	CreatedById        sql.NullInt64 `json:"created_by_id,omitempty"`
	CreatedFromIP      string        `json:"created_from_ip"`
	LastModifiedTime   time.Time     `json:"last_modified_time"`
	LastModifiedById   sql.NullInt64 `json:"last_modified_by_id,omitempty"`
	LastModifiedFromIP string        `json:"last_modified_from_ip"`
	State              int8          `json:"state"`
	OldId              uint64        `json:"old_id"`
}

type BulletinBoardItemRepository interface {
	Insert(ctx context.Context, u *BulletinBoardItem) error
	UpdateById(ctx context.Context, u *BulletinBoardItem) error
	GetById(ctx context.Context, id uint64) (*BulletinBoardItem, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*BulletinBoardItem, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *BulletinBoardItem) error
}
