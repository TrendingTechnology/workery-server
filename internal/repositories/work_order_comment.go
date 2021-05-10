package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type WorkOrderCommentRepo struct {
	db *sql.DB
}

func NewWorkOrderCommentRepo(db *sql.DB) *WorkOrderCommentRepo {
	return &WorkOrderCommentRepo{
		db: db,
	}
}

func (r *WorkOrderCommentRepo) Insert(ctx context.Context, m *models.WorkOrderComment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO work_order_comments (
        uuid, tenant_id, order_id, comment_id, old_id
    ) VALUES (
        $1, $2, $3, $4, $5
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.OrderId, m.CommentId, m.OldId,
	)
	return err
}

func (r *WorkOrderCommentRepo) UpdateById(ctx context.Context, m *models.WorkOrderComment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        work_order_comments
    SET
        tenant_id = $1, order_id = $2, comment_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.OrderId, m.CommentId, m.Id,
	)
	return err
}

func (r *WorkOrderCommentRepo) GetById(ctx context.Context, id uint64) (*models.WorkOrderComment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderComment)

	query := `
    SELECT
        id, uuid, tenant_id, order_id, comment_id
	FROM
        work_order_comments
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.OrderId, &m.CommentId,
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

func (r *WorkOrderCommentRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.WorkOrderComment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderComment)

	query := `
    SELECT
        id, uuid, tenant_id, order_id, comment_id
	FROM
        work_order_comments
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.OrderId, &m.CommentId,
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

func (r *WorkOrderCommentRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        work_order_comments
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

func (r *WorkOrderCommentRepo) InsertOrUpdateById(ctx context.Context, m *models.WorkOrderComment) error {
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
