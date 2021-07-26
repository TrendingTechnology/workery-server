package idos

import (
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
)

type LiteTenantFilterIDO struct {
	State      null.Int `json:"state"`
	LastSeenId uint64   `json:"last_seen_id"`
	Limit      uint64   `json:"limit"`
}

type LiteTenantIDO struct {
	Id         uint64 `db:"id"`
	SchemaName string `db:"schema_name"`
	Name       string `db:"name"`
	State      int8   `db:"state" json:"state"`
}

type LiteTenantListResponseIDO struct {
	NextPageToken uint64           `json:"next_page_token,omitempty"`
	TotalSize     uint64           `json:"total_size"`
	Results       []*LiteTenantIDO `json:"results"`
}

func NewLiteTenantListResponseIDO(arr []*models.LiteTenant, count uint64) *LiteTenantListResponseIDO {
	// Calculate next id.
	var nextPageToken uint64
	if len(arr) > 0 {
		lastRecord := arr[len(arr)-1]
		nextPageToken = lastRecord.Id
	}

	res := &LiteTenantListResponseIDO{ // Return through HTTP.
		TotalSize:     count,
		Results:       toLiteTenantIDOArray(arr),
		NextPageToken: nextPageToken,
	}

	return res
}

// DEPRECATED
func NewLiteTenantFilter(ido *LiteTenantFilterIDO) *models.LiteTenantFilter {
	return &models.LiteTenantFilter{
		State:      ido.State,
		LastSeenId: ido.LastSeenId,
		Limit:      ido.Limit,
	}
}

func toLiteTenantIDOArray(arr []*models.LiteTenant) []*LiteTenantIDO {
	var s []*LiteTenantIDO
	for _, v := range arr {
		s = append(s, toLiteTenantIDO(v))
	}
	return s
}

func toLiteTenantIDO(m *models.LiteTenant) *LiteTenantIDO {
	return &LiteTenantIDO{
		Id:         m.Id,
		SchemaName: m.SchemaName,
		Name:       m.Name,
		State:      m.State,
	}
}
