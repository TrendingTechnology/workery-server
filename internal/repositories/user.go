package repositories

import (
    "context"
	"database/sql"
    "time"

    "github.com/over55/workery-server/internal/models"
)

type UserRepo struct {
    db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
    return &UserRepo{
        db: db,
    }
}

func (r *UserRepo) Insert(ctx context.Context, m *models.User) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (uuid, email, first_name, last_name, password_hash, state, timezone, created_time, session_uuid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid,
		m.Email,
		m.FirstName,
        m.LastName,
        m.PasswordHash,
        m.State,
        m.Timezone,
        m.CreatedTime,
        m.SessionUuid,
	)
	return err
}

func (r *UserRepo) Update(ctx context.Context, m *models.User) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "UPDATE users SET email = $1, first_name = $2, last_name = $3, password_hash = $4, state = $5, timezone = $6, created_time = $7, session_uuid = $8 WHERE id = $9"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Email,
		m.FirstName,
		m.LastName,
        m.PasswordHash,
        m.State,
        m.Timezone,
        m.CreatedTime,
        m.SessionUuid,
        m.Id,
	)
	return err
}

func (r *UserRepo) GetById(ctx context.Context, id uint64) (*models.User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.User)

	query := "SELECT id, uuid, email, first_name, last_name, password_hash, state, timezone, created_time, session_uuid FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(
        &m.Id,
        &m.Uuid,
		&m.Email,
		&m.FirstName,
        &m.LastName,
        &m.PasswordHash,
        &m.State,
        &m.Timezone,
        &m.CreatedTime,
        &m.SessionUuid,
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

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.User)

	query := "SELECT id, uuid, email, first_name, last_name, password_hash, state, timezone, created_time, session_uuid FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(
        &m.Id,
        &m.Uuid,
		&m.Email,
		&m.FirstName,
        &m.LastName,
        &m.PasswordHash,
        &m.State,
        &m.Timezone,
        &m.CreatedTime,
        &m.SessionUuid,
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

func (r *UserRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

    var exists bool

    query := `SELECT 1 FROM users WHERE id = $1;`

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

func (r *UserRepo) InsertOrUpdate(ctx context.Context, m *models.User) error {
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
