package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

const (
	AssociateAwayLogActiveState   = 1
	AssociateAwayLogInactiveState = 0
)

type AssociateAwayLog struct {
	Id                   uint64      `json:"id"`
	Uuid                 string      `json:"uuid"`
	TenantId             uint64      `json:"tenant_id"`
	AssociateId          uint64      `json:"associate_id"`
	AssociateName        string      `json:"associate_name,omitempty"`
	AssociateLexicalName string      `json:"associate_lexical_name,omitempty"`
	Reason               int8        `json:"reason"`
	ReasonOther          null.String `json:"reason_other"`
	UntilFurtherNotice   bool        `json:"until_further_notice"`
	UntilDate            null.Time   `json:"until_date"`
	StartDate            null.Time   `json:"start_date"`
	State                int8        `json:"state"`
	CreatedTime          time.Time   `json:"created_time"`
	CreatedById          null.Int    `json:"created_by_id"`
	CreatedFromIP        string      `json:"created_from_ip"`
	LastModifiedTime     time.Time   `json:"last_modified_time"`
	LastModifiedById     null.Int    `json:"last_modified_by_id"`
	LastModifiedFromIP   string      `json:"last_modified_from_ip"`
	OldId                uint64      `json:"old_id"`
}

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `AssociateAwayLog` model.
type AssociateAwayLogFilter struct {
	TenantId   uint64 `json:"tenant_id"`
	States     []int8 `json:"states"`
	LastSeenId uint64 `json:"last_seen_id"`
	Limit      uint64 `json:"limit"`
}

type AssociateAwayLogRepository interface {
	Insert(ctx context.Context, u *AssociateAwayLog) error
	UpdateById(ctx context.Context, u *AssociateAwayLog) error
	GetById(ctx context.Context, id uint64) (*AssociateAwayLog, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*AssociateAwayLog, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *AssociateAwayLog) error
	ListByFilter(ctx context.Context, filter *AssociateAwayLogFilter) ([]*AssociateAwayLog, error)
	CountByFilter(ctx context.Context, filter *AssociateAwayLogFilter) (uint64, error)
}
