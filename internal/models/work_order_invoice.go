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
	Id          uint64 `json:"id"` // 	OrderId   uint64 `json:"order_id"`
	Uuid        string `json:"uuid"`
	TenantId    uint64 `json:"tenant_id"`

	InvoiceId string `json:"invoice_id"`
	InvoiceDate time.Time `json:"invoice_date"`
	AssociateName string `json:"associate_name"`
	AssociateTelephone string `json:"associate_telephone"`
	ClientName string `json:"client_name"`
	ClientTelephone string `json:"client_telephone"`
	ClientEmail null.String `json:"client_email"`
	Line01Qty int8 `json:"line_01_qty"`
	Line01Desc string `json:"line_01_desc"`
	Line01PriceCurrency string `json:"line_01_price_currency"`
	Line01Price float64 `json:"line_01_price"`
	Line01AmountCurrency string `json:"line_01_amount_currency"`
	Line01Amount float64 `json:"line_01_amount"`
	Line02Qty null.Int `json:"line_02_qty"` // Make `int8`
	Line02Desc string `json:"line_02_desc"`
	Line02PriceCurrency string `json:"line_02_price_currency"`
	Line02Price float64 `json:"line_02_price"`
	Line02AmountCurrency string `json:"line_02_amount_currency"`
	Line02Amount float64 `json:"line_02_amount"`
	Line03Qty null.Int `json:"line_03_qty"` // Make `int8`
	Line03Desc string `json:"line_03_desc"`
	Line03PriceCurrency string `json:"line_03_price_currency"`
	Line03Price float64 `json:"line_03_price"`
	Line03AmountCurrency string `json:"line_03_amount_currency"`
	Line03Amount float64 `json:"line_03_amount"`
	Line04Qty null.Int `json:"line_04_qty"` // Make `int8`
	Line04Desc string `json:"line_04_desc"`
	Line04PriceCurrency string `json:"line_04_price_currency"`
	Line04Price float64 `json:"line_04_price"`
	Line04AmountCurrency string `json:"line_04_amount_currency"`
	Line04Amount float64 `json:"line_04_amount"`
	Line05Qty null.Int `json:"line_05_qty"` // Make `int8`
	Line05Desc string `json:"line_05_desc"`
	Line05PriceCurrency string `json:"line_05_price_currency"`
	Line05Price float64 `json:"line_05_price"`
	Line05AmountCurrency string `json:"line_05_amount_currency"`
	Line05Amount float64 `json:"line_05_amount"`
	Line06Qty null.Int `json:"line_06_qty"` // Make `int8`
	Line06Desc string `json:"line_06_desc"`
	Line06PriceCurrency string `json:"line_06_price_currency"`
	Line06Price float64 `json:"line_06_price"`
	Line06AmountCurrency string `json:"line_06_amount_currency"`
	Line06Amount float64 `json:"line_06_amount"`
	Line07Qty null.Int `json:"line_07_qty"` // Make `int8`
	Line07Desc string `json:"line_07_desc"`
	Line07PriceCurrency string `json:"line_07_price_currency"`
	Line07Price float64 `json:"line_07_price"`
	Line07AmountCurrency string `json:"line_07_amount_currency"`
	Line07Amount float64 `json:"line_07_amount"`
	Line08Qty null.Int `json:"line_08_qty"` // Make `int8`
	Line08Desc string `json:"line_08_desc"`
	Line08PriceCurrency string `json:"line_08_price_currency"`
	Line08Price float64 `json:"line_08_price"`
	Line08AmountCurrency string `json:"line_08_amount_currency"`
	Line08Amount float64 `json:"line_08_amount"`
	Line09Qty null.Int `json:"line_09_qty"` // Make `int8`
	Line09Desc string `json:"line_09_desc"`
	Line09PriceCurrency string `json:"line_09_price_currency"`
	Line09Price float64 `json:"line_09_price"`
	Line09AmountCurrency string `json:"line_09_amount_currency"`
	Line09Amount float64 `json:"line_09_amount"`
	Line10Qty null.Int `json:"line_10_qty"` // Make `int8`
	Line10Desc string `json:"line_10_desc"`
	Line10PriceCurrency string `json:"line_10_price_currency"`
	Line10Price float64 `json:"line_10_price"`
	Line10AmountCurrency string `json:"line_10_amount_currency"`
	Line10Amount float64 `json:"line_10_amount"`
	Line11Qty null.Int `json:"line_11_qty"` // Make `int8`
	Line11Desc string `json:"line_11_desc"`
	Line11PriceCurrency string `json:"line_11_price_currency"`
	Line11Price float64 `json:"line_11_price"`
	Line11AmountCurrency string `json:"line_11_amount_currency"`
	Line11Amount float64 `json:"line_11_amount"`
	Line12Qty null.Int `json:"line_12_qty"` // Make `int8`
	Line12Desc string `json:"line_12_desc"`
	Line12PriceCurrency string `json:"line_12_price_currency"`
	Line12Price float64 `json:"line_12_price"`
	Line12AmountCurrency string `json:"line_12_amount_currency"`
	Line12Amount float64 `json:"line_12_amount"`
	Line13Qty null.Int `json:"line_13_qty"` // Make `int8`
	Line13Desc string `json:"line_13_desc"`
	Line13PriceCurrency string `json:"line_13_price_currency"`
	Line13Price float64 `json:"line_13_price"`
	Line13AmountCurrency string `json:"line_13_amount_currency"`
	Line13Amount float64 `json:"line_13_amount"`
	Line14Qty null.Int `json:"line_14_qty"` // Make `int8`
	Line14Desc string `json:"line_14_desc"`
	Line14PriceCurrency string `json:"line_14_price_currency"`
	Line14Price float64 `json:"line_14_price"`
	Line14AmountCurrency string `json:"line_14_amount_currency"`
	Line14Amount float64 `json:"line_14_amount"`
	Line15Qty null.Int `json:"line_15_qty"` // Make `int8`
	Line15Desc string `json:"line_15_desc"`
	Line15PriceCurrency string `json:"line_15_price_currency"`
	Line15Price float64 `json:"line_15_price"`
	Line15AmountCurrency string `json:"line_15_amount_currency"`
	Line15Amount float64 `json:"line_15_amount"`
	InvoiceQuoteDays int8 `json:"invoice_quote_days"`
	InvoiceAssociateTax null.String `json:"invoice_associate_tax"`
	InvoiceQuoteDate time.Time `json:"invoice_quote_date"`
	InvoiceCustomersApproval string `json:"invoice_customers_approval"`
	Line01Notes null.String `json:"line_01_notes"`
	Line02Notes null.String `json:"line_02_notes"`
	TotalLabourCurrency string `json:"total_labour_currency"`
	TotalLabour float64 `json:"total_labour"`
	TotalMaterialsCurrency string `json:"total_materials_currency"`
	TotalMaterials float64 `json:"total_materials"`
	OtherCostsCurrency string `json:"other_costs_currency"`
	OtherCosts float64 `json:"other_costs"`
	AmountDueCurrency string `json:"amount_due_currency"`
	TaxCurrency string `json:"tax_currency"`
	Tax float64 `json:"tax"`
	TotalCurrency string `json:"total_currency"`
	Total float64 `json:"total"`
	DepositCurrency string `json:"deposit_currency"`
	PaymentAmountCurrency string `json:"payment_amount_currency"`
	PaymentAmount float64 `json:"payment_amount"`
	PaymentDate time.Time `json:"payment_date"`
	IsCash bool `json:"is_cash"`
	IsCheque bool `json:"is_cheque"`
	IsDebit bool `json:"is_debit"`
	IsCredit bool `json:"is_credit"`
	IsOther bool `json:"is_other"`
	ClientSignature string `json:"client_signature"`
	AssociateSignDate time.Time `json:"associate_sign_date"`
	AssociateSignature string `json:"associate_signature"`
	WorkOrderId   uint64 `json:"work_order_id"`
	CreatedAt time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	CreatedById   uint64 `json:"created_by_id"`
	LastModifiedById   uint64 `json:"last_modified_by_id"`
	CreatedFrom string `json:"created_from"`
	CreatedFromIsPublic bool `json:"created_from_is_public"`
	LastModifiedFrom string `json:"last_modified_from"`
	LastModifiedFromIsPublic bool `json:"last_modified_from_is_public"`
	ClientAddress string `json:"client_address"`
	RevisionVersion int8 `json:"revision_version"`
	Deposit float64 `json:"deposit"`
	AmountDue float64 `json:"amount_due"`
	SubTotal float64 `json:"sub_total"`
	SubTotalCurrency string `json:"sub_total_currency"`

    State       int8 `json:"state"`	 // IsArchived bool `json:"is_archived"`
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
