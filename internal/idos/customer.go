package idos

import (
	"time"
	
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
)

type LiteCustomerFilterIDO struct {
	TenantId  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder null.String `json:"sort_order"`
	SortField null.String `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"last_seen_id"`
	Limit     uint64      `json:"limit"`
}

type LiteCustomerListResponseIDO struct {
	NextId  uint64                 `json:"next_id,omitempty"`
	Count   uint64                 `json:"count"`
	Results []*models.LiteCustomer `json:"results"`
}

func NewLiteCustomerListResponseIDO(arr []*models.LiteCustomer, count uint64) *LiteCustomerListResponseIDO {
	// Calculate next id.
	var nextId uint64
	if len(arr) > 0 {
		lastRecord := arr[len(arr)-1]
		nextId = lastRecord.Id
	}

	res := &LiteCustomerListResponseIDO{ // Return through HTTP.
		Count:   count,
		Results: arr,
		NextId:  nextId,
	}

	return res
}


type CustomerIDO struct {
	// -- customer.py
	Id                      uint64    `json:"id"`
	Uuid                    string    `json:"uuid"`
	TenantId                uint64    `json:"tenant_id"`
	UserId                  uint64    `json:"user_id"`
	TypeOf                  int8      `json:"type_of"`
	IndexedText             string    `json:"indexed_text"`
	IsOkToEmail             bool      `json:"is_ok_to_email"`
	IsOkToText              bool      `json:"is_ok_to_text"`
	IsBusiness              bool      `json:"is_business"`
	IsSenior                bool      `json:"is_senior"`
	IsSupport               bool      `json:"is_support"`
	JobInfoRead             string    `json:"job_info_read"`
	HowHearId               uint64    `json:"how_hear_id"`
	HowHearOld              int8      `json:"how_hear_old"`
	HowHearOther            string    `json:"how_hear_other"`
	State                   int8      `json:"state"`
	DeactivationReason      int8      `json:"deactivation_reason"`
	DeactivationReasonOther string    `json:"deactivation_reason_other"`
	CreatedTime             time.Time `json:"created_time"`
	CreatedById             uint64    `json:"created_by_id"`
	CreatedFromIP           string    `json:"created_from_ip"`
	LastModifiedTime        time.Time `json:"last_modified_time"`
	LastModifiedById        uint64    `json:"last_modified_by_id"`
	LastModifiedFromIP      string    `json:"last_modified_from_ip"`
	OrganizationName        string    `json:"organization_name"`
	OrganizationTypeOf      int8      `json:"organization_type_of"`
	OldId                   uint64    `json:"old_id"`

	// -- abstract_postal_address.py
	AddressCountry      string `json:"address_country"`
	AddressRegion       string `json:"address_region"`
	AddressLocality     string `json:"address_locality"`
	PostOfficeBoxNumber string `json:"post_office_box_number"`
	PostalCode          string `json:"postal_code"`
	StreetAddress       string `json:"street_address"`
	StreetAddressExtra  string `json:"street_address_extra"`

	// -- abstract_person.py
	GivenName   string       `json:"given_name"`
	MiddleName  string       `json:"middle_name"`
	LastName    string       `json:"last_name"`
	Name        string       `json:"name,omitempty"`
	LexicalName string       `json:"lexical_name,omitempty"`
	Birthdate   null.Time `json:"birthdate"`
	JoinDate    null.Time `json:"join_date"`
	Nationality string       `json:"nationality"`
	Gender      string       `json:"gender"`
	TaxId       string       `json:"tax_id"`

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

func NewCustomerIDO(m *models.Customer) *CustomerIDO {
	return &CustomerIDO{
		// -- customer.py
		Id:                      m.Id,
		Uuid:                    m.Uuid,
		TenantId:                m.TenantId,
		UserId:                  m.UserId,
		TypeOf:                  m.TypeOf,
		IndexedText:             m.IndexedText,
		IsOkToEmail:             m.IsOkToEmail,
		IsOkToText:              m.IsOkToText,
		IsBusiness:              m.IsBusiness,
		IsSenior:                m.IsSenior,
		IsSupport:               m.IsSupport,
		JobInfoRead:             m.JobInfoRead,
		HowHearId:               m.HowHearId,
		HowHearOld:              m.HowHearOld,
		HowHearOther:            m.HowHearOther,
		State:                   m.State,
		DeactivationReason:      m.DeactivationReason,
		DeactivationReasonOther: m.DeactivationReasonOther,
		CreatedTime:             m.CreatedTime,
		CreatedById:             m.CreatedById,
		CreatedFromIP:           m.CreatedFromIP,
		LastModifiedTime:        m.LastModifiedTime,
		LastModifiedById:        m.LastModifiedById,
		LastModifiedFromIP:      m.LastModifiedFromIP,
		OrganizationName:        m.OrganizationName,
		OrganizationTypeOf:      m.OrganizationTypeOf,

		// -- abstract_postal_address.py
		AddressCountry:          m.AddressCountry,
		AddressRegion:           m.AddressRegion,
		AddressLocality:         m.AddressLocality,
		PostOfficeBoxNumber:     m.PostOfficeBoxNumber,
		PostalCode:              m.PostalCode,
		StreetAddress:           m.StreetAddress,
		StreetAddressExtra:      m.StreetAddressExtra,

		// -- abstract_person.py
		GivenName:               m.GivenName,
		MiddleName:              m.MiddleName,
		LastName:                m.LastName,
		Name:                    m.Name,
		LexicalName:             m.LexicalName,
		Birthdate:               m.Birthdate,
		JoinDate:                m.JoinDate,
		Nationality:             m.Nationality,
		Gender:                  m.Gender,
		TaxId:                   m.TaxId,

		// -- abstract_geo_coorindate.py
		Elevation:               m.Elevation,
		Latitude:                m.Latitude,
		Longitude:               m.Longitude,

		// -- abstract_contact_point.py
		AreaServed:              m.AreaServed,
		AvailableLanguage:       m.AvailableLanguage,
		ContactType:             m.ContactType,
		Email:                   m.Email,
		FaxNumber:               m.FaxNumber,
		Telephone:               m.Telephone,
		TelephoneTypeOf:         m.TelephoneTypeOf,
		TelephoneExtension:      m.TelephoneExtension,
		OtherTelephone:          m.OtherTelephone,
		OtherTelephoneExtension: m.OtherTelephoneExtension,
		OtherTelephoneTypeOf:    m.OtherTelephoneTypeOf,

		// Name     string  `json:"name"`
		// Url     string  `json:"url"`
		// OrganizationId sql.NullInt64  `json:"organization_id"`
		// OwnerId sql.NullInt64  `json:"owner_id"`
		// IsBlacklisted bool `json:"is_blacklisted"`
		// AvatarImageId sql.NullInt64  `json:"avatar_image_id"`
	}
}
