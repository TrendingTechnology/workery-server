package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type WorkOrderInvoiceRepo struct {
	db *sql.DB
}

func NewWorkOrderInvoiceRepo(db *sql.DB) *WorkOrderInvoiceRepo {
	return &WorkOrderInvoiceRepo{
		db: db,
	}
}

func (r *WorkOrderInvoiceRepo) Insert(ctx context.Context, m *models.WorkOrderInvoice) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO work_order_invoices (
        uuid, tenant_id, old_id
    ) VALUES (
        $1, $2, $3
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.OldId,
	)
	return err
}

func (r *WorkOrderInvoiceRepo) UpdateById(ctx context.Context, m *models.WorkOrderInvoice) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        work_order_invoices
    SET
        tenant_id = $1
    WHERE
        id = $2`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.Id,
	)
	return err
}

func (r *WorkOrderInvoiceRepo) GetById(ctx context.Context, id uint64) (*models.WorkOrderInvoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderInvoice)

	query := `
    SELECT
        id, uuid, tenant_id
	FROM
        work_order_invoices
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId,
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

func (r *WorkOrderInvoiceRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.WorkOrderInvoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderInvoice)

	query := `
    SELECT
        id, uuid, tenant_id
	FROM
        work_order_invoices
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId,
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

func (r *WorkOrderInvoiceRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        work_order_invoices
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

func (r *WorkOrderInvoiceRepo) InsertOrUpdateById(ctx context.Context, m *models.WorkOrderInvoice) error {
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
