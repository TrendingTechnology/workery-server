package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	// "strconv"

	// "github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	customerETLSchemaName string
)

func init() {
	customerETLCmd.Flags().StringVarP(&customerETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	customerETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(customerETLCmd)
}

var customerETLCmd = &cobra.Command{
	Use:   "etl_customer",
	Short: "Import the customer data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportCustomer()
	},
}

func doRunImportCustomer() {
	// Load up our NEW database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our OLD database.
	oldDBHost := os.Getenv("WORKERY_OLD_DB_HOST")
	oldDBPort := os.Getenv("WORKERY_OLD_DB_PORT")
	oldDBUser := os.Getenv("WORKERY_OLD_DB_USER")
	oldDBPassword := os.Getenv("WORKERY_OLD_DB_PASSWORD")
	oldDBName := os.Getenv("WORKERY_OLD_DB_NAME")
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, customerETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	ur := repositories.NewUserRepo(db)
	om := repositories.NewCustomerRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, customerETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runCustomerETL(ctx, tenant.Id, ur, om, oldDb)
}

type OldUCustomer struct {
	Created                  time.Time       `json:"created"`
	LastModified             time.Time       `json:"last_modified"`
	AlternateName            sql.NullString  `json:"alternate_name"`
	Description              sql.NullString  `json:"description"`
	Name                     sql.NullString  `json:"name"`
	Url                      sql.NullString  `json:"url"`
	AreaServed               sql.NullString  `json:"area_served"`
	AvailableLanguage        sql.NullString  `json:"available_language"`
	ContactType              sql.NullString  `json:"contact_type"`
	Email                    sql.NullString  `json:"email"`
	FaxNumber                sql.NullString  `json:"fax_number"`
	ProductSupported         sql.NullString  `json:"product_supported"`
	Telephone                sql.NullString  `json:"telephone"`
	TelephoneTypeOf          int8            `json:"telephone_type_of"`
	TelephoneExtension       sql.NullString  `json:"telephone_extension"`
	OtherTelephone           sql.NullString  `json:"other_telephone"`
	OtherTelephoneExtension  sql.NullString  `json:"other_telephone_extension"`
	OtherTelephoneTypeOf     int8            `json:"other_telephone_type_of"`
	AddressCountry           string          `json:"address_country"`
	AddressLocality          string          `json:"address_locality"`
	AddressRegion            string          `json:"address_region"`
	PostOfficeBoxNumber      sql.NullString  `json:"post_office_box_number"`
	PostalCode               sql.NullString  `json:"postal_code"`
	StreetAddress            string          `json:"street_address"`
	StreetAddressExtra       sql.NullString  `json:"street_address_extra"`
	Elevation                sql.NullFloat64 `json:"elevation"`
	Latitude                 sql.NullFloat64 `json:"latitude"`
	Longitude                sql.NullFloat64 `json:"longitude"`
	GivenName                sql.NullString  `json:"given_name"`
	MiddleName               sql.NullString  `json:"middle_name"`
	LastName                 sql.NullString  `json:"last_name"`
	Birthdate                sql.NullTime    `json:"birthdate"`
	JoinDate                 sql.NullTime    `json:"join_date"`
	Nationality              sql.NullString  `json:"nationality"`
	Gender                   sql.NullString  `json:"gender"`
	TaxId                    sql.NullString  `json:"tax_id"`
	Id                       uint64          `json:"id"`
	IndexedText              sql.NullString  `json:"indexed_text"`
	TypeOf                   int8            `json:"type_of"`
	IsOkToEmail              bool            `json:"is_ok_to_email"`
	IsOkToText               bool            `json:"is_ok_to_text"`
	IsBusiness               bool            `json:"is_business"`
	IsSenior                 bool            `json:"is_senior"`
	IsSupport                bool            `json:"is_support"`
	JobInfoRead              sql.NullString  `json:"job_info_read"`
	CreatedFrom              sql.NullString  `json:"created_from"`
	CreatedFromIsPublic      bool            `json:"created_from_is_public"`
	LastModifiedFrom         sql.NullString  `json:"last_modified_from"`
	LastModifiedFromIsPublic bool            `json:"last_modified_from_is_public"`
	IsArchived               bool            `json:"is_archived"`
	CreatedById              sql.NullInt64   `json:"created_by_id"`
	LastModifiedById         sql.NullInt64   `json:"last_modified_by_id"`
	OrganizationId           sql.NullInt64   `json:"organization_id"`
	OwnerId                  sql.NullInt64   `json:"owner_id"`
	HowHearOther             string          `json:"how_hear_other"`
	IsBlacklisted            bool            `json:"is_blacklisted"`
	DeactivationReason       int8            `json:"deactivation_reason"`
	DeactivationReasonOther  string          `json:"deactivation_reason_other"`
	State                    string          `json:"state"`
	HowHearId                sql.NullInt64   `json:"how_hear_id"`
	HowHearOld               int8            `json:"how_hear_old"`
	OrganizationName         sql.NullString  `json:"organization_name"`
	OrganizationTypeOf       int8            `json:"organization_type_of"`
	AvatarImageId            sql.NullInt64   `json:"avatar_image_id"`
}

