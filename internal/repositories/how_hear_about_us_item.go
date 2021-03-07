package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type HowHearAboutUsItemRepo struct {
	db *sql.DB
}

func NewHowHearAboutUsItemRepo(db *sql.DB) *HowHearAboutUsItemRepo {
	return &HowHearAboutUsItemRepo{
		db: db,
	}
}

func (r *HowHearAboutUsItemRepo) Insert(ctx context.Context, m *models.HowHearAboutUsItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO how_hear_about_us_items (
        uuid, tenant_id, text, sort_number, is_for_associate, is_for_customer,
		is_for_staff, is_for_partner, state, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.Text, m.SortNumber, m.IsForAssociate, m.IsForCustomer,
		m.IsForStaff, m.IsForPartner, m.State, m.OldId,
	)
	return err
}

func (r *HowHearAboutUsItemRepo) UpdateById(ctx context.Context, m *models.HowHearAboutUsItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        how_hear_about_us_items
    SET
        tenant_id = $1, text = $2, sort_number = $3, is_for_associate = $4,
		is_for_customer = $5, is_for_staff = $6, is_for_partner = $7, state = $8
    WHERE
        id = $9`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.Text, m.SortNumber, m.IsForAssociate, m.IsForCustomer,
		m.IsForStaff, m.IsForPartner, m.State, m.Id,
	)
	return err
}

func (r *HowHearAboutUsItemRepo) GetById(ctx context.Context, id uint64) (*models.HowHearAboutUsItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.HowHearAboutUsItem)

	query := `
    SELECT
        id, uuid, tenant_id, text, sort_number, is_for_associate, is_for_customer,
		is_for_staff, is_for_partner, state, old_id
	FROM
        how_hear_about_us_items
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.SortNumber, &m.IsForAssociate, &m.IsForCustomer,
		&m.IsForStaff, &m.IsForPartner, &m.State, &m.OldId,
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

func (r *HowHearAboutUsItemRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.HowHearAboutUsItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.HowHearAboutUsItem)

	query := `
    SELECT
        id, uuid, tenant_id, text, sort_number, is_for_associate, is_for_customer,
		is_for_staff, is_for_partner, state, old_id
	FROM
        how_hear_about_us_items
    WHERE
        old_id = $1`
	err := r.db.QueryRowContext(ctx, query, oldId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.SortNumber, &m.IsForAssociate, &m.IsForCustomer,
		&m.IsForStaff, &m.IsForPartner, &m.State, &m.OldId,
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

func (r *HowHearAboutUsItemRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        how_hear_about_us_items
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

func (r *HowHearAboutUsItemRepo) InsertOrUpdateById(ctx context.Context, m *models.HowHearAboutUsItem) error {
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
