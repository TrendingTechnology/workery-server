package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	null "gopkg.in/guregu/null.v4"

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
	r := repositories.NewHowHearAboutUsItemRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, customerETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runCustomerETL(ctx, tenant.Id, ur, om, r, oldDb)
}

type OldUCustomer struct {
	Created                  time.Time   `json:"created"`
	LastModified             time.Time   `json:"last_modified"`
	AlternateName            null.String `json:"alternate_name"`
	Description              null.String `json:"description"`
	Name                     null.String `json:"name"`
	Url                      null.String `json:"url"`
	AreaServed               null.String `json:"area_served"`
	AvailableLanguage        null.String `json:"available_language"`
	ContactType              null.String `json:"contact_type"`
	Email                    null.String `json:"email"`
	FaxNumber                null.String `json:"fax_number"`
	ProductSupported         null.String `json:"product_supported"`
	Telephone                null.String `json:"telephone"`
	TelephoneTypeOf          int8        `json:"telephone_type_of"`
	TelephoneExtension       null.String `json:"telephone_extension"`
	OtherTelephone           null.String `json:"other_telephone"`
	OtherTelephoneExtension  null.String `json:"other_telephone_extension"`
	OtherTelephoneTypeOf     int8        `json:"other_telephone_type_of"`
	AddressCountry           string      `json:"address_country"`
	AddressLocality          string      `json:"address_locality"`
	AddressRegion            string      `json:"address_region"`
	PostOfficeBoxNumber      null.String `json:"post_office_box_number"`
	PostalCode               null.String `json:"postal_code"`
	StreetAddress            string      `json:"street_address"`
	StreetAddressExtra       null.String `json:"street_address_extra"`
	Elevation                null.Float  `json:"elevation"`
	Latitude                 null.Float  `json:"latitude"`
	Longitude                null.Float  `json:"longitude"`
	GivenName                null.String `json:"given_name"`
	MiddleName               null.String `json:"middle_name"`
	LastName                 null.String `json:"last_name"`
	Birthdate                null.Time   `json:"birthdate"`
	JoinDate                 null.Time   `json:"join_date"`
	Nationality              null.String `json:"nationality"`
	Gender                   null.String `json:"gender"`
	TaxId                    null.String `json:"tax_id"`
	Id                       uint64      `json:"id"`
	IndexedText              null.String `json:"indexed_text"`
	TypeOf                   int8        `json:"type_of"`
	IsOkToEmail              bool        `json:"is_ok_to_email"`
	IsOkToText               bool        `json:"is_ok_to_text"`
	IsBusiness               bool        `json:"is_business"`
	IsSenior                 bool        `json:"is_senior"`
	IsSupport                bool        `json:"is_support"`
	JobInfoRead              null.String `json:"job_info_read"`
	CreatedFrom              null.String `json:"created_from"`
	CreatedFromIsPublic      bool        `json:"created_from_is_public"`
	LastModifiedFrom         null.String `json:"last_modified_from"`
	LastModifiedFromIsPublic bool        `json:"last_modified_from_is_public"`
	IsArchived               bool        `json:"is_archived"`
	CreatedById              null.Int    `json:"created_by_id"`
	LastModifiedById         null.Int    `json:"last_modified_by_id"`
	OrganizationId           null.Int    `json:"organization_id"`
	OwnerId                  null.Int    `json:"owner_id"`
	HowHearOther             string      `json:"how_hear_other"`
	IsBlacklisted            bool        `json:"is_blacklisted"`
	DeactivationReason       int8        `json:"deactivation_reason"`
	DeactivationReasonOther  string      `json:"deactivation_reason_other"`
	State                    string      `json:"state"`
	HowHearId                null.Int    `json:"how_hear_id"`
	HowHearOld               int8        `json:"how_hear_old"`
	OrganizationName         null.String `json:"organization_name"`
	OrganizationTypeOf       int8        `json:"organization_type_of"`
	AvatarImageId            null.Int    `json:"avatar_image_id"`
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
	ORDER BY
	    id
	ASC
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

func runCustomerETL(
	ctx context.Context,
	tenantId uint64,
	ur *repositories.UserRepo,
	omr *repositories.CustomerRepo,
	r *repositories.HowHearAboutUsItemRepo,
	oldDb *sql.DB,
) {
	customers, err := ListAllCustomers(oldDb)
	if err != nil {
		log.Fatal("(0000)", err)
	}
	for _, om := range customers {
		insertCustomerETL(ctx, tenantId, ur, omr, r, om)
	}
}

func insertCustomerETL(
	ctx context.Context,
	tid uint64,
	ur *repositories.UserRepo,
	omr *repositories.CustomerRepo,
	r *repositories.HowHearAboutUsItemRepo,
	om *OldUCustomer,
) {
	//
	// Set the `state`.
	//

	var state int8 = 0
	if om.State == "active" {
		state = models.CustomerActiveState
	}

	// Variable used to keep the ID of the user record in our database.
	userId := uint64(om.OwnerId.Int64)

	//
	// Generate our full name / lexical full name.
	//

	var name string
	var lexicalName string
	if om.MiddleName.Valid {
		name = om.GivenName.String + " " + om.MiddleName.String + " " + om.LastName.String
		lexicalName = om.LastName.String + ", " + om.MiddleName.String + ", " + om.GivenName.String
	} else {
		name = om.GivenName.String + " " + om.LastName.String
		lexicalName = om.LastName.String + ", " + om.GivenName.String
	}
	lexicalName = strings.Replace(lexicalName, ", ,", ",", 0)
	lexicalName = strings.Replace(lexicalName, "  ", " ", 0)
	lexicalName = strings.Replace(lexicalName, ", , ", ", ", 0)
	name = strings.Replace(name, "  ", " ", 0)

	// CASE 1: User record exists in our database.
	if om.OwnerId.Valid {
		// log.Println(om.Id)
		user, err := ur.GetByOldId(ctx, userId)
		if err != nil {
			log.Fatal("(A)", err)
		}
		if user == nil {
			log.Fatal("(B) User is null")
		}
		userId = user.Id

		// CASE 2: Record D.N.E.
	} else {
		var email string

		// CASE 2A: Email specified
		if om.Email.Valid {
			email = om.Email.String

			// CASE 2B: Email is not specified
		} else {
			customerIdStr := strconv.FormatUint(om.Id, 10)
			email = "customer+" + customerIdStr + "@workery.ca"
		}

		user, err := ur.GetByEmail(ctx, email)
		if err != nil {
			log.Panic("(C)", err)
		}
		if user == nil {
			um := &models.User{
				Uuid:        uuid.NewString(),
				FirstName:   om.GivenName.String,
				LastName:    om.LastName.String,
				Name:        name,
				LexicalName: lexicalName,
				Email:       email,
				// JoinedTime:     om.DateJoined,
				State:    state,
				Timezone: "America/Toronto",
				// CreatedTime:       om.DateJoined,
				// ModifiedTime:      om.LastModified,
				Salt:              "",
				WasEmailActivated: false,
				PrAccessCode:      "",
				PrExpiryTime:      time.Now(),
				TenantId:          tid,
				RoleId:            5, // Customer
			}
			err = ur.InsertOrUpdateByEmail(ctx, um)
			if err != nil {
				log.Panic("(D)", err)
			}
			user, err = ur.GetByEmail(ctx, email)
			if err != nil {
				log.Panic("(E)", err)
			}
		}
		userId = user.Id
	}

	//
	// Get `createdById` and `createdByName` values.
	//

	var createdById null.Int
	var createdByName null.String
	if om.CreatedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(om.CreatedById.ValueOrZero()))

		if err != nil {
			log.Panic("ur.GetIdByOldId", err)
		}
		user, err := ur.GetById(ctx, userId)
		if err != nil {
			log.Panic("ur.GetById", err)
		}

		if user != nil {
			createdById = null.IntFrom(int64(userId))
			createdByName = null.StringFrom(user.Name)
		} else {
			log.Println("WARNING: D.N.E.")
		}

		// // For debugging purposes only.
		// log.Println("createdById:", createdById)
		// log.Println("createdByName:", createdByName)
	}

	//
	// Get `lastModifiedById` and `lastModifiedByName` values.
	//

	var lastModifiedById null.Int
	var lastModifiedByName null.String
	if om.LastModifiedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(om.LastModifiedById.ValueOrZero()))
		if err != nil {
			log.Panic("ur.GetIdByOldId", err)
		}
		user, err := ur.GetById(ctx, userId)
		if err != nil {
			log.Panic("ur.GetById", err)
		}

		if user != nil {
			lastModifiedById = null.IntFrom(int64(userId))
			lastModifiedByName = null.StringFrom(user.Name)
		} else {
			log.Println("WARNING: D.N.E.")
		}

		// // For debugging purposes only.
		// log.Println("lastModifiedById:", lastModifiedById)
		// log.Println("lastModifiedByName:", lastModifiedByName)
	}

	//
	// Compile the `full address` and `address url`.
	//

	address := ""
	if om.StreetAddress != "" && om.StreetAddress != "-" {
		address += om.StreetAddress
	}
	if om.StreetAddressExtra.IsZero() != false && om.StreetAddressExtra.ValueOrZero() != "" {
		address += om.StreetAddressExtra.ValueOrZero()
	}
	if om.StreetAddress != "" && om.StreetAddress != "-" {
		address += ", "
	}
	address += om.AddressLocality
	address += ", " + om.AddressRegion
	address += ", " + om.AddressCountry
	fullAddressWithoutPostalCode := address
	fullAddressWithPostalCode := "-"
	fullAddressUrl := ""
	if om.PostalCode.String != "" {
		fullAddressWithPostalCode = address + ", " + om.PostalCode.String
		fullAddressUrl = "https://www.google.com/maps/place/" + fullAddressWithPostalCode
	} else {
		fullAddressUrl = "https://www.google.com/maps/place/" + fullAddressWithoutPostalCode
	}

	//
	// Compile the `how hear` text.
	//

	howHearId := uint64(om.HowHearId.Int64)
	howHearText := ""
	howHear, err := r.GetById(ctx, howHearId)
	if err != nil {
		log.Fatal(err)
		return
	}
	if howHearId == 1 {
		if om.HowHearOther == "" {
			howHearText = "-"
		} else {
			howHearText = om.HowHearOther
		}
	} else {
		howHearText = howHear.Text
	}

	//
	// Insert our `Customer` data.
	//

	m := &models.Customer{
		OldId:                        om.Id,
		Uuid:                         uuid.NewString(),
		TenantId:                     tid,
		UserId:                       userId,
		TypeOf:                       om.TypeOf,
		IndexedText:                  om.IndexedText.String,
		IsOkToEmail:                  om.IsOkToEmail,
		IsOkToText:                   om.IsOkToText,
		IsBusiness:                   om.IsBusiness,
		IsSenior:                     om.IsSenior,
		IsSupport:                    om.IsSupport,
		JobInfoRead:                  om.JobInfoRead.String,
		HowHearId:                    uint64(om.HowHearId.Int64),
		HowHearOld:                   om.HowHearOld,
		HowHearOther:                 om.HowHearOther,
		HowHearText:                  howHearText,
		State:                        state,
		DeactivationReason:           om.DeactivationReason,
		DeactivationReasonOther:      om.DeactivationReasonOther,
		CreatedTime:                  om.Created,
		CreatedById:                  createdById,
		CreatedByName:                createdByName,
		CreatedFromIP:                om.CreatedFrom.String,
		LastModifiedTime:             om.LastModified,
		LastModifiedById:             lastModifiedById,
		LastModifiedByName:           lastModifiedByName,
		LastModifiedFromIP:           om.LastModifiedFrom.String,
		OrganizationName:             om.OrganizationName.String,
		AddressCountry:               om.AddressCountry,
		AddressRegion:                om.AddressRegion,
		AddressLocality:              om.AddressLocality,
		PostOfficeBoxNumber:          om.PostOfficeBoxNumber.String,
		PostalCode:                   om.PostalCode.String,
		StreetAddress:                om.StreetAddress,
		StreetAddressExtra:           om.StreetAddressExtra.String,
		FullAddressWithoutPostalCode: fullAddressWithoutPostalCode,
		FullAddressWithPostalCode:    fullAddressWithPostalCode,
		FullAddressUrl:               fullAddressUrl,
		GivenName:                    om.GivenName.String,
		MiddleName:                   om.MiddleName.String,
		LastName:                     om.LastName.String,
		Name:                         name,
		LexicalName:                  lexicalName,
		Birthdate:                    om.Birthdate,
		JoinDate:                     om.JoinDate,
		Nationality:                  om.Nationality.String,
		Gender:                       om.Gender.String,
		TaxId:                        om.TaxId.String,
		Elevation:                    om.Elevation.Float64,
		Latitude:                     om.Latitude.Float64,
		Longitude:                    om.Longitude.Float64,
		AreaServed:                   om.AreaServed.String,
		AvailableLanguage:            om.AvailableLanguage.String,
		ContactType:                  om.ContactType.String,
		Email:                        om.Email.String,
		FaxNumber:                    om.FaxNumber.String,
		Telephone:                    om.Telephone.String,
		TelephoneTypeOf:              om.TelephoneTypeOf,
		TelephoneExtension:           om.TelephoneExtension.String,
		OtherTelephone:               om.OtherTelephone.String,
		OtherTelephoneExtension:      om.OtherTelephoneExtension.String,
		OtherTelephoneTypeOf:         om.OtherTelephoneTypeOf,
	}

	// fmt.Println(m) // For debugging purposes only.

	err = omr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", om.Id)
}
