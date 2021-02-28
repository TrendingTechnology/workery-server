package repositories

import (
    "context"
	"database/sql"
    "time"

    "github.com/over55/workery-server/internal/models"
)

type TenantRepo struct {
    db *sql.DB
}

func NewTenantRepo(db *sql.DB) *TenantRepo {
    return &TenantRepo{
        db: db,
    }
}

func (r *TenantRepo) Insert(ctx context.Context, m *models.Tenant) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO tenants (
        uuid, name, state, timezone, created_time, modified_time
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    )
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid,
		m.Name,
        m.State,
        m.Timezone,
        m.CreatedTime,
        m.ModifiedTime,
	)
	return err
}

func (r *TenantRepo) Update(ctx context.Context, m *models.Tenant) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        tenants
    SET
        name = $1, state = $2, timezone = $3, created_time = $4, modified_time = $5
    WHERE
        id = $6
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Name,
        m.State,
        m.Timezone,
        m.CreatedTime,
        m.ModifiedTime,
        m.Id,
	)
	return err
}

func (r *TenantRepo) GetById(ctx context.Context, id uint64) (*models.Tenant, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Tenant)

	query := `
    SELECT
        id, uuid, name, state, timezone, created_time, modified_time
    FROM
        tenants
    WHERE
        id = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
        &m.Id,
        &m.Uuid,
        &m.Name,
        &m.State,
        &m.Timezone,
        &m.CreatedTime,
        &m.ModifiedTime,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *TenantRepo) GetByName(ctx context.Context, name string) (*models.Tenant, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Tenant)

	query := `
    SELECT
        id, uuid, name, state, timezone, created_time, modified_time
    FROM
        tenants
    WHERE
        name = $1
    `
	err := r.db.QueryRowContext(ctx, query, name).Scan(
        &m.Id,
        &m.Uuid,
        &m.Name,
        &m.State,
        &m.Timezone,
        &m.CreatedTime,
        &m.ModifiedTime,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that name.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *TenantRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

    var exists bool

    query := `
    SELECT
        1
    FROM
        tenants
    WHERE
        id = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *TenantRepo) CheckIfExistsByName(ctx context.Context, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

    var exists bool

    query := `
    SELECT
        1
    FROM
        tenants
    WHERE
        name = $1
    `

	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *TenantRepo) InsertOrUpdate(ctx context.Context, m *models.Tenant) error {
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
    return r.Update(ctx, m)
}
