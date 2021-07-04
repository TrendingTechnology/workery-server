package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type PublicImageUploadRepo struct {
	db *sql.DB
}

func NewPublicImageUploadRepo(db *sql.DB) *PublicImageUploadRepo {
	return &PublicImageUploadRepo{
		db: db,
	}
}

func (r *PublicImageUploadRepo) Insert(ctx context.Context, m *models.PublicImageUpload) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO public_image_uploads (
        uuid, tenant_id, image_file, created_time, created_from_ip,
		last_modified_time, last_modified_from_ip, created_by_id,
		last_modified_by_id, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//    time.Time `json:""`
	//  string    `json:""`
	//         uint64    `json:""`
	//    uint64    `json:""`
	// State              int8      `json:"state"`
	// OldId              uint64    `json:"old_id"`

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.ImageFile, m.CreatedTime, m.CreatedFromIP,
		m.LastModifiedTime,
		m.LastModifiedFromIP,
		m.CreatedById,
		m.LastModifiedById,
		m.State, m.OldId,
	)
	return err
}

func (r *PublicImageUploadRepo) UpdateById(ctx context.Context, m *models.PublicImageUpload) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        tags
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
		m.TenantId, m.State, m.Id,
	)
	return err
}

func (r *PublicImageUploadRepo) GetById(ctx context.Context, id uint64) (*models.PublicImageUpload, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.PublicImageUpload)

	query := `
    SELECT
        id, uuid, tenant_id, text, description, state
	FROM
        tags
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.State,
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

func (r *PublicImageUploadRepo) GetIdByOldId(ctx context.Context, tenantId uint64, oldId uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
	SELECT
		id
	FROM
		tags
	WHERE
		tenant_id = $1 AND old_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, tenantId, oldId).Scan(&newId)
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

func (r *PublicImageUploadRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        tags
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

func (r *PublicImageUploadRepo) InsertOrUpdateById(ctx context.Context, m *models.PublicImageUpload) error {
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
