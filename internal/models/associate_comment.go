package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type AssociateComment struct {
	Id                     uint64 `json:"id"`
	Uuid                   string `json:"uuid"`
	TenantId               uint64 `json:"tenant_id"`
	AssociateId            uint64 `json:"associate_id"`
	CommentId              uint64 `json:"comment_id"`
	OldId                  uint64 `json:"old_id"`
}

type AssociateCommentRepository interface {
	Insert(ctx context.Context, u *AssociateComment) error
	UpdateById(ctx context.Context, u *AssociateComment) error
	GetById(ctx context.Context, id uint64) (*AssociateComment, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*AssociateComment, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *AssociateComment) error
}
