package models

import (
	"context"
	"time"

	null "gopkg.in/guregu/null.v4"
)

// State
//---------------------
// 1 = Active
// 0 = Inactive

type Partner struct {
	// -- partner.py
	Id                            uint64      `json:"id"`
	Uuid                          string      `json:"uuid"`
	TenantId                      uint64      `json:"tenant_id"`
	UserId                        uint64      `json:"user_id"`
	TypeOf                        int8        `json:"type_of"`
	OrganizationName              string      `json:"organization_name"`
	OrganizationTypeOf            int8        `json:"organization_type_of"`
	Business                      string      `json:"business"`
	IndexedText                   string      `json:"indexed_text"`
	IsOkToEmail                   bool        `json:"is_ok_to_email"`
	IsOkToText                    bool        `json:"is_ok_to_text"`
	HourlySalaryDesired           int8        `json:"hourly_salary_desired"`
	LimitSpecial                  string      `json:"limit_special"`
	DuesDate                      time.Time   `json:"dues_date"`
	CommercialInsuranceExpiryDate time.Time   `json:"commercial_insurance_expiry_date"`
	AutoInsuranceExpiryDate       time.Time   `json:"auto_insurance_expiry_date"`
	WsibNumber                    string      `json:"wsib_number"`
	WsibInsuranceDate             time.Time   `json:"wsib_insurance_date"`
	PoliceCheck                   time.Time   `json:"police_check"`
	DriversLicenseClass           string      `json:"drivers_license_class"`
	HowHearOld                    int8        `json:"how_hear_old"`
	HowHearId                     uint64      `json:"how_hear_id"`
	HowHearOther                  string      `json:"how_hear_other"`
	HowHearText                   string      `json:"how_hear_text"` // Referenced value from `HowHearAboutUsItem`.
	State                         int8        `json:"state"`
	DeactivationReason            int8        `json:"deactivation_reason"`
	DeactivationReasonOther       string      `json:"deactivation_reason_other"`
	CreatedTime                   time.Time   `json:"created_time"`
	CreatedById                   null.Int    `json:"created_by_id"`
	CreatedByName                 null.String `json:"created_by_name"`
	CreatedFromIP                 string      `json:"created_from_ip"`
	LastModifiedTime              time.Time   `json:"last_modified_time"`
	LastModifiedById              null.Int    `json:"last_modified_by_id"`
	LastModifiedByName            null.String `json:"last_modified_by_name"`
	LastModifiedFromIP            string      `json:"last_modified_from_ip"`
	Score                         float64     `json:"score"`
	OldId                         uint64      `json:"old_id"`
	ServiceFeeId                  uint64      `json:"service_fee_id"`

	// -- abstract_postal_address.py
	AddressCountry               string `json:"address_country"`
	AddressRegion                string `json:"address_region"`
	AddressLocality              string `json:"address_locality"`
	PostOfficeBoxNumber          string `json:"post_office_box_number"`
	PostalCode                   string `json:"postal_code"`
	StreetAddress                string `json:"street_address"`
	StreetAddressExtra           string `json:"street_address_extra"`
	FullAddressWithoutPostalCode string `json:"full_address_without_postal_code,omitempty"` // Compiled value
	FullAddressWithPostalCode    string `json:"full_address_with_postal_code,omitempty"`    // Compiled value
	FullAddressUrl               string `json:"full_address_url,omitempty"`                 // Compiled value

	// -- abstract_person.py
	GivenName   string    `json:"given_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	Name        string    `json:"name,omitempty"`
	LexicalName string    `json:"lexical_name,omitempty"`
	Birthdate   null.Time `json:"birthdate"`
	JoinDate    null.Time `json:"join_date"`
	Nationality string    `json:"nationality"`
	Gender      string    `json:"gender"`
	TaxId       string    `json:"tax_id"`

	// -- abstract_geo_coorindate.py
	Elevation float64 `json:"elevation"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	// -- abstract_contact_point.py
	AreaServed              string `json:"area_served"`
	AvailableLanguage       string `json:"available_language"`
	ContactType             string `json:"contact_type"`
	Email                   string `json:"email"`
	FaxNumber               string `json:"fax_number"`
	Telephone               string `json:"telephone"`
	TelephoneTypeOf         int8   `json:"telephone_type_of"`
	TelephoneExtension      string `json:"telephone_extension"`
	OtherTelephone          string `json:"other_telephone"`
	OtherTelephoneExtension string `json:"other_telephone_extension"`
	OtherTelephoneTypeOf    int8   `json:"other_telephone_type_of"`

	// Name     string  `json:"name"`
	// Url     string  `json:"url"`
	// OrganizationId sql.NullInt64  `json:"organization_id"`
	// OwnerId sql.NullInt64  `json:"owner_id"`
	// IsBlacklisted bool `json:"is_blacklisted"`
	// AvatarImageId sql.NullInt64  `json:"avatar_image_id"`
}

type PartnerRepository interface {
	Insert(ctx context.Context, u *Partner) error
	UpdateById(ctx context.Context, u *Partner) error
	GetById(ctx context.Context, id uint64) (*Partner, error)
	GetIdByOldId(ctx context.Context, tid uint64, oid uint64) (uint64, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	InsertOrUpdateById(ctx context.Context, u *Partner) error
}