func ListAllCustomers(db *sql.DB) ([]*OldUCustomer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, created, last_modified, alternate_name, description, name, url,
		area_served, available_language, contact_type, email, fax_number,
		product_supported, telephone, telephone_type_of, telephone_extension,
		other_telephone, other_telephone_extension, other_telephone_type_of,
		address_country, address_locality, address_region, post_office_box_number,
		postal_code, street_address, street_address_extra, elevation, latitude,
		longitude, given_name, middle_name, last_name, birthdate, join_date,
		nationality, gender, tax_id, id, indexed_text, type_of, is_ok_to_email,
		is_ok_to_text, is_business, is_senior, is_support, job_info_read,
		created_from, created_from_is_public, last_modified_from,
		last_modified_from_is_public, is_archived, created_by_id, last_modified_by_id,
		organization_id, owner_id, how_hear_other, is_blacklisted, deactivation_reason,
		deactivation_reason_other, state, how_hear_id, how_hear_old, organization_name,
		organization_type_of, avatar_image_id
	FROM
	    workery_customers
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUCustomer
	defer rows.Close()
	for rows.Next() {
		m := new(OldUCustomer)
		err = rows.Scan(
			&m.Id, &m.Created, &m.LastModified, &m.AlternateName, &m.Description, &m.Name, &m.Url,
			&m.AreaServed, &m.AvailableLanguage, &m.ContactType, &m.Email, &m.FaxNumber,
			&m.ProductSupported, &m.Telephone, &m.TelephoneTypeOf, &m.TelephoneExtension,
			&m.OtherTelephone, &m.OtherTelephoneExtension, &m.OtherTelephoneTypeOf,
			&m.AddressCountry, &m.AddressLocality, &m.AddressRegion, &m.PostOfficeBoxNumber,
			&m.PostalCode, &m.StreetAddress, &m.StreetAddressExtra, &m.Elevation,
			&m.Latitude, &m.Longitude, &m.GivenName, &m.MiddleName, &m.LastName,
			&m.Birthdate, &m.JoinDate, &m.Nationality, &m.Gender, &m.TaxId, &m.Id, &m.IndexedText, &m.TypeOf,
			&m.IsOkToEmail, &m.IsOkToText, &m.IsBusiness, &m.IsSenior, &m.IsSupport,
			&m.JobInfoRead, &m.CreatedFrom, &m.CreatedFromIsPublic, &m.LastModifiedFrom,
			&m.LastModifiedFromIsPublic, &m.IsArchived, &m.CreatedById, &m.LastModifiedById,
			&m.OrganizationId, &m.OwnerId, &m.HowHearOther, &m.IsBlacklisted, &m.DeactivationReason,
			&m.DeactivationReasonOther, &m.State, &m.HowHearId, &m.HowHearOld, &m.OrganizationName,
			&m.OrganizationTypeOf, &m.AvatarImageId,
		)
		if err != nil {
			log.Fatal("(AA)", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("(BB)", err)
	}
	return arr, err
}

func runCustomerETL(ctx context.Context, tenantId uint64, ur *repositories.UserRepo, omr *repositories.CustomerRepo, oldDb *sql.DB) {
	customers, err := ListAllCustomers(oldDb)
	if err != nil {
		log.Fatal("(0000)", err)
	}
	for _, om := range customers {
		insertCustomerETL(ctx, tenantId, ur, omr, om)
	}
}

func insertCustomerETL(ctx context.Context, tid uint64, ur *repositories.UserRepo, omr *repositories.CustomerRepo, om *OldUCustomer) {
	// var state int8 = 1
	// if om.IsArchived == true {
	// 	state = 0
	// }

    userId := uint64(om.OwnerId.Int64)
	var err error
	var user *models.User
	if om.OwnerId.Valid {
		log.Println(om.Id)
		user, err = ur.GetByOldId(ctx, om.Id)
		if err != nil {
			return
		}
		userId = user.Id
	}
	fmt.Println("Imported UserID#", userId)

	// m := &models.Customer{
	// 	OldId:    om.Id,
	// 	Uuid:     uuid.NewString(),
	// 	TenantId: tid,
	// 	UserId:   userId,
	// 	TypeOf:   om.TypeOf,
	// 	State:    state,
	// }
	//
	// fmt.Println("Imported UserID#", m.UserId)
	//
	// // err := om.Insert(ctx, m)
	// // if err != nil {
	// // 	log.Panic(err)
	// // }
	// // fmt.Println("Imported ID#", om.Id)
}
