package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type CustomerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

func (r *CustomerRepo) Insert(ctx context.Context, m *models.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO customers (
        uuid, tenant_id, user_id
    ) VALUES (
        $1, $2, $3
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.UserId,
	)
	return err
}

func (r *CustomerRepo) UpdateById(ctx context.Context, m *models.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        customers
    SET
        tenant_id = $1, user_id = $2
    WHERE
        id = $3`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.UserId, m.Id,
	)
	return err
}

func (r *CustomerRepo) GetById(ctx context.Context, id uint64) (*models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Customer)

	query := `
    SELECT
        id, uuid, tenant_id, user_id
	FROM
        customers
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.UserId,
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

func (r *CustomerRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        customers
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

func (r *CustomerRepo) InsertOrUpdateById(ctx context.Context, m *models.Customer) error {
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
