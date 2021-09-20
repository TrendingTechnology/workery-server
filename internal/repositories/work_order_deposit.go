package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type WorkOrderDepositRepo struct {
	db *sql.DB
}

func NewWorkOrderDepositRepo(db *sql.DB) *WorkOrderDepositRepo {
	return &WorkOrderDepositRepo{
		db: db,
	}
}

func (r *WorkOrderDepositRepo) Insert(ctx context.Context, m *models.WorkOrderDeposit) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO work_order_deposits (
        uuid, tenant_id, paid_at, deposit_method, paid_to, currency,
		amount, paid_for, created_time, last_modified_time, created_by_id,
		created_by_name, last_modified_by_id, last_modified_by_name, order_id, created_from_ip,
		last_modified_from_ip, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		$17, $18, $19
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.PaidAt, m.DepositMethod, m.PaidTo, m.Currency,
		m.Amount, m.PaidFor, m.CreatedTime, m.LastModifiedTime, m.CreatedById,
		m.CreatedByName, m.LastModifiedById, m.LastModifiedByName, m.OrderId, m.CreatedFromIP,
		m.LastModifiedFromIP, m.State, m.OldId,
	)
	return err
}

func (r *WorkOrderDepositRepo) UpdateById(ctx context.Context, m *models.WorkOrderDeposit) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        work_order_deposits
    SET
        tenant_id = $1, text = $2, description = $3, state = $4
    WHERE
        id = $5`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId,
		// m.Text, m.Description,
		m.State, m.Id,
	)
	return err
}

func (r *WorkOrderDepositRepo) GetById(ctx context.Context, id uint64) (*models.WorkOrderDeposit, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderDeposit)

	query := `
    SELECT
        id, uuid, tenant_id, text, description, state
	FROM
        work_order_deposits
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId,
		// &m.Text, &m.Description,
		&m.State,
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

func (r *WorkOrderDepositRepo) GetIdByOldId(ctx context.Context, tenantId uint64, oldId uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
	SELECT
		id
	FROM
		work_order_deposits
	WHERE
		tenant_id = $1 AND old_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, tenantId, oldId).Scan(&newId)
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

func (r *WorkOrderDepositRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        work_order_deposits
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

func (r *WorkOrderDepositRepo) InsertOrUpdateById(ctx context.Context, m *models.WorkOrderDeposit) error {
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
