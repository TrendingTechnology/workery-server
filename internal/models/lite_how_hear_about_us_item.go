package models

import (
	"context"
	// "time"
	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteHowHearAboutUsItem` model.
type LiteHowHearAboutUsItemFilter struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type LiteHowHearAboutUsItem struct {
	Id        uint64    `json:"id"`
	TenantId  uint64    `json:"tenant_id"`
	Text           string `json:"text"`
	SortNumber     int8   `json:"sort_number"`
	IsForAssociate bool   `json:"is_for_associate"`
	IsForCustomer  bool   `json:"is_for_customer"`
	IsForStaff     bool   `json:"is_for_staff"`
	IsForPartner   bool   `json:"is_for_partner"`
	State          int8   `json:"state"`
}

type LiteHowHearAboutUsItemRepository interface {
	ListByFilter(ctx context.Context, filter *LiteHowHearAboutUsItemFilter) ([]*LiteHowHearAboutUsItem, error)
	CountByFilter(ctx context.Context, filter *LiteHowHearAboutUsItemFilter) (uint64, error)
}
