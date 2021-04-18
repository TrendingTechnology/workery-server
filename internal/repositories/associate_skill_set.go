package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type AssociateSkillSetRepo struct {
	db *sql.DB
}

func NewAssociateSkillSetRepo(db *sql.DB) *AssociateSkillSetRepo {
	return &AssociateSkillSetRepo{
		db: db,
	}
}

func (r *AssociateSkillSetRepo) Insert(ctx context.Context, m *models.AssociateSkillSet) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO associate_skill_sets (
        uuid, tenant_id, associate_id, skill_set_id, old_id
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
		m.Uuid, m.TenantId, m.AssociateId, m.SkillSetId, m.OldId,
	)
	return err
}

func (r *AssociateSkillSetRepo) UpdateById(ctx context.Context, m *models.AssociateSkillSet) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        associate_skill_sets
    SET
        tenant_id = $1, associate_id = $2, skill_set_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.AssociateId, m.SkillSetId, m.Id,
	)
	return err
}

func (r *AssociateSkillSetRepo) GetById(ctx context.Context, id uint64) (*models.AssociateSkillSet, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AssociateSkillSet)

	query := `
    SELECT
        id, uuid, tenant_id, associate_id, skill_set_id
	FROM
        associate_skill_sets
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.AssociateId, &m.SkillSetId,
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

func (r *AssociateSkillSetRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.AssociateSkillSet, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AssociateSkillSet)

	query := `
    SELECT
        id, uuid, tenant_id, associate_id, skill_set_id
	FROM
        associate_skill_sets
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.AssociateId, &m.SkillSetId,
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

func (r *AssociateSkillSetRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        associate_skill_sets
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

func (r *AssociateSkillSetRepo) InsertOrUpdateById(ctx context.Context, m *models.AssociateSkillSet) error {
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
