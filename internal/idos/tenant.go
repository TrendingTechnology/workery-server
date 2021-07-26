package idos

import (
	"time"
	// null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
)

type TenantIDO struct {
	Id                      uint64    `json:"id"`
	Uuid                    string    `json:"uuid"`
	SchemaName              string    `json:"schema_name"`
	AlternateName           string    `json:"alternate_name"`
	Description             string    `json:"description"`
	Name                    string    `json:"name"`
	Url                     string    `json:"url"`
	State                   int8      `json:"state"`
	Timezone                string    `json:"timestamp"`
	CreatedTime             time.Time `json:"created_time"`
	ModifiedTime            time.Time `json:"modified_time"`
	AddressCountry          string    `json:"address_country"`
	AddressRegion           string    `json:"address_region"`
	AddressLocality         string    `json:"address_locality"`
	PostOfficeBoxNumber     string    `json:"post_office_box_number"`
	PostalCode              string    `json:"postal_code"`
	StreetAddress           string    `json:"street_address"`
	StreetAddressExtra      string    `json:"street_address_extra"`
	Elevation               float64   `json:"elevation"`
	Latitude                float64   `json:"latitude"`
	Longitude               float64   `json:"longitude"`
	AreaServed              string    `json:"area_served"`
	AvailableLanguage       string    `json:"available_language"`
	ContactType             string    `json:"contact_type"`
	Email                   string    `json:"email"`
	FaxNumber               string    `json:"fax_number"`
	Telephone               string    `json:"telephone"`
	TelephoneTypeOf         int8      `json:"telephone_type_of"`
	TelephoneExtension      string    `json:"telephone_extension"`
	OtherTelephone          string    `json:"other_telephone"`
	OtherTelephoneExtension string    `json:"other_telephone_extension"`
	OtherTelephoneTypeOf    int8      `json:"other_telephone_type_of"`
}

func NewTenantIDO(m *models.Tenant) *TenantIDO {
	return &TenantIDO{
		Id:                      m.Id,
		Uuid:                    m.Uuid,
		SchemaName:              m.SchemaName,
		AlternateName:           m.AlternateName,
		Description:             m.Description,
		Name:                    m.Name,
		Url:                     m.Url,
		State:                   m.State,
		Timezone:                m.Timezone,
		CreatedTime:             m.CreatedTime,
		ModifiedTime:            m.ModifiedTime,
		AddressCountry:          m.AddressCountry,
		AddressRegion:           m.AddressRegion,
		AddressLocality:         m.AddressLocality,
		PostOfficeBoxNumber:     m.PostOfficeBoxNumber,
		PostalCode:              m.PostalCode,
		StreetAddress:           m.StreetAddress,
		StreetAddressExtra:      m.StreetAddressExtra,
		Elevation:               m.Elevation,
		Latitude:                m.Latitude,
		Longitude:               m.Longitude,
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
	}
}
