package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type BulletinBoardItemRepo struct {
	db *sql.DB
}

func NewBulletinBoardItemRepo(db *sql.DB) *BulletinBoardItemRepo {
	return &BulletinBoardItemRepo{
		db: db,
	}
}

func (r *BulletinBoardItemRepo) Insert(ctx context.Context, m *models.BulletinBoardItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO bulletin_board_items (
        uuid, tenant_id, text, created_time, created_by_id, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_from_ip, state,
		old_id
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
		m.Uuid, m.TenantId, m.Text, m.CreatedTime, m.CreatedById,
		m.CreatedFromIP, m.LastModifiedTime, m.LastModifiedById,
		m.LastModifiedFromIP, m.State, m.OldId,
	)
	return err
}

func (r *BulletinBoardItemRepo) UpdateById(ctx context.Context, m *models.BulletinBoardItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        bulletin_board_items
    SET
        tenant_id = $1, text = $2, created_time = $3, created_by_id = $4,
		created_from_ip = $5, last_modified_time = $6, last_modified_by_id = $7,
		last_modified_from_ip = $8, state = $9
    WHERE
        id = $10`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.Text, m.CreatedTime, m.CreatedById,
		m.CreatedFromIP, m.LastModifiedTime, m.LastModifiedById,
		m.LastModifiedFromIP, m.State, m.Id,
	)
	return err
}

func (r *BulletinBoardItemRepo) GetById(ctx context.Context, id uint64) (*models.BulletinBoardItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.BulletinBoardItem)

	query := `
    SELECT
        id, uuid, tenant_id, text, created_time, created_by_id, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_from_ip, state
	FROM
        bulletin_board_items
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.CreatedTime, &m.CreatedById,
		&m.CreatedFromIP, &m.LastModifiedTime, &m.LastModifiedById,
		&m.LastModifiedFromIP, &m.State,
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

func (r *BulletinBoardItemRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        bulletin_board_items
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

func (r *BulletinBoardItemRepo) InsertOrUpdateById(ctx context.Context, m *models.BulletinBoardItem) error {
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

func (r *BulletinBoardItemRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.BulletinBoardItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.BulletinBoardItem)

	query := `
    SELECT
	    id, uuid, tenant_id, text, created_time, created_by_id, created_from_ip,
    	last_modified_time, last_modified_by_id, last_modified_from_ip, state
    FROM
	    bulletin_board_items
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.CreatedTime, &m.CreatedById,
		&m.CreatedFromIP, &m.LastModifiedTime, &m.LastModifiedById,
		&m.LastModifiedFromIP, &m.State,
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
