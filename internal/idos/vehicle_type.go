package idos

import (
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
)

type LiteVehicleTypeFilterIDO struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder null.String `json:"sort_order"`
	SortField null.String `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"last_seen_id"`
	Limit     uint64      `json:"limit"`
}

type LiteVehicleTypeListResponseIDO struct {
	NextId  uint64                    `json:"next_id,omitempty"`
	Count   uint64                    `json:"count"`
	Results []*models.LiteVehicleType `json:"results"`
}

func NewLiteVehicleTypeListResponseIDO(arr []*models.LiteVehicleType, count uint64) *LiteVehicleTypeListResponseIDO {
	// Calculate next id.
	var nextId uint64
	if len(arr) > 0 {
		lastRecord := arr[len(arr)-1]
		nextId = lastRecord.Id
	}

	res := &LiteVehicleTypeListResponseIDO{ // Return through HTTP.
		Count:   count,
		Results: arr,
		NextId:  nextId,
	}

	return res
}
