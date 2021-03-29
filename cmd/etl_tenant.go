package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

func init() {
	rootCmd.AddCommand(tenantETLCmd)
}

var tenantETLCmd = &cobra.Command{
	Use:   "etl_tenant",
	Short: "Import the tenant data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportTenant()
	},
}

func doRunImportTenant() {
	// Load up our new database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewTenantRepo(db)

	// Load up our old database.
	oldDBHost := os.Getenv("WORKERY_OLD_DB_HOST")
	oldDBPort := os.Getenv("WORKERY_OLD_DB_PORT")
	oldDBUser := os.Getenv("WORKERY_OLD_DB_USER")
	oldDBPassword := os.Getenv("WORKERY_OLD_DB_PASSWORD")
	oldDBName := os.Getenv("WORKERY_OLD_DB_NAME")
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Begin the operation.
	runTenantETL(r, oldDb)
}

type OldTenant struct {
	Id                      uint64          `json:"id"`
	SchemaName              string          `json:"schema_name"`
	Created                 time.Time       `json:"created"`
	LastModified            time.Time       `json:"last_modified"`
	AlternateName           string          `json:"alternate_name"`
	Description             string          `json:"description"`
	Name                    string          `json:"name"`
	Url                     sql.NullString  `json:"url"`
	AreaServed              sql.NullString  `json:"area_served"`
	AvailableLanguage       sql.NullString  `json:"available_language"`
	ContactType             sql.NullString  `json:"contact_type"`
	Email                   sql.NullString  `json:"email"`
	FaxNumber               sql.NullString  `json:"fax_number"`
	Telephone               sql.NullString  `json:"telephone"`
	TelephoneTypeOf         int8            `json:"telephone_type_of"`
	TelephoneExtension      sql.NullString  `json:"telephone_extension"`
	OtherTelephone          sql.NullString  `json:"other_telephone"`
	OtherTelephoneExtension sql.NullString  `json:"other_telephone_extension"`
	OtherTelephoneTypeOf    int8            `json:"other_telephone_type_of"`
	AddressCountry          string          `json:"address_country"`
	AddressRegion           string          `json:"address_region"`
	AddressLocality         string          `json:"address_locality"`
	PostOfficeBoxNumber     string          `json:"post_office_box_number"`
	PostalCode              string          `json:"postal_code"`
	StreetAddress           string          `json:"street_address"`
	StreetAddressExtra      string          `json:"street_address_extra"`
	Elevation               sql.NullFloat64 `json:"elevation"`
	Latitude                sql.NullFloat64 `json:"latitude"`
	Longitude               sql.NullFloat64 `json:"longitude"`
	TimezoneName            string          `json:"timestamp_name"`
	IsArchived              bool            `json:"is_archived"`
}

// Function returns a paginated list of all type element items.
func ListAllTenants(db *sql.DB) ([]*OldTenant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, schema_name, created, last_modified, alternate_name, description,
		name, url, area_served, available_language, contact_type, email,
		fax_number, telephone, telephone_type_of, telephone_extension,
		other_telephone, other_telephone_extension, other_telephone_type_of,
		address_country, address_region, address_locality, post_office_box_number,
		postal_code, street_address, street_address_extra, elevation, latitude,
		longitude, timezone_name, is_archived
	FROM
	    workery_franchises
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldTenant
	defer rows.Close()
	for rows.Next() {
		m := new(OldTenant)
		err = rows.Scan(
			&m.Id,
			&m.SchemaName,
			&m.Created,
			&m.LastModified,
			&m.AlternateName,
			&m.Description,
			&m.Name,
			&m.Url,
			&m.AreaServed,
			&m.AvailableLanguage,
			&m.ContactType,
			&m.Email,
			&m.FaxNumber,
			&m.Telephone,
			&m.TelephoneTypeOf,
			&m.TelephoneExtension,
			&m.OtherTelephone,
			&m.OtherTelephoneExtension,
			&m.OtherTelephoneTypeOf,
			&m.AddressCountry,
			&m.AddressRegion,
			&m.AddressLocality,
			&m.PostOfficeBoxNumber,
			&m.PostalCode,
			&m.StreetAddress,
			&m.StreetAddressExtra,
			&m.Elevation,
			&m.Latitude,
			&m.Longitude,
			&m.TimezoneName,
			&m.IsArchived,
		)
		if err != nil {
			panic(err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func runTenantETL(r *repositories.TenantRepo, oldDb *sql.DB) {
	tenants, err := ListAllTenants(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range tenants {
		runTenantInsert(v, r)
	}
}

func runTenantInsert(ot *OldTenant, r *repositories.TenantRepo) {
	m := &models.Tenant{
		OldId:              ot.Id,
		Uuid:               uuid.NewString(),
		AlternateName:      ot.AlternateName,
		Description:        ot.Description,
		Name:               ot.Name,
		Url:                ot.Url.String,
		State:              1,
		Timezone:           "America/Toronto",
		CreatedTime:        ot.Created,
		ModifiedTime:       ot.LastModified,
		AddressCountry:     ot.AddressCountry,
		AddressRegion:      ot.AddressRegion,
		AddressLocality:    ot.AddressLocality,
		PostalCode:         ot.PostalCode,
		StreetAddress:      ot.StreetAddress,
		StreetAddressExtra: ot.StreetAddressExtra,
		SchemaName:         ot.SchemaName,
	}
	ctx := context.Background()
	err := r.InsertOrUpdateBySchemaName(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", ot.Id)
}
