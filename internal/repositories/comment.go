package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

func (r *CommentRepo) Insert(ctx context.Context, m *models.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO comments (
        uuid, tenant_id, created_time, created_by_id, created_from_ip, last_modified_time, last_modified_by_id, last_modified_from_ip, text, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.CreatedTime, m.CreatedById, m.CreatedFromIP, m.LastModifiedTime,
		m.LastModifiedById, m.LastModifiedFromIP, m.Text, m.State, m.OldId,
	)
	return err
}

func (r *CommentRepo) UpdateById(ctx context.Context, m *models.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        comments
    SET
        tenant_id = $1, created_time = $2, created_by_id = $3, created_from_ip = $4,
		last_modified_time = $5, last_modified_by_id = $6, last_modified_from_ip = $7, text = $8, state = $9
    WHERE
        id = $10`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.CreatedTime, m.CreatedById, m.CreatedFromIP, m.LastModifiedTime,
		m.LastModifiedById, m.LastModifiedFromIP, m.Text, m.State, m.Id,
	)
	return err
}

func (r *CommentRepo) GetById(ctx context.Context, id uint64) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Comment)

	query := `
    SELECT
        id, uuid, tenant_id, created_time, created_by_id, created_from_ip, last_modified_by_id,
		last_modified_time, last_modified_from_ip, text, state
	FROM
        comments
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.CreatedTime, &m.CreatedById, &m.CreatedFromIP,
		&m.LastModifiedTime, &m.LastModifiedById, &m.LastModifiedFromIP, &m.Text, &m.State, &m.Id,
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

func (r *CommentRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        comments
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

func (r *CommentRepo) InsertOrUpdateById(ctx context.Context, m *models.Comment) error {
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
