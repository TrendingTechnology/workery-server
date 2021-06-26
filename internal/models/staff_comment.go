package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type StaffComment struct {
	Id        uint64 `json:"id"`
	Uuid      string `json:"uuid"`
	TenantId  uint64 `json:"tenant_id"`
	StaffId   uint64 `json:"staff_id"`
	CommentId uint64 `json:"comment_id"`
	OldId     uint64 `json:"old_id"`
}

type StaffCommentRepository interface {
	Insert(ctx context.Context, u *StaffComment) error
	UpdateById(ctx context.Context, u *StaffComment) error
	GetById(ctx context.Context, id uint64) (*StaffComment, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *StaffComment) error
}
