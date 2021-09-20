package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type PrivateFileRepo struct {
	db *sql.DB
}

func NewPrivateFileRepo(db *sql.DB) *PrivateFileRepo {
	return &PrivateFileRepo{
		db: db,
	}
}

func (r *PrivateFileRepo) Insert(ctx context.Context, m *models.PrivateFile) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO private_files (
        uuid, tenant_id, s3_key, title, description, indexed_text, created_time,
		created_from_ip, created_by_id, created_by_name, last_modified_time, last_modified_by_id, last_modified_by_name,
		last_modified_from_ip, associate_id, customer_id, partner_id, staff_id,
		work_order_id, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		$17, $18, $19, $20, $21
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid,
		m.TenantId,
		m.S3Key,
		m.Title,
		m.Description,
		m.IndexedText,
		m.CreatedTime,
		m.CreatedFromIP,
		m.CreatedById,
		m.CreatedByName,
		m.LastModifiedTime,
		m.LastModifiedById,
		m.LastModifiedByName,
		m.LastModifiedFromIP,
		m.AssociateId,
		m.CustomerId,
		m.PartnerId,
		m.StaffId,
		m.WorkOrderId,
		m.State,
		m.OldId,
	)
	return err
}

func (r *PrivateFileRepo) UpdateById(ctx context.Context, m *models.PrivateFile) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        private_files
    SET
        tenant_id = $1, s3_key = $2, title = $3, description = $4,
		indexed_text = $5, last_modified_time = $6, last_modified_by_id = $7,
		last_modified_from_ip = $8, associate_id = $9, customer_id = $10,
		partner_id = $11, staff_id = $12, work_order_id = $13, state = $14,
		last_modified_by_name = $15
    WHERE
        id = $16`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.S3Key, m.Title, m.Description, m.IndexedText,
		m.LastModifiedTime, m.LastModifiedById, m.LastModifiedFromIP,
		m.AssociateId, m.CustomerId, m.PartnerId, m.StaffId, m.WorkOrderId,
		m.State, m.LastModifiedByName, m.Id,
	)
	return err
}

func (r *PrivateFileRepo) GetById(ctx context.Context, id uint64) (*models.PrivateFile, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.PrivateFile)

	query := `
    SELECT
        id, uuid, tenant_id, text, description, state
	FROM
        private_files
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

func (r *PrivateFileRepo) GetIdByOldId(ctx context.Context, tenantId uint64, oldId uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
	SELECT
		id
	FROM
		private_files
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

func (r *PrivateFileRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        private_files
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

func (r *PrivateFileRepo) CheckIfExistsByS3Key(ctx context.Context, s3Key string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        private_files
    WHERE
        s3_key = $1`
	err := r.db.QueryRowContext(ctx, query, s3Key).Scan(&exists)
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

func (r *PrivateFileRepo) InsertOrUpdateById(ctx context.Context, m *models.PrivateFile) error {
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
