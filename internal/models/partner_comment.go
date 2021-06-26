package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type PartnerComment struct {
	Id        uint64 `json:"id"`
	Uuid      string `json:"uuid"`
	TenantId  uint64 `json:"tenant_id"`
	PartnerId uint64 `json:"partner_id"`
	CommentId uint64 `json:"comment_id"`
	OldId     uint64 `json:"old_id"`
}

type PartnerCommentRepository interface {
	Insert(ctx context.Context, u *PartnerComment) error
	UpdateById(ctx context.Context, u *PartnerComment) error
	GetById(ctx context.Context, id uint64) (*PartnerComment, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*PartnerComment, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *PartnerComment) error
}
