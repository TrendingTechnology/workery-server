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
	Id                 uint64      `json:"id"`
	Uuid               string      `json:"uuid"`
	TenantId           uint64      `json:"tenant_id"`
	AssociateId        uint64      `json:"associate_id"`
	Reason             int8        `json:"reason"`
	ReasonOther        null.String `json:"reason_other"`
	UntilFurtherNotice bool        `json:"until_further_notice"`
	UntilDate          null.Time   `json:"until_date"`
	StartDate          null.Time   `json:"start_date"`
	State              int8        `json:"state"`
	CreatedTime        time.Time   `json:"created_time"`
	CreatedById        uint64      `json:"created_by_id"`
	CreatedFromIP      string      `json:"created_from_ip"`
	LastModifiedTime   time.Time   `json:"last_modified_time"`
	LastModifiedById   uint64      `json:"last_modified_by_id"`
	LastModifiedFromIP string      `json:"last_modified_from_ip"`
	OldId              uint64      `json:"old_id"`
}

type AssociateAwayLogRepository interface {
	Insert(ctx context.Context, u *AssociateAwayLog) error
	UpdateById(ctx context.Context, u *AssociateAwayLog) error
	GetById(ctx context.Context, id uint64) (*AssociateAwayLog, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*AssociateAwayLog, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *AssociateAwayLog) error
}
