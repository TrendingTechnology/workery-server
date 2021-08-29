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
        uuid, tenant_id, user_id, type_of, indexed_text, is_ok_to_email,
		is_ok_to_text, is_business, is_senior, is_support, job_info_read,
		how_hear_id, how_hear_old, how_hear_other, state, deactivation_reason,
		deactivation_reason_other, created_time, created_by_id, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_from_ip,
		organization_name, old_id, address_country, address_region,
		address_locality, post_office_box_number, postal_code, street_address,
		street_address_extra, given_name, middle_name, last_name, birthdate,
		join_date, nationality, gender, tax_id, elevation, latitude, longitude,
		area_served, available_language, contact_type, email, fax_number,
		telephone, telephone_type_of, telephone_extension, other_telephone,
		other_telephone_extension, other_telephone_type_of, name, lexical_name
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
		$31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
		$45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.UserId, m.TypeOf, m.IndexedText, m.IsOkToEmail,
		m.IsOkToText, m.IsBusiness, m.IsSenior, m.IsSupport, m.JobInfoRead,
		m.HowHearId, m.HowHearOld, m.HowHearOther, m.State, m.DeactivationReason,
		m.DeactivationReasonOther, m.CreatedTime, m.CreatedById, m.CreatedFromIP,
		m.LastModifiedTime, m.LastModifiedById, m.LastModifiedFromIP,
		m.OrganizationName, m.OldId, m.AddressCountry, m.AddressRegion,
		m.AddressLocality, m.PostOfficeBoxNumber, m.PostalCode, m.StreetAddress,
		m.StreetAddressExtra, m.GivenName, m.MiddleName, m.LastName, m.Birthdate,
		m.JoinDate, m.Nationality, m.Gender, m.TaxId, m.Elevation, m.Latitude, m.Longitude,
		m.AreaServed, m.AvailableLanguage, m.ContactType, m.Email, m.FaxNumber,
		m.Telephone, m.TelephoneTypeOf, m.TelephoneExtension, m.OtherTelephone,
		m.OtherTelephoneExtension, m.OtherTelephoneTypeOf, m.Name, m.LexicalName,
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
        id, uuid, tenant_id, user_id, type_of, indexed_text, is_ok_to_email,
		is_ok_to_text, is_business, is_senior, is_support, job_info_read,
		how_hear_id, how_hear_old, how_hear_other, state, deactivation_reason,
		deactivation_reason_other, created_time, created_by_id, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_from_ip,
		organization_name, old_id, address_country, address_region,
		address_locality, post_office_box_number, postal_code, street_address,
		street_address_extra, given_name, middle_name, last_name, birthdate,
		join_date, nationality, gender, tax_id, elevation, latitude, longitude,
		area_served, available_language, contact_type, email, fax_number,
		telephone, telephone_type_of, telephone_extension, other_telephone,
		other_telephone_extension, other_telephone_type_of, name, lexical_name
	FROM
        customers
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.UserId, &m.TypeOf, &m.IndexedText, &m.IsOkToEmail,
		&m.IsOkToText, &m.IsBusiness, &m.IsSenior, &m.IsSupport, &m.JobInfoRead,
		&m.HowHearId, &m.HowHearOld, &m.HowHearOther, &m.State, &m.DeactivationReason,
		&m.DeactivationReasonOther, &m.CreatedTime, &m.CreatedById, &m.CreatedFromIP,
		&m.LastModifiedTime, &m.LastModifiedById, &m.LastModifiedFromIP,
		&m.OrganizationName, &m.OldId, &m.AddressCountry, &m.AddressRegion,
		&m.AddressLocality, &m.PostOfficeBoxNumber, &m.PostalCode, &m.StreetAddress,
		&m.StreetAddressExtra, &m.GivenName, &m.MiddleName, &m.LastName, &m.Birthdate,
		&m.JoinDate, &m.Nationality, &m.Gender, &m.TaxId, &m.Elevation, &m.Latitude, &m.Longitude,
		&m.AreaServed, &m.AvailableLanguage, &m.ContactType, &m.Email, &m.FaxNumber,
		&m.Telephone, &m.TelephoneTypeOf, &m.TelephoneExtension, &m.OtherTelephone,
		&m.OtherTelephoneExtension, &m.OtherTelephoneTypeOf, &m.Name, &m.LexicalName,
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

func (r *CustomerRepo) GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
    SELECT
        id
    FROM
        customers
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
