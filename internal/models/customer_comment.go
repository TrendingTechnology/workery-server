package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type CustomerComment struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`
	CustomerId    uint64 `json:"customer_id"`
	CommentId        uint64 `json:"comment_id"`
	OldId       uint64  `json:"old_id"`
}

type CustomerCommentRepository interface {
	Insert(ctx context.Context, u *CustomerComment) error
	UpdateById(ctx context.Context, u *CustomerComment) error
	GetById(ctx context.Context, id uint64) (*CustomerComment, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *CustomerComment) error
}
