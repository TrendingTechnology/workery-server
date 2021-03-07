package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type InsuranceRequirementRepo struct {
	db *sql.DB
}

func NewInsuranceRequirementRepo(db *sql.DB) *InsuranceRequirementRepo {
	return &InsuranceRequirementRepo{
		db: db,
	}
}

func (r *InsuranceRequirementRepo) Insert(ctx context.Context, m *models.InsuranceRequirement) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO insurance_requirements (
        uuid, tenant_id, text, description, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.Text, m.Description, m.State, m.OldId,
	)
	return err
}

func (r *InsuranceRequirementRepo) UpdateById(ctx context.Context, m *models.InsuranceRequirement) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        insurance_requirements
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
		m.TenantId, m.Text, m.Description, m.State, m.Id,
	)
	return err
}

func (r *InsuranceRequirementRepo) GetById(ctx context.Context, id uint64) (*models.InsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.InsuranceRequirement)

	query := `
    SELECT
        id, uuid, tenant_id, text, description, state
	FROM
        insurance_requirements
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.Description, &m.State,
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

func (r *InsuranceRequirementRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.InsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.InsuranceRequirement)

	query := `
    SELECT
        id, uuid, tenant_id, text, description, state
	FROM
        insurance_requirements
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.Description, &m.State,
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

func (r *InsuranceRequirementRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        insurance_requirements
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

func (r *InsuranceRequirementRepo) InsertOrUpdateById(ctx context.Context, m *models.InsuranceRequirement) error {
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
