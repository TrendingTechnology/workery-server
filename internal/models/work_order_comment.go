package models

import (
	"context"
	"time"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type WorkOrderComment struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`
	OrderId     uint64 `json:"order_id"`
	CommentId   uint64 `json:"comment_id"`
	CreatedTime time.Time   `json:"created_time"`
	OldId       uint64 `json:"old_id"`
}

type WorkOrderCommentRepository interface {
	Insert(ctx context.Context, u *WorkOrderComment) error
	UpdateById(ctx context.Context, u *WorkOrderComment) error
	GetById(ctx context.Context, id uint64) (*WorkOrderComment, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*WorkOrderComment, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrderComment) error
}
