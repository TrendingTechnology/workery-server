package models

import (
	"context"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type WorkOrderInvoice struct {
	Id          uint64 `json:"id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`
	OrderId     uint64 `json:"order_id"`
	InvoiceId       uint64 `json:"invoice_id"`
	OldId       uint64 `json:"old_id"`
}

type WorkOrderInvoiceRepository interface {
	Insert(ctx context.Context, u *WorkOrderInvoice) error
	UpdateById(ctx context.Context, u *WorkOrderInvoice) error
	GetById(ctx context.Context, id uint64) (*WorkOrderInvoice, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*WorkOrderInvoice, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrderInvoice) error
}
