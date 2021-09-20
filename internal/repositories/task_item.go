package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type TaskItemRepo struct {
	db *sql.DB
}

func NewTaskItemRepo(db *sql.DB) *TaskItemRepo {
	return &TaskItemRepo{
		db: db,
	}
}

func (r *TaskItemRepo) Insert(ctx context.Context, m *models.TaskItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO task_items (
        uuid, tenant_id, type_of, title, description, due_date, is_closed,
		was_postponed, closing_reason, closing_reason_other, created_time,
		created_from_ip, created_by_id, created_by_name, last_modified_time,
		last_modified_from_ip, last_modified_by_id, last_modified_by_name, order_id, order_type_of,
		ongoing_order_id, state, customer_id, customer_name, customer_lexical_name,
		associate_id, associate_name, associate_lexical_name, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.TypeOf, m.Title, m.Description, m.DueDate, m.IsClosed,
		m.WasPostponed, m.ClosingReason, m.ClosingReasonOther, m.CreatedTime,
		m.CreatedFromIP, m.CreatedById, m.CreatedByName, m.LastModifiedTime, m.LastModifiedFromIP,
		m.LastModifiedById, m.LastModifiedByName, m.OrderId, m.OrderTypeOf, m.OngoingOrderId,
		m.State, m.CustomerId, m.CustomerName, m.CustomerLexicalName,
		m.AssociateId, m.AssociateName, m.AssociateLexicalName, m.OldId,
	)
	return err
}

func (r *TaskItemRepo) UpdateById(ctx context.Context, m *models.TaskItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        task_items
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
		m.TenantId, m.Description, m.State, m.Id,
	)
	return err
}

func (r *TaskItemRepo) GetById(ctx context.Context, id uint64) (*models.TaskItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.TaskItem)

	query := `
    SELECT
        id, uuid, tenant_id, text, description, state
	FROM
        task_items
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Description, &m.State,
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

func (r *TaskItemRepo) GetIdByOldId(ctx context.Context, tenantId uint64, oldId uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
	SELECT
		id
	FROM
		task_items
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

func (r *TaskItemRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        task_items
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

func (r *TaskItemRepo) InsertOrUpdateById(ctx context.Context, m *models.TaskItem) error {
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
