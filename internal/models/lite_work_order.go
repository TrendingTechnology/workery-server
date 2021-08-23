package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `LiteWorkOrder` model.
type LiteWorkOrderFilter struct {
	TenantId         uint64   `json:"tenant_id"`
	States           []int8   `json:"states"`
	LastModifiedById null.Int `json:"last_modified_by_id"`
	SortOrder        string      `json:"sort_order"`
	SortField        string      `json:"sort_field"`
	Search           null.String `json:"search"`
	Offset           uint64      `json:"offset"`
	Limit            uint64   `json:"limit"`
}

type LiteWorkOrder struct {
	Id uint64 `json:"id"`
	// Uuid                              string      `json:"uuid"`
	TenantId    uint64   `json:"tenant_id"`
	CustomerId  uint64   `json:"customer_id"`
	AssociateId null.Int `json:"associate_id"`
	// Description                       string      `json:"description"`
	AssignmentDate null.Time   `json:"assignment_date"`
	IsOngoing                         bool        `json:"is_ongoing"`
	// IsHomeSupportService              bool        `json:"is_home_support_service"`
	StartDate                         time.Time   `json:"start_date"`
	// CompletionDate                    null.Time   `json:"completion_date"`
	// Hours                             float64     `json:"hours"`
	TypeOf                            int8        `json:"type_of"`
	// IndexedText                       string      `json:"indexed_text"`
	// ClosingReason                     int8        `json:"closing_reason"`
	// ClosingReasonOther                null.String `json:"closing_reason_other"`
	State int8 `json:"state"`
	// Currency                          string      `json:"currency"`
	// WasJobSatisfactory                bool        `json:"was_job_satisfactory"`
	// WasJobFinishedOnTimeAndOnBudget   bool        `json:"was_job_finished_on_time_and_on_budget"`
	// WasAssociatePunctual              bool        `json:"was_associate_punctual"`
	// WasAssociateProfessional          bool        `json:"was_associate_professional"`
	// WouldCustomerReferOurOrganization bool        `json:"would_customer_refer_our_organization"`
	// Score                             int8        `json:"score"`
	// InvoiceDate                       null.Time   `json:"invoice_date"`
	// InvoiceQuoteAmount                float64     `json:"invoice_quote_amount"`
	// InvoiceLabourAmount               float64     `json:"invoice_labour_amount"`
	// InvoiceMaterialAmount             float64     `json:"invoice_material_amount"`
	// InvoiceTaxAmount                  float64     `json:"invoice_tax_amount"`
	// InvoiceTotalAmount                float64     `json:"invoice_total_amount"`
	// InvoiceServiceFeeAmount           float64     `json:"invoice_service_fee_amount"`
	// InvoiceServiceFeePaymentDate      null.Time   `json:"invoice_service_fee_payment_date"`
	// CreatedTime                       time.Time   `json:"created_time"`
	// CreatedById                       null.Int    `json:"created_by_id"`
	// CreatedFromIP                     null.String `json:"created_from_ip"`
	// LastModifiedTime                  time.Time   `json:"last_modified_time"`
	// LastModifiedById                  null.Int    `json:"last_modified_by_id"`
	// LastModifiedFromIP                null.String `json:"last_modified_from_ip"`
	// OldId                             uint64      `json:"old_id"`
	// InvoiceServiceFeeId               null.Int    `json:"invoice_service_fee_id"`
	// LatestPendingTaskId               null.Int    `json:"latest_pending_task_id"`
	// OngoingWorkOrderId                null.Int    `json:"ongoing_work_order_id"`
	// WasSurveyConducted                bool        `json:"was_survey_conducted"`
	// WasThereFinancialsInputted        bool        `json:"was_there_financials_inputted"`
	// InvoiceActualServiceFeeAmountPaid float64     `json:"invoice_actual_service_fee_amount_paid"`
	// InvoiceBalanceOwingAmount         float64     `json:"invoice_balance_owing_amount"`
	// InvoiceQuotedLabourAmount         float64     `json:"invoice_quoted_labour_amount"`
	// InvoiceQuotedMaterialAmount       float64     `json:"invoice_quoted_material_amount"`
	// InvoiceTotalQuoteAmount           float64     `json:"invoice_total_quote_amount"`
	// Visits                            int8        `json:"visits"`
	// InvoiceIds                        null.String `json:"invoice_ids"`
	// NoSurveyConductedReason           null.Int    `json:"no_survey_conducted_reason"`
	// NoSurveyConductedReasonOther      null.String `json:"no_survey_conducted_reason_other"`
	// ClonedFromId                      null.Int    `json:"cloned_from_id"`
	// InvoiceDepositAmount              float64     `json:"invoice_deposit_amount"`
	// InvoiceOtherCostsAmount           float64     `json:"invoice_other_costs_amount"`
	// InvoiceQuotedOtherCostsAmount     float64     `json:"invoice_quoted_other_costs_amount"`
	// InvoicePaidTo                     null.Int    `json:"invoice_paid_to"`
	// InvoiceAmountDue                  float64     `json:"invoice_amount_due"`
	// InvoiceSubTotalAmount             float64     `json:"invoice_sub_total_amount"`
	// ClosingReasonComment              string      `json:"closing_reason_comment"`
}

type LiteWorkOrderRepository interface {
	ListByFilter(ctx context.Context, filter *LiteWorkOrderFilter) ([]*LiteWorkOrder, error)
	CountByFilter(ctx context.Context, filter *LiteWorkOrderFilter) (uint64, error)
}
