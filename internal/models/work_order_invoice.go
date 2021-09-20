package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type WorkOrderInvoice struct {
	Id       uint64 `json:"id"` // 	OrderId   uint64 `json:"order_id"`
	Uuid     string `json:"uuid"`
	TenantId uint64 `json:"tenant_id"`

	OrderId                  uint64      `json:"order_id"`
	InvoiceId                string      `json:"invoice_id"`
	InvoiceDate              time.Time   `json:"invoice_date"`
	AssociateName            string      `json:"associate_name"`
	AssociateTelephone       string      `json:"associate_telephone"`
	ClientName               string      `json:"client_name"`
	ClientTelephone          string      `json:"client_telephone"`
	ClientEmail              null.String `json:"client_email"`
	Line01Qty                int8        `json:"line_01_qty"`
	Line01Desc               string      `json:"line_01_desc"`
	Line01Price              float64     `json:"line_01_price"`
	Line01Amount             float64     `json:"line_01_amount"`
	Line02Qty                null.Int    `json:"line_02_qty"` // Make `int8`
	Line02Desc               null.String `json:"line_02_desc"`
	Line02Price              null.Float  `json:"line_02_price"`
	Line02Amount             null.Float  `json:"line_02_amount"`
	Line03Qty                null.Int    `json:"line_03_qty"` // Make `int8`
	Line03Desc               null.String `json:"line_03_desc"`
	Line03Price              null.Float  `json:"line_03_price"`
	Line03Amount             null.Float  `json:"line_03_amount"`
	Line04Qty                null.Int    `json:"line_04_qty"` // Make `int8`
	Line04Desc               null.String `json:"line_04_desc"`
	Line04Price              null.Float  `json:"line_04_price"`
	Line04Amount             null.Float  `json:"line_04_amount"`
	Line05Qty                null.Int    `json:"line_05_qty"` // Make `int8`
	Line05Desc               null.String `json:"line_05_desc"`
	Line05Price              null.Float  `json:"line_05_price"`
	Line05Amount             null.Float  `json:"line_05_amount"`
	Line06Qty                null.Int    `json:"line_06_qty"` // Make `int8`
	Line06Desc               null.String `json:"line_06_desc"`
	Line06Price              null.Float  `json:"line_06_price"`
	Line06Amount             null.Float  `json:"line_06_amount"`
	Line07Qty                null.Int    `json:"line_07_qty"` // Make `int8`
	Line07Desc               null.String `json:"line_07_desc"`
	Line07Price              null.Float  `json:"line_07_price"`
	Line07Amount             null.Float  `json:"line_07_amount"`
	Line08Qty                null.Int    `json:"line_08_qty"` // Make `int8`
	Line08Desc               null.String `json:"line_08_desc"`
	Line08Price              null.Float  `json:"line_08_price"`
	Line08Amount             null.Float  `json:"line_08_amount"`
	Line09Qty                null.Int    `json:"line_09_qty"` // Make `int8`
	Line09Desc               null.String `json:"line_09_desc"`
	Line09Price              null.Float  `json:"line_09_price"`
	Line09Amount             null.Float  `json:"line_09_amount"`
	Line10Qty                null.Int    `json:"line_10_qty"` // Make `int8`
	Line10Desc               null.String `json:"line_10_desc"`
	Line10Price              null.Float  `json:"line_10_price"`
	Line10Amount             null.Float  `json:"line_10_amount"`
	Line11Qty                null.Int    `json:"line_11_qty"` // Make `int8`
	Line11Desc               null.String `json:"line_11_desc"`
	Line11Price              null.Float  `json:"line_11_price"`
	Line11Amount             null.Float  `json:"line_11_amount"`
	Line12Qty                null.Int    `json:"line_12_qty"` // Make `int8`
	Line12Desc               null.String `json:"line_12_desc"`
	Line12Price              null.Float  `json:"line_12_price"`
	Line12Amount             null.Float  `json:"line_12_amount"`
	Line13Qty                null.Int    `json:"line_13_qty"` // Make `int8`
	Line13Desc               null.String `json:"line_13_desc"`
	Line13Price              null.Float  `json:"line_13_price"`
	Line13Amount             null.Float  `json:"line_13_amount"`
	Line14Qty                null.Int    `json:"line_14_qty"` // Make `int8`
	Line14Desc               null.String `json:"line_14_desc"`
	Line14Price              null.Float  `json:"line_14_price"`
	Line14Amount             null.Float  `json:"line_14_amount"`
	Line15Qty                null.Int    `json:"line_15_qty"` // Make `int8`
	Line15Desc               null.String `json:"line_15_desc"`
	Line15Price              null.Float  `json:"line_15_price"`
	Line15Amount             null.Float  `json:"line_15_amount"`
	InvoiceQuoteDays         int8        `json:"invoice_quote_days"`
	InvoiceAssociateTax      null.String `json:"invoice_associate_tax"`
	InvoiceQuoteDate         time.Time   `json:"invoice_quote_date"`
	InvoiceCustomersApproval string      `json:"invoice_customers_approval"`
	Line01Notes              null.String `json:"line_01_notes"`
	Line02Notes              null.String `json:"line_02_notes"`
	TotalLabour              float64     `json:"total_labour"`
	TotalMaterials           float64     `json:"total_materials"`
	OtherCosts               float64     `json:"other_costs"`
	Tax                      float64     `json:"tax"`
	Total                    float64     `json:"total"`
	PaymentAmount            float64     `json:"payment_amount"`
	PaymentDate              time.Time   `json:"payment_date"`
	IsCash                   bool        `json:"is_cash"`
	IsCheque                 bool        `json:"is_cheque"`
	IsDebit                  bool        `json:"is_debit"`
	IsCredit                 bool        `json:"is_credit"`
	IsOther                  bool        `json:"is_other"`
	ClientSignature          string      `json:"client_signature"`
	AssociateSignDate        time.Time   `json:"associate_sign_date"`
	AssociateSignature       string      `json:"associate_signature"`
	WorkOrderId              uint64      `json:"work_order_id"`
	CreatedTime              time.Time   `json:"created_time"`
	LastModifiedTime         time.Time   `json:"last_modified_time"`
	CreatedById              uint64      `json:"created_by_id"`
	CreatedByName            null.String `json:"created_by_name"`
	LastModifiedById         uint64      `json:"last_modified_by_id"`
	LastModifiedByName       null.String `json:"last_modified_by_name"`
	CreatedFrom              string      `json:"created_from"`
	CreatedFromIsPublic      bool        `json:"created_from_is_public"`
	LastModifiedFrom         string      `json:"last_modified_from"`
	LastModifiedFromIsPublic bool        `json:"last_modified_from_is_public"`
	ClientAddress            string      `json:"client_address"`
	RevisionVersion          int8        `json:"revision_version"`
	Deposit                  float64     `json:"deposit"`
	AmountDue                float64     `json:"amount_due"`
	SubTotal                 float64     `json:"sub_total"`

	State int8   `json:"state"` // IsArchived bool `json:"is_archived"`
	OldId uint64 `json:"old_id"`
}

type WorkOrderInvoiceRepository interface {
	Insert(ctx context.Context, u *WorkOrderInvoice) error
	UpdateById(ctx context.Context, u *WorkOrderInvoice) error
	GetById(ctx context.Context, id uint64) (*WorkOrderInvoice, error)
	GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*WorkOrderInvoice, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *WorkOrderInvoice) error
}
