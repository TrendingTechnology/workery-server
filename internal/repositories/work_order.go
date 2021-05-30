package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type WorkOrderRepo struct {
	db *sql.DB
}

func NewWorkOrderRepo(db *sql.DB) *WorkOrderRepo {
	return &WorkOrderRepo{
		db: db,
	}
}

func (r *WorkOrderRepo) Insert(ctx context.Context, m *models.WorkOrder) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO work_orders (
        uuid, tenant_id, customer_id, associate_id, description, assignment_date,
        is_ongoing, is_home_support_service, start_date, completion_date, hours,
		indexed_text, closing_reason, closing_reason_other, state, currency,
		was_job_satisfactory, was_job_finished_on_time_and_on_budget, was_associate_punctual,
		was_associate_professional, would_customer_refer_our_organization, score,
        invoice_date, invoice_quote_amount, invoice_labour_amount, invoice_material_amount,
		invoice_tax_amount, invoice_total_amount, invoice_service_fee_amount, invoice_service_fee_payment_date,
        created_time, created_by_id, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_from_ip, old_id,
		invoice_service_fee_id, latest_pending_task_id, ongoing_work_order_id,
		was_survey_conducted, was_there_financials_inputted, invoice_actual_service_fee_amount_paid,
		invoice_balance_owing_amount, invoice_quoted_labour_amount, invoice_quoted_material_amount,
		invoice_total_quote_amount, visits, invoice_ids, no_survey_conducted_reason,
		no_survey_conducted_reason_other, cloned_from_id, invoice_deposit_amount,
		invoice_other_costs_amount, invoice_quoted_other_costs_amount, invoice_paid_to,
		invoice_amount_due, invoice_sub_total_amount, closing_reason_comment
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
		$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
		$33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47,
		$48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.CustomerId, m.AssociateId, m.Description,
		m.AssignmentDate, m.IsOngoing, m.IsHomeSupportService, m.StartDate, m.CompletionDate, m.Hours,
		m.IndexedText, m.ClosingReason, m.ClosingReasonOther, m.State, m.Currency,
		m.WasJobSatisfactory, m.WasJobFinishedOnTimeAndOnBudget, m.WasAssociatePunctual,
		m.WasAssociateProfessional, m.WouldCustomerReferOurOrganization, m.Score,
		m.InvoiceDate, m.InvoiceQuoteAmount, m.InvoiceLabourAmount, m.InvoiceMaterialAmount,
		m.InvoiceTaxAmount, m.InvoiceTotalAmount, m.InvoiceServiceFeeAmount, m.InvoiceServiceFeePaymentDate,
		m.CreatedTime, m.CreatedById, m.CreatedFromIP,
		m.LastModifiedTime, m.LastModifiedById, m.LastModifiedFromIP, m.OldId,
		m.InvoiceServiceFeeId, m.LatestPendingTaskId, m.OngoingWorkOrderId,
		m.WasSurveyConducted, m.WasThereFinancialsInputted, m.InvoiceActualServiceFeeAmountPaid,
		m.InvoiceBalanceOwingAmount, m.InvoiceQuotedLabourAmount, m.InvoiceQuotedMaterialAmount,
		m.InvoiceTotalQuoteAmount, m.Visits, m.InvoiceIds, m.NoSurveyConductedReason,
		m.NoSurveyConductedReasonOther, m.ClonedFromId, m.InvoiceDepositAmount,
		m.InvoiceOtherCostsAmount, m.InvoiceQuotedOtherCostsAmount, m.InvoicePaidTo,
		m.InvoiceAmountDue, m.InvoiceSubTotalAmount, m.ClosingReasonComment,
	)
	return err
}

func (r *WorkOrderRepo) UpdateById(ctx context.Context, m *models.WorkOrder) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        work_orders
    SET
        tenant_id = $1, customer_id = $2, associate_id = $3, state = $4,
		last_modified_time = $5, last_modified_by_id = $6, last_modified_from_ip = $7
    WHERE
        id = $8`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.CustomerId, m.AssociateId, m.State, m.LastModifiedTime,
		m.LastModifiedById, m.LastModifiedFromIP, m.Id,
	)
	return err
}

func (r *WorkOrderRepo) GetById(ctx context.Context, id uint64) (*models.WorkOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrder)

	query := `
    SELECT
        id, uuid, tenant_id, customer_id, associate_id, state,
		created_time, created_by_id, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_from_ip
	FROM
        work_orders
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.CustomerId, &m.AssociateId,
		&m.State, &m.CreatedTime, &m.CreatedById, &m.CreatedFromIP, &m.LastModifiedTime,
		&m.LastModifiedById, &m.LastModifiedFromIP,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *WorkOrderRepo) GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
    SELECT
        id
    FROM
        work_orders
    WHERE
		tenant_id = $1
	AND
	    old_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, tid, oid).Scan(&newId)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return 0, nil
		} else { // CASE 2 OF 2: All other errors.
			return 0, err
		}
	}
	return newId, nil
}

func (r *WorkOrderRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        work_orders
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *WorkOrderRepo) InsertOrUpdateById(ctx context.Context, m *models.WorkOrder) error {
	if m.Id == 0 {
		return r.Insert(ctx, m)
	}

	doesExist, err := r.CheckIfExistsById(ctx, m.Id)
	if err != nil {
		return err
	}

	if doesExist == false {
		return r.Insert(ctx, m)
	}
	return r.UpdateById(ctx, m)
}
