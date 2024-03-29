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

type Comment struct {
	Id                 uint64      `json:"id"`
	Uuid               string      `json:"uuid"`
	TenantId           uint64      `json:"tenant_id"`
	CreatedTime        time.Time   `json:"created_time"`
	CreatedById        null.Int    `json:"created_by_id,omitempty"`
	CreatedByName      null.String `json:"created_by_name,omitempty"`
	CreatedFromIP      string      `json:"created_from_ip"`
	LastModifiedTime   time.Time   `json:"last_modified_time"`
	LastModifiedById   null.Int    `json:"last_modified_by_id,omitempty"`
	LastModifiedByName null.String `json:"last_modified_by_name,omitempty"`
	LastModifiedFromIP string      `json:"last_modified_from_ip,omitempty"`
	Text               string      `json:"text"`
	State              int8        `json:"state"`
	OldId              uint64      `json:"old_id"`
}

type CommentRepository interface {
	Insert(ctx context.Context, u *Comment) error
	UpdateById(ctx context.Context, u *Comment) error
	GetById(ctx context.Context, id uint64) (*Comment, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *Comment) error
}
