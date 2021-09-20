package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type OngoingWorkOrderRepo struct {
	db *sql.DB
}

func NewOngoingWorkOrderRepo(db *sql.DB) *OngoingWorkOrderRepo {
	return &OngoingWorkOrderRepo{
		db: db,
	}
}

func (r *OngoingWorkOrderRepo) Insert(ctx context.Context, m *models.OngoingWorkOrder) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO ongoing_work_orders (
        uuid,
		tenant_id,
		customer_id,
		customer_name,
		customer_lexical_name,
		associate_id,
		associate_name,
		associate_lexical_name,
		state,
		created_time,
		created_by_id,
		created_by_name,
		created_from_ip,
		last_modified_time,
		last_modified_by_id,
		last_modified_by_name,
		last_modified_from_ip, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid,
		m.TenantId,
		m.CustomerId,
		m.CustomerName,
		m.CustomerLexicalName,
		m.AssociateId,
		m.AssociateName,
		m.AssociateLexicalName,
		m.State,
		m.CreatedTime,
		m.CreatedById,
		m.CreatedByName,
		m.CreatedFromIP,
		m.LastModifiedTime,
		m.LastModifiedById,
		m.LastModifiedByName,
		m.LastModifiedFromIP,
		m.OldId,
	)
	return err
}

func (r *OngoingWorkOrderRepo) UpdateById(ctx context.Context, m *models.OngoingWorkOrder) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        ongoing_work_orders
    SET
        tenant_id = $1,
		customer_id = $2,
		associate_id = $3,
		state = $4,
		last_modified_time = $5,
		last_modified_by_id = $6,
		last_modified_by_name = $7,
		last_modified_from_ip = $8
    WHERE
        id = $9`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId,
		m.CustomerId,
		m.AssociateId,
		m.State,
		m.LastModifiedTime,
		m.LastModifiedById,
		m.LastModifiedByName,
		m.LastModifiedFromIP,
		m.Id,
	)
	return err
}

func (r *OngoingWorkOrderRepo) GetById(ctx context.Context, id uint64) (*models.OngoingWorkOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.OngoingWorkOrder)

	query := `
    SELECT
        id,
		uuid,
		tenant_id,
		customer_id,
		customer_name,
		customer_lexical_name,
		associate_id,
		associate_name,
		associate_lexical_name,
		state,
		created_time,
		created_by_id,
		created_by_name,
		created_from_ip,
		last_modified_time,
		last_modified_by_id,
		last_modified_by_name,
		last_modified_from_ip
	FROM
        ongoing_work_orders
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id,
		&m.Uuid,
		&m.TenantId,
		&m.CustomerId,
		&m.CustomerName,
		&m.CustomerLexicalName,
		&m.AssociateId,
		&m.AssociateName,
		&m.AssociateLexicalName,
		&m.State,
		&m.CreatedTime,
		&m.CreatedById,
		&m.CreatedByName,
		&m.CreatedFromIP,
		&m.LastModifiedTime,
		&m.LastModifiedById,
		&m.LastModifiedByName,
		&m.LastModifiedFromIP,
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

func (r *OngoingWorkOrderRepo) GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
    SELECT
        id
    FROM
        ongoing_work_orders
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

func (r *OngoingWorkOrderRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        ongoing_work_orders
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

func (r *OngoingWorkOrderRepo) InsertOrUpdateById(ctx context.Context, m *models.OngoingWorkOrder) error {
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
