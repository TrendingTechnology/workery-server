package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type SkillSetInsuranceRequirementRepo struct {
	db *sql.DB
}

func NewSkillSetInsuranceRequirementRepo(db *sql.DB) *SkillSetInsuranceRequirementRepo {
	return &SkillSetInsuranceRequirementRepo{
		db: db,
	}
}

func (r *SkillSetInsuranceRequirementRepo) Insert(ctx context.Context, m *models.SkillSetInsuranceRequirement) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO skill_set_insurance_requirements (
        uuid, tenant_id, skill_set_id, insurance_requirement_id, old_id
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
		m.Uuid, m.TenantId, m.SkillSetId, m.InsuranceRequirementId, m.OldId,
	)
	return err
}

func (r *SkillSetInsuranceRequirementRepo) UpdateById(ctx context.Context, m *models.SkillSetInsuranceRequirement) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        skill_set_insurance_requirements
    SET
        tenant_id = $1, skill_set_id = $2, insurance_requirement_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.SkillSetId, m.InsuranceRequirementId, m.Id,
	)
	return err
}

func (r *SkillSetInsuranceRequirementRepo) GetById(ctx context.Context, id uint64) (*models.SkillSetInsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.SkillSetInsuranceRequirement)

	query := `
    SELECT
        id, uuid, tenant_id, skill_set_id, insurance_requirement_id
	FROM
        skill_set_insurance_requirements
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.SkillSetId, &m.InsuranceRequirementId,
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

func (r *SkillSetInsuranceRequirementRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.SkillSetInsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.SkillSetInsuranceRequirement)

	query := `
    SELECT
        id, uuid, tenant_id, skill_set_id, insurance_requirement_id
	FROM
        skill_set_insurance_requirements
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.SkillSetId, &m.InsuranceRequirementId,
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

func (r *SkillSetInsuranceRequirementRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        skill_set_insurance_requirements
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

func (r *SkillSetInsuranceRequirementRepo) InsertOrUpdateById(ctx context.Context, m *models.SkillSetInsuranceRequirement) error {
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
