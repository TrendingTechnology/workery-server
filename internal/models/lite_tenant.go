package models

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

type LiteTenantFilter struct {
	State      null.Int `json:"state"`
	LastSeenId uint64   `json:"last_seen_id"`
	Limit      uint64   `json:"limit"`
}

type LiteTenant struct {
	Id         uint64 `db:"id"`
	SchemaName string `db:"schema_name"`
	Name       string `db:"name"`
	State      int8   `db:"state" json:"state"`
}

type LiteTenantRepository interface {
	ListAllIds(ctx context.Context) ([]uint64, error)
	ListByFilter(ctx context.Context, f *LiteTenantFilter) ([]*LiteTenant, error)
	CountByFilter(ctx context.Context, f *LiteTenantFilter) (uint64, error)
}
