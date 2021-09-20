package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type ActivitySheetItemRepo struct {
	db *sql.DB
}

func NewActivitySheetItemRepo(db *sql.DB) *ActivitySheetItemRepo {
	return &ActivitySheetItemRepo{
		db: db,
	}
}

func (r *ActivitySheetItemRepo) Insert(ctx context.Context, m *models.ActivitySheetItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO activity_sheet_items (
        uuid,
		tenant_id,
		comment,
		created_time,
		created_by_id,
		created_by_name,
		created_from_ip,
		associate_id,
		associate_name,
		associate_lexical_name,
		order_id,
		state,
		ongoing_order_id,
		old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.Comment, m.CreatedTime, m.CreatedById, m.CreatedByName,
		m.CreatedFromIP, m.AssociateId, m.AssociateName, m.AssociateLexicalName, m.OrderId, m.State, m.OngoingOrderId,
		m.OldId,
	)
	return err
}

func (r *ActivitySheetItemRepo) UpdateById(ctx context.Context, m *models.ActivitySheetItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        activity_sheet_items
    SET
        comment = $1, associate_id = $2, order_id = $3, state = $4, ongoing_order_id = $5
    WHERE
        id = $6`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Comment, m.AssociateId, m.OrderId, m.State, m.OngoingOrderId,
	)
	return err
}

func (r *ActivitySheetItemRepo) GetById(ctx context.Context, id uint64) (*models.ActivitySheetItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.ActivitySheetItem)

	query := `
    SELECT
        id, uuid, tenant_id,
		comment,
		created_time,
		created_by_id,
		created_by_name,
		created_from_ip,
		associate_id,
		associate_name,
		associate_lexical_name,
		order_id,
		state,
		ongoing_order_id
	FROM
        activity_sheet_items
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Comment, &m.CreatedTime, &m.CreatedById, &m.CreatedByName,
		&m.CreatedFromIP, &m.AssociateId, &m.AssociateName, &m.AssociateLexicalName, &m.OrderId, &m.State, &m.OngoingOrderId,
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

func (r *ActivitySheetItemRepo) GetIdByOldId(ctx context.Context, tenantId uint64, oldId uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
	SELECT
		id
	FROM
		activity_sheet_items
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

func (r *ActivitySheetItemRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        activity_sheet_items
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

func (r *ActivitySheetItemRepo) InsertOrUpdateById(ctx context.Context, m *models.ActivitySheetItem) error {
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
