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
        uuid, name, state, timezone, created_time, modified_time,
        address_country, address_region, address_locality, post_office_box_number,
        postal_code, street_address, street_address_extra, elevation, latitude,
        longitude, area_served, available_language, contact_type, email, fax_number,
        telephone, telephone_type_of, telephone_extension, other_telephone,
        other_telephone_extension, other_telephone_type_of
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
        $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27,
    )
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.Name, m.State, m.Timezone, m.CreatedTime, m.ModifiedTime,
		m.AddressCountry, m.AddressRegion, m.AddressLocality, m.PostOfficeBoxNumber,
		m.PostalCode, m.StreetAddress, m.StreetAddressExtra, m.Elevation, m.Latitude,
		m.Longitude, m.AreaServed, m.AvailableLanguage, m.ContactType, m.Email, m.FaxNumber,
		m.Telephone, m.TelephoneTypeOf, m.TelephoneExtension, m.OtherTelephone,
		m.OtherTelephoneExtension, m.OtherTelephoneTypeOf,
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
        name = $1, state = $2, timezone = $3, created_time = $4, modified_time = $5,
        address_country = $6, address_region = $7, address_locality = $8,
        post_office_box_number = $9, postal_code = $10, street_address = $11,
        street_address_extra = $12, elevation = $13, latitude = $14, longitude = $15,
        area_served = $16, available_language = $17, contact_type = $18, email = $19,
        fax_number = $20, telephone = $20, telephone_type_of = $21,
        telephone_extension = $22, other_telephone = $23, other_telephone_extension = $24,
        other_telephone_type_of = $25
    WHERE
        id = $26
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Name, m.State, m.Timezone, m.CreatedTime, m.ModifiedTime,
		m.AddressCountry, m.AddressRegion, m.AddressLocality,
        m.PostOfficeBoxNumber, m.PostalCode, m.StreetAddress,
		m.StreetAddressExtra, m.Elevation, m.Latitude, m.Longitude,
        m.AreaServed, m.AvailableLanguage, m.ContactType, m.Email,
		m.FaxNumber, m.Telephone, m.TelephoneTypeOf,
        m.TelephoneExtension, m.OtherTelephone, m.OtherTelephoneExtension,
		m.OtherTelephoneTypeOf, m.Id,
	)
	return err
}

func (r *TenantRepo) GetById(ctx context.Context, id uint64) (*models.Tenant, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Tenant)

	query := `
    SELECT
        id, uuid, name, state, timezone, created_time, modified_time,
        address_country, address_region, address_locality, post_office_box_number,
        postal_code, street_address, street_address_extra, elevation, latitude,
        longitude, area_served, available_language, contact_type, email, fax_number,
        telephone, telephone_type_of, telephone_extension, other_telephone,
        other_telephone_extension, other_telephone_type_of
    FROM
        tenants
    WHERE
        id = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.Name, &m.State, &m.Timezone, &m.CreatedTime, &m.ModifiedTime,
        &m.AddressCountry, &m.AddressRegion, &m.AddressLocality,
        &m.PostOfficeBoxNumber, &m.PostalCode, &m.StreetAddress,
		&m.StreetAddressExtra, &m.Elevation, &m.Latitude, &m.Longitude,
        &m.AreaServed, &m.AvailableLanguage, &m.ContactType, &m.Email,
		&m.FaxNumber, &m.Telephone, &m.TelephoneTypeOf,
        &m.TelephoneExtension, &m.OtherTelephone, &m.OtherTelephoneExtension,
		&m.OtherTelephoneTypeOf,
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
        id, uuid, name, state, timezone, created_time, modified_time,
        address_country, address_region, address_locality, post_office_box_number,
        postal_code, street_address, street_address_extra, elevation, latitude,
        longitude, area_served, available_language, contact_type, email, fax_number,
        telephone, telephone_type_of, telephone_extension, other_telephone,
        other_telephone_extension, other_telephone_type_of
    FROM
        tenants
    WHERE
        name = $1
    `
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&m.Id, &m.Uuid, &m.Name, &m.State, &m.Timezone, &m.CreatedTime, &m.ModifiedTime,
        &m.AddressCountry, &m.AddressRegion, &m.AddressLocality,
        &m.PostOfficeBoxNumber, &m.PostalCode, &m.StreetAddress,
		&m.StreetAddressExtra, &m.Elevation, &m.Latitude, &m.Longitude,
        &m.AreaServed, &m.AvailableLanguage, &m.ContactType, &m.Email,
		&m.FaxNumber, &m.Telephone, &m.TelephoneTypeOf,
        &m.TelephoneExtension, &m.OtherTelephone, &m.OtherTelephoneExtension,
		&m.OtherTelephoneTypeOf,
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
