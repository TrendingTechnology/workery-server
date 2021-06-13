package models

import (
	"context"
	"database/sql"
	"time"
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
	Created                              time.Time       `json:"created"`
	LastModified                         time.Time       `json:"last_modified"`
	AlternateName                        sql.NullString  `json:"alternate_name"`
	Description                          sql.NullString  `json:"description"`
	Name                                 sql.NullString  `json:"name"`
	Url                                  sql.NullString  `json:"url"`
	AreaServed                           sql.NullString  `json:"area_served"`
	AvailableLanguage                    sql.NullString  `json:"available_language"`
	ContactType                          sql.NullString  `json:"contact_type"`
	Email                                sql.NullString  `json:"email"`
	FaxNumber                            sql.NullString  `json:"fax_number"`
	ProductSupported                     sql.NullString  `json:"product_supported"`
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
	CreatedFrom                          sql.NullString  `json:"created_from"`
	CreatedFromIsPublic                  bool            `json:"created_from_is_public"`
	LastModifiedFrom                     sql.NullString  `json:"last_modified_from"`
	LastModifiedFromIsPublic             bool            `json:"last_modified_from_is_public"`
	IsArchived                           bool            `json:"is_archived"`
	CreatedById                          sql.NullInt64   `json:"created_by_id"`
	LastModifiedById                     sql.NullInt64   `json:"last_modified_by_id"`
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
}

type StaffRepository interface {
	Insert(ctx context.Context, u *Staff) error
	UpdateById(ctx context.Context, u *Staff) error
	GetById(ctx context.Context, id uint64) (*Staff, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *Staff) error
}
