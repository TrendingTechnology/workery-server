package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type AssociateInsuranceRequirementRepo struct {
	db *sql.DB
}

func NewAssociateInsuranceRequirementRepo(db *sql.DB) *AssociateInsuranceRequirementRepo {
	return &AssociateInsuranceRequirementRepo{
		db: db,
	}
}

func (r *AssociateInsuranceRequirementRepo) Insert(ctx context.Context, m *models.AssociateInsuranceRequirement) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO associate_insurance_requirements (
        uuid, tenant_id, associate_id, insurance_requirement_id, old_id
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
		m.Uuid, m.TenantId, m.AssociateId, m.InsuranceRequirementId, m.OldId,
	)
	return err
}

func (r *AssociateInsuranceRequirementRepo) UpdateById(ctx context.Context, m *models.AssociateInsuranceRequirement) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        associate_insurance_requirements
    SET
        tenant_id = $1, associate_id = $2, insurance_requirement_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.AssociateId, m.InsuranceRequirementId, m.Id,
	)
	return err
}

func (r *AssociateInsuranceRequirementRepo) GetById(ctx context.Context, id uint64) (*models.AssociateInsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AssociateInsuranceRequirement)

	query := `
    SELECT
        id, uuid, tenant_id, associate_id, insurance_requirement_id
	FROM
        associate_insurance_requirements
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.AssociateId, &m.InsuranceRequirementId,
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

func (r *AssociateInsuranceRequirementRepo) GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
    SELECT
        id
    FROM
        associate_insurance_requirements
    WHERE
		tenant_id = $1
	AND
	    old_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, tid, oid).Scan(&newId)
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

func (r *AssociateInsuranceRequirementRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        associate_insurance_requirements
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

func (r *AssociateInsuranceRequirementRepo) InsertOrUpdateById(ctx context.Context, m *models.AssociateInsuranceRequirement) error {
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
