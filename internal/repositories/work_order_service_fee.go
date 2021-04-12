package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type WorkOrderServiceFeeRepo struct {
	db *sql.DB
}

func NewWorkOrderServiceFeeRepo(db *sql.DB) *WorkOrderServiceFeeRepo {
	return &WorkOrderServiceFeeRepo{
		db: db,
	}
}

func (r *WorkOrderServiceFeeRepo) Insert(ctx context.Context, m *models.WorkOrderServiceFee) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO work_order_service_fees (
        uuid, tenant_id, title, description, percentage, created_time,
		created_by_id, created_from_ip, last_modified_time, last_modified_by_id,
		last_modified_from_ip, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.Title, m.Description, m.Percentage,
		m.CreatedTime, m.CreatedById, m.CreatedFromIP,
		m.LastModifiedTime, m.LastModifiedById, m.LastModifiedFromIP,
		m.State, m.OldId,
	)
	return err
}

func (r *WorkOrderServiceFeeRepo) UpdateById(ctx context.Context, m *models.WorkOrderServiceFee) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        work_order_service_fees
    SET
        tenant_id = $1, title = $2, description = $3, state = $4,
		percentage = $5, created_time = $6, created_by_id = $7,
		created_from_ip = $8, last_modified_time = $9, last_modified_by_id = $10,
		last_modified_from_ip = $11,
    WHERE
        id = $12`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.Title, m.Description, m.State, m.Percentage, m.CreatedTime, m.CreatedById, m.CreatedFromIP,
		m.LastModifiedTime, m.LastModifiedById, m.LastModifiedFromIP, m.Id,
	)
	return err
}

func (r *WorkOrderServiceFeeRepo) GetById(ctx context.Context, id uint64) (*models.WorkOrderServiceFee, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderServiceFee)

	query := `
    SELECT
        id, uuid, tenant_id, title, description, percentage, created_time,
		created_by_id, created_from_ip, last_modified_time, last_modified_by_id,
		last_modified_from_ip, state
	FROM
        work_order_service_fees
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Title, &m.Description, &m.State,
		&m.Percentage, &m.CreatedTime, &m.CreatedById, &m.CreatedFromIP,
		&m.LastModifiedTime, &m.LastModifiedById, &m.LastModifiedFromIP,
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

func (r *WorkOrderServiceFeeRepo) GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
    SELECT
        id
    FROM
        work_order_service_fees
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

func (r *WorkOrderServiceFeeRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        work_order_service_fees
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

func (r *WorkOrderServiceFeeRepo) InsertOrUpdateById(ctx context.Context, m *models.WorkOrderServiceFee) error {
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
