package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type PartnerCommentRepo struct {
	db *sql.DB
}

func NewPartnerCommentRepo(db *sql.DB) *PartnerCommentRepo {
	return &PartnerCommentRepo{
		db: db,
	}
}

func (r *PartnerCommentRepo) Insert(ctx context.Context, m *models.PartnerComment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO partner_comments (
        uuid, tenant_id, partner_id, comment_id, old_id
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
		m.Uuid, m.TenantId, m.PartnerId, m.CommentId, m.OldId,
	)
	return err
}

func (r *PartnerCommentRepo) UpdateById(ctx context.Context, m *models.PartnerComment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        partner_comments
    SET
        tenant_id = $1, partner_id = $2, comment_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.PartnerId, m.CommentId, m.Id,
	)
	return err
}

func (r *PartnerCommentRepo) GetById(ctx context.Context, id uint64) (*models.PartnerComment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.PartnerComment)

	query := `
    SELECT
        id, uuid, tenant_id, partner_id, comment_id
	FROM
        partner_comments
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.PartnerId, &m.CommentId,
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

func (r *PartnerCommentRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.PartnerComment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.PartnerComment)

	query := `
    SELECT
        id, uuid, tenant_id, partner_id, comment_id
	FROM
        partner_comments
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.PartnerId, &m.CommentId,
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

func (r *PartnerCommentRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        partner_comments
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

func (r *PartnerCommentRepo) InsertOrUpdateById(ctx context.Context, m *models.PartnerComment) error {
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
