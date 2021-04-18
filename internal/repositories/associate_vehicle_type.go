package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type AssociateVehicleTypeRepo struct {
	db *sql.DB
}

func NewAssociateVehicleTypeRepo(db *sql.DB) *AssociateVehicleTypeRepo {
	return &AssociateVehicleTypeRepo{
		db: db,
	}
}

func (r *AssociateVehicleTypeRepo) Insert(ctx context.Context, m *models.AssociateVehicleType) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO associate_vehicle_types (
        uuid, tenant_id, associate_id, vehicle_type_id, old_id
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
		m.Uuid, m.TenantId, m.AssociateId, m.VehicleTypeId, m.OldId,
	)
	return err
}

func (r *AssociateVehicleTypeRepo) UpdateById(ctx context.Context, m *models.AssociateVehicleType) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        associate_vehicle_types
    SET
        tenant_id = $1, associate_id = $2, vehicle_type_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.AssociateId, m.VehicleTypeId, m.Id,
	)
	return err
}

func (r *AssociateVehicleTypeRepo) GetById(ctx context.Context, id uint64) (*models.AssociateVehicleType, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AssociateVehicleType)

	query := `
    SELECT
        id, uuid, tenant_id, associate_id, vehicle_type_id
	FROM
        associate_vehicle_types
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.AssociateId, &m.VehicleTypeId,
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

func (r *AssociateVehicleTypeRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.AssociateVehicleType, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AssociateVehicleType)

	query := `
    SELECT
        id, uuid, tenant_id, associate_id, vehicle_type_id
	FROM
        associate_vehicle_types
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.AssociateId, &m.VehicleTypeId,
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

func (r *AssociateVehicleTypeRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        associate_vehicle_types
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

func (r *AssociateVehicleTypeRepo) InsertOrUpdateById(ctx context.Context, m *models.AssociateVehicleType) error {
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
