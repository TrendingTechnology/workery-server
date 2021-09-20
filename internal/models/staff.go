package models

import (
	"context"
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
	CreatedTime                          time.Time   `json:"created_time"`
	LastModifiedTime                     time.Time   `json:"last_modified_time"`
	AvailableLanguage                    null.String `json:"available_language"`
	ContactType                          null.String `json:"contact_type"`
	Email                                null.String `json:"email"`
	FaxNumber                            null.String `json:"fax_number"`
	Telephone                            null.String `json:"telephone"`
	TelephoneTypeOf                      int8        `json:"telephone_type_of"`
	TelephoneExtension                   null.String `json:"telephone_extension"`
	OtherTelephone                       null.String `json:"other_telephone"`
	OtherTelephoneExtension              null.String `json:"other_telephone_extension"`
	OtherTelephoneTypeOf                 int8        `json:"other_telephone_type_of"`
	AddressCountry                       string      `json:"address_country"`
	AddressLocality                      string      `json:"address_locality"`
	AddressRegion                        string      `json:"address_region"`
	PostOfficeBoxNumber                  null.String `json:"post_office_box_number"`
	PostalCode                           null.String `json:"postal_code"`
	StreetAddress                        string      `json:"street_address"`
	StreetAddressExtra                   null.String `json:"street_address_extra"`
	FullAddressWithoutPostalCode         string      `json:"full_address_without_postal_code,omitempty"` // Compiled value
	FullAddressWithPostalCode            string      `json:"full_address_with_postal_code,omitempty"`    // Compiled value
	FullAddressUrl                       string      `json:"full_address_url,omitempty"`                 // Compiled value
	Elevation                            null.Float  `json:"elevation"`
	Latitude                             null.Float  `json:"latitude"`
	Longitude                            null.Float  `json:"longitude"`
	GivenName                            null.String `json:"given_name"`
	MiddleName                           null.String `json:"middle_name"`
	LastName                             null.String `json:"last_name"`
	Name                                 string      `json:"name,omitempty"`
	LexicalName                          string      `json:"lexical_name,omitempty"`
	Birthdate                            null.Time   `json:"birthdate"`
	JoinDate                             null.Time   `json:"join_date"`
	Nationality                          null.String `json:"nationality"`
	Gender                               null.String `json:"gender"`
	TaxId                                null.String `json:"tax_id"`
	TenantId                             uint64      `json:"tenant_id"`
	Id                                   uint64      `json:"id"`
	Uuid                                 string      `json:"uuid"`
	IndexedText                          null.String `json:"indexed_text"`
	CreatedFromIP                        null.String `json:"created_from_ip"`
	LastModifiedFromIP                   null.String `json:"last_modified_from_ip"`
	State                                int8        `json:"state"`
	CreatedById                          null.Int    `json:"created_by_id"`
	CreatedByName                        null.String `json:"created_by_name"`
	LastModifiedById                     null.Int    `json:"last_modified_by_id"`
	LastModifiedByName                   null.String `json:"last_modified_by_name"`
	UserId                               uint64      `json:"user_id"`
	HowHearId                            null.Int    `json:"how_hear_id"`
	HowHearOther                         null.String `json:"how_hear_other"`
	HowHearText                          string      `json:"how_hear_text"` // Referenced value from `HowHearAboutUsItem`.
	AvatarImageId                        null.Int    `json:"avatar_image_id"`
	PersonalEmail                        null.String `json:"personal_email"`
	EmergencyContactAlternativeTelephone null.String `json:"emergency_contact_alternative_telephone"`
	EmergencyContactName                 null.String `json:"emergency_contact_name"`
	EmergencyContactRelationship         null.String `json:"emergency_contact_relationship"`
	EmergencyContactTelephone            null.String `json:"emergency_contact_telephone"`
	PoliceCheck                          null.Time   `json:"police_check"`
	OldId                                uint64      `json:"old_id"`
	IsOkToEmail                          bool        `json:"is_ok_to_email"`
	IsOkToText                           bool        `json:"is_ok_to_text"`
}

type StaffRepository interface {
	Insert(ctx context.Context, u *Staff) error
	UpdateById(ctx context.Context, u *Staff) error
	GetById(ctx context.Context, id uint64) (*Staff, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *Staff) error
}
