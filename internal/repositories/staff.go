package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type StaffRepo struct {
	db *sql.DB
}

func NewStaffRepo(db *sql.DB) *StaffRepo {
	return &StaffRepo{
		db: db,
	}
}

func (r *StaffRepo) Insert(ctx context.Context, m *models.Staff) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO staff (
		old_id, created_time, last_modified_time, available_language, contact_type, email, fax_number,
		telephone, telephone_type_of, telephone_extension,
		other_telephone, other_telephone_extension, other_telephone_type_of,
		address_country, address_locality, address_region, post_office_box_number,
		postal_code, street_address, street_address_extra, elevation, latitude,
		longitude, given_name, middle_name, last_name, birthdate, join_date,
		nationality, gender, tax_id, indexed_text,
		created_from_ip, last_modified_from_ip,
		state, created_by_id, last_modified_by_id,
		user_id, how_hear_other, how_hear_id, avatar_image_id, personal_email,
		emergency_contact_alternative_telephone, emergency_contact_name,
		emergency_contact_relationship, emergency_contact_telephone, police_check,
		uuid, tenant_id, is_ok_to_email, is_ok_to_text
	) VALUES (
	    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
		$31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
		$45, $46, $47, $48, $49, $50, $51
	)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.OldId, m.CreatedTime, m.LastModifiedTime, m.AvailableLanguage, m.ContactType, m.Email, m.FaxNumber,
		m.Telephone, m.TelephoneTypeOf, m.TelephoneExtension,
		m.OtherTelephone, m.OtherTelephoneExtension, m.OtherTelephoneTypeOf,
		m.AddressCountry, m.AddressLocality, m.AddressRegion, m.PostOfficeBoxNumber,
		m.PostalCode, m.StreetAddress, m.StreetAddressExtra, m.Elevation,
		m.Latitude, m.Longitude, m.GivenName, m.MiddleName, m.LastName,
		m.Birthdate, m.JoinDate, m.Nationality, m.Gender, m.TaxId, m.IndexedText,
		m.CreatedFromIP, m.LastModifiedFromIP,
		m.State, m.CreatedById, m.LastModifiedById,
		m.UserId, m.HowHearOther, m.HowHearId, m.AvatarImageId, m.PersonalEmail,
		m.EmergencyContactAlternativeTelephone, m.EmergencyContactName,
		m.EmergencyContactRelationship, m.EmergencyContactTelephone, m.PoliceCheck,
		m.Uuid, m.TenantId, m.IsOkToEmail, m.IsOkToText,
	)
	return err
}

func (r *StaffRepo) UpdateById(ctx context.Context, m *models.Staff) error {
	return nil
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	//
	// query := `
	// UPDATE
	//     staff
	// SET
	//     tenant_id = $1, user_id = $2
	// WHERE
	//     id = $3`
	// stmt, err := r.db.PrepareContext(ctx, query)
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()
	//
	// _, err = stmt.ExecContext(
	// 	ctx,
	// 	m.TenantId, m.UserId, m.Id,
	// )
	// return err
}

func (r *StaffRepo) GetById(ctx context.Context, id uint64) (*models.Staff, error) {
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	//
	// m := new(models.Staff)
	//
	// query := `
	// SELECT
	//     id, uuid, tenant_id, user_id
	// FROM
	//     staff
	// WHERE
	//     id = $1`
	// err := r.db.QueryRowContext(ctx, query, id).Scan(
	// 	&m.Id, &m.Uuid, &m.TenantId, &m.UserId,
	// )
	// if err != nil {
	// 	// CASE 1 OF 2: Cannot find record with that email.
	// 	if err == sql.ErrNoRows {
	// 		return nil, nil
	// 	} else { // CASE 2 OF 2: All other errors.
	// 		return nil, err
	// 	}
	// }
	// return m, nil
	return nil, nil
}

func (r *StaffRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        staff
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

func (r *StaffRepo) GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var newId uint64

	query := `
    SELECT
        id
    FROM
        staff
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

func (r *StaffRepo) InsertOrUpdateById(ctx context.Context, m *models.Staff) error {
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
