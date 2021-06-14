package models

import (
	"context"
	"database/sql"
	"time"

	null "gopkg.in/guregu/null.v4"
)

//---------------------
// state
//---------------------
// 1 = Active
// 0 = Inactive

//---------------------
// type_of
//---------------------
// 1 = Unknown Staff
// 2 = Residential Staff
// 3 = Commercial Staff

//---------------------
// organization_type_of
//---------------------
// 1 = Unknown Organization Type | UNKNOWN_ORGANIZATION_TYPE_OF_ID
// 2 = Private Organization Type | PRIVATE_ORGANIZATION_TYPE_OF_ID
// 3 = Non-Profit Organization Type | NON_PROFIT_ORGANIZATION_TYPE_OF_ID
// 4 = Government Organization | GOVERNMENT_ORGANIZATION_TYPE_OF_ID

type Staff struct {
	CreatedTime                          time.Time       `json:"created_time"`
	LastModifiedTime                     time.Time       `json:"last_modified_time"`
	AvailableLanguage                    sql.NullString  `json:"available_language"`
	ContactType                          sql.NullString  `json:"contact_type"`
	Email                                sql.NullString  `json:"email"`
	FaxNumber                            sql.NullString  `json:"fax_number"`
	Telephone                            sql.NullString  `json:"telephone"`
	TelephoneTypeOf                      int8            `json:"telephone_type_of"`
	TelephoneExtension                   sql.NullString  `json:"telephone_extension"`
	OtherTelephone                       sql.NullString  `json:"other_telephone"`
	OtherTelephoneExtension              sql.NullString  `json:"other_telephone_extension"`
	OtherTelephoneTypeOf                 int8            `json:"other_telephone_type_of"`
	AddressCountry                       string          `json:"address_country"`
	AddressLocality                      string          `json:"address_locality"`
	AddressRegion                        string          `json:"address_region"`
	PostOfficeBoxNumber                  sql.NullString  `json:"post_office_box_number"`
	PostalCode                           sql.NullString  `json:"postal_code"`
	StreetAddress                        string          `json:"street_address"`
	StreetAddressExtra                   sql.NullString  `json:"street_address_extra"`
	Elevation                            sql.NullFloat64 `json:"elevation"`
	Latitude                             sql.NullFloat64 `json:"latitude"`
	Longitude                            sql.NullFloat64 `json:"longitude"`
	GivenName                            sql.NullString  `json:"given_name"`
	MiddleName                           sql.NullString  `json:"middle_name"`
	LastName                             sql.NullString  `json:"last_name"`
	Birthdate                            sql.NullTime    `json:"birthdate"`
	JoinDate                             sql.NullTime    `json:"join_date"`
	Nationality                          sql.NullString  `json:"nationality"`
	Gender                               sql.NullString  `json:"gender"`
	TaxId                                sql.NullString  `json:"tax_id"`
	TenantId                             uint64          `json:"tenant_id"`
	Id                                   uint64          `json:"id"`
	Uuid                                 string          `json:"uuid"`
	IndexedText                          sql.NullString  `json:"indexed_text"`
	CreatedFromIP                        sql.NullString  `json:"created_from_ip"`
	LastModifiedFromIP                   sql.NullString  `json:"last_modified_from_ip"`
	State                                int8            `json:"state"`
	CreatedById                          null.Int        `json:"created_by_id"`
	LastModifiedById                     null.Int        `json:"last_modified_by_id"`
	UserId                               uint64          `json:"user_id"`
	HowHearOther                         sql.NullString  `json:"how_hear_other"`
	HowHearId                            sql.NullInt64   `json:"how_hear_id"`
	AvatarImageId                        sql.NullInt64   `json:"avatar_image_id"`
	PersonalEmail                        sql.NullString  `json:"personal_email"`
	EmergencyContactAlternativeTelephone sql.NullString  `json:"emergency_contact_alternative_telephone"`
	EmergencyContactName                 sql.NullString  `json:"emergency_contact_name"`
	EmergencyContactRelationship         sql.NullString  `json:"emergency_contact_relationship"`
	EmergencyContactTelephone            sql.NullString  `json:"emergency_contact_telephone"`
	PoliceCheck                          sql.NullTime    `json:"police_check"`
	OldId                                uint64          `json:"old_id"`
	IsOkToEmail                          bool            `json:"is_ok_to_email"`
	IsOkToText                           bool            `json:"is_ok_to_text"`
}

type StaffRepository interface {
	Insert(ctx context.Context, u *Staff) error
	UpdateById(ctx context.Context, u *Staff) error
	GetById(ctx context.Context, id uint64) (*Staff, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *Staff) error
}
