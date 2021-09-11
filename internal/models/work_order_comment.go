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

type WorkOrderComment struct {
	Id          uint64      `json:"id"`
	Uuid        string      `json:"uuid"`
	TenantId    uint64      `json:"tenant_id"`
	OrderId     uint64      `json:"order_id"`
	CommentId   uint64      `json:"comment_id"`
	CreatedTime time.Time   `json:"created_time"`
	Text        null.String `json:"text,omitempty"`
	OldId       uint64      `json:"old_id"`
}

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `WorkOrderComment` model.
type WorkOrderCommentFilter struct {
	TenantId    uint64    `json:"tenant_id"`
	States      []int8    `json:"states"`
	CreatedTime null.Time `json:"created_time"`
	LastSeenId  uint64    `json:"last_seen_id"`
	Limit       uint64    `json:"limit"`
}

type WorkOrderCommentRepository interface {
	Insert(ctx context.Context, u *WorkOrderComment) error
	UpdateById(ctx context.Context, u *WorkOrderComment) error
	GetById(ctx context.Context, id uint64) (*WorkOrderComment, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*WorkOrderComment, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrderComment) error
	ListByFilter(ctx context.Context, filter *WorkOrderCommentFilter) ([]*WorkOrderComment, error)
	CountByFilter(ctx context.Context, filter *WorkOrderCommentFilter) (uint64, error)
}
