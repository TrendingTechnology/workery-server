package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

const (
	BulletinBoardItemActiveState   = 1
	BulletinBoardItemArchivedState = 0
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `BulletinBoardItem` model.
type BulletinBoardItemFilter struct {
	TenantId   uint64 `json:"tenant_id"`
	States     []int8 `json:"states"`
	LastSeenId uint64 `json:"last_seen_id"`
	Limit      uint64 `json:"limit"`
}

type BulletinBoardItem struct {
	Id                 uint64      `json:"id"`
	Uuid               string      `json:"uuid"`
	TenantId           uint64      `json:"tenant_id"`
	Text               string      `json:"text"`
	CreatedTime        time.Time   `json:"created_time"`
	CreatedById        null.Int    `json:"created_by_id,omitempty"`
	CreatedByName      null.String `json:"created_by_name,omitempty"`
	CreatedFromIP      string      `json:"created_from_ip"`
	LastModifiedTime   time.Time   `json:"last_modified_time"`
	LastModifiedById   null.Int    `json:"last_modified_by_id,omitempty"`
	LastModifiedByName null.String `json:"last_modified_by_name,omitempty"`
	LastModifiedFromIP string      `json:"last_modified_from_ip"`
	State              int8        `json:"state"`
	OldId              uint64      `json:"old_id"`
}

type BulletinBoardItemRepository interface {
	Insert(ctx context.Context, u *BulletinBoardItem) error
	UpdateById(ctx context.Context, u *BulletinBoardItem) error
	GetById(ctx context.Context, id uint64) (*BulletinBoardItem, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*BulletinBoardItem, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *BulletinBoardItem) error
	ListByFilter(ctx context.Context, filter *BulletinBoardItemFilter) ([]*BulletinBoardItem, error)
	CountByFilter(ctx context.Context, filter *BulletinBoardItemFilter) (uint64, error)
}
