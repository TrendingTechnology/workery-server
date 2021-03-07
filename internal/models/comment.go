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

type Comment struct {
	Id                uint64    `json:"id"`
	Uuid              string    `json:"uuid"`
	TenantId          uint64    `json:"tenant_id"`
	CreatedTime       time.Time `json:"created_time"`
	CreatedById       sql.NullInt64 `json:"created_by_id,omitempty"`
	LastModifiedTime  time.Time `json:"last_modified_time"`
	LastModifiedById  sql.NullInt64 `json:"last_modified_by_id,omitempty"`
	Text              string    `json:"text"`
	State             int8      `json:"state"`
    OldId             uint64    `json:"old_id"`
}

type CommentRepository interface {
	Insert(ctx context.Context, u *Comment) error
	UpdateById(ctx context.Context, u *Comment) error
	GetById(ctx context.Context, id uint64) (*Comment, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*Comment, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *Comment) error
}
