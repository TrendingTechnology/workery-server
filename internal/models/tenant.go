package models

import (
	"context"
	"time"
)

type Tenant struct {
	Id                      uint64    `json:"id"`
	Uuid                    string    `json:"uuid"`
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

type TenantRepository interface {
	Insert(ctx context.Context, u *Tenant) error
	Update(ctx context.Context, u *Tenant) error
	GetById(ctx context.Context, id uint64) (*Tenant, error)
	GetByName(ctx context.Context, name string) (*Tenant, error)
	CheckIfExistsById(ctx context.Context, id uint64) (bool, error)
	CheckIfExistsByName(ctx context.Context, name string) (bool, error)
	InsertOrUpdate(ctx context.Context, u *Tenant) error
}
