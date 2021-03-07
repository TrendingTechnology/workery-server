package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type HowHearAboutUsItem struct {
	Id                uint64    `json:"id"`
	Uuid              string    `json:"uuid"`
	TenantId          uint64    `json:"tenant_id"`
	Text              string    `json:"text"`
	SortNumber        int8      `json:"sort_number"`
	IsForAssociate    bool      `json:"is_for_associate"`
	IsForCustomer     bool      `json:"is_for_customer"`
	IsForStaff        bool      `json:"is_for_staff"`
	IsForPartner      bool      `json:"is_for_partner"`
	State             int8      `json:"state"`
    OldId             uint64    `json:"old_id"`
}

type HowHearAboutUsItemRepository interface {
	Insert(ctx context.Context, u *HowHearAboutUsItem) error
	UpdateById(ctx context.Context, u *HowHearAboutUsItem) error
	GetById(ctx context.Context, id uint64) (*HowHearAboutUsItem, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*HowHearAboutUsItem, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *HowHearAboutUsItem) error
}
