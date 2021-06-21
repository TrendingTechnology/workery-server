package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type StaffTagRepo struct {
	db *sql.DB
}

func NewStaffTagRepo(db *sql.DB) *StaffTagRepo {
	return &StaffTagRepo{
		db: db,
	}
}

func (r *StaffTagRepo) Insert(ctx context.Context, m *models.StaffTag) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO staff_tags (
        uuid, tenant_id, staff_id, tag_id, old_id
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
		m.Uuid, m.TenantId, m.StaffId, m.TagId, m.OldId,
	)
	return err
}

func (r *StaffTagRepo) UpdateById(ctx context.Context, m *models.StaffTag) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        staff_tags
    SET
        tenant_id = $1, staff_id = $2, tag_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.StaffId, m.TagId, m.Id,
	)
	return err
}

func (r *StaffTagRepo) GetById(ctx context.Context, id uint64) (*models.StaffTag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.StaffTag)

	query := `
    SELECT
        id, uuid, tenant_id, staff_id, tag_id
	FROM
        staff_tags
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.StaffId, &m.TagId,
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

func (r *StaffTagRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        staff_tags
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

func (r *StaffTagRepo) InsertOrUpdateById(ctx context.Context, m *models.StaffTag) error {
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
