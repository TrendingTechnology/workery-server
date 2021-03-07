package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type SkillSetRepo struct {
	db *sql.DB
}

func NewSkillSetRepo(db *sql.DB) *SkillSetRepo {
	return &SkillSetRepo{
		db: db,
	}
}

func (r *SkillSetRepo) Insert(ctx context.Context, m *models.SkillSet) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO skill_sets (
        uuid, tenant_id, category, sub_category, description, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.Category, m.SubCategory, m.Description, m.State, m.OldId,
	)
	return err
}

func (r *SkillSetRepo) UpdateById(ctx context.Context, m *models.SkillSet) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        skill_sets
    SET
        tenant_id = $1, category = $2, sub_category = $3, description = $4, state = $5
    WHERE
        id = $6`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.Category, m.SubCategory, m.Description, m.State, m.Id,
	)
	return err
}

func (r *SkillSetRepo) GetById(ctx context.Context, id uint64) (*models.SkillSet, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.SkillSet)

	query := `
    SELECT
        id, uuid, tenant_id, category, sub_category, description, state
	FROM
        skill_sets
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Category, &m.SubCategory, &m.Description, &m.State,
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

func (r *SkillSetRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.SkillSet, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.SkillSet)

	query := `
    SELECT
        id, uuid, tenant_id, category, sub_category, description, state
	FROM
        skill_sets
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Category, &m.SubCategory, &m.Description, &m.State,
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

func (r *SkillSetRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        skill_sets
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

func (r *SkillSetRepo) InsertOrUpdateById(ctx context.Context, m *models.SkillSet) error {
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
