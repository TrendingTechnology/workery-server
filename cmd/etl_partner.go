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
	partnerETLSchemaName string
)

func init() {
	partnerETLCmd.Flags().StringVarP(&partnerETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	partnerETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(partnerETLCmd)
}

var partnerETLCmd = &cobra.Command{
	Use:   "etl_partner",
	Short: "Import the partner data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportPartner()
	},
}

func doRunImportPartner() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, partnerETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tenantRepo := repositories.NewTenantRepo(db)
	userRepo := repositories.NewUserRepo(db)
	partnerRepo := repositories.NewPartnerRepo(db)
	serviceFeeRepo := repositories.NewWorkOrderServiceFeeRepo(db)
	r := repositories.NewHowHearAboutUsItemRepo(db)

	// Lookup the tenant.
	tenant, err := tenantRepo.GetBySchemaName(ctx, partnerETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runPartnerETL(ctx, tenant.Id, userRepo, partnerRepo, serviceFeeRepo, r, oldDb)
}

type OldUPartner struct {
	Created                 time.Time   `json:"created"`
	LastModified            time.Time   `json:"last_modified"`
	AlternateName           null.String `json:"alternate_name"`
	Description             null.String `json:"description"`
	Name                    null.String `json:"name"`
	Url                     null.String `json:"url"`
	AreaServed              null.String `json:"area_served"`
	AvailableLanguage       null.String `json:"available_language"`
	ContactType             null.String `json:"contact_type"`
	Email                   null.String `json:"email"`
	FaxNumber               null.String `json:"fax_number"`
	ProductSupported        null.String `json:"product_supported"`
	Telephone               null.String `json:"telephone"`
	TelephoneTypeOf         int8        `json:"telephone_type_of"`
	TelephoneExtension      null.String `json:"telephone_extension"`
	OtherTelephone          null.String `json:"other_telephone"`
	OtherTelephoneExtension null.String `json:"other_telephone_extension"`
	OtherTelephoneTypeOf    int8        `json:"other_telephone_type_of"`
	AddressCountry          string      `json:"address_country"`
	AddressLocality         string      `json:"address_locality"`
	AddressRegion           string      `json:"address_region"`
	PostOfficeBoxNumber     null.String `json:"post_office_box_number"`
	PostalCode              null.String `json:"postal_code"`
	StreetAddress           string      `json:"street_address"`
	StreetAddressExtra      null.String `json:"street_address_extra"`
	Elevation               null.Float  `json:"elevation"`
	Latitude                null.Float  `json:"latitude"`
	Longitude               null.Float  `json:"longitude"`
	GivenName               null.String `json:"given_name"`
	MiddleName              null.String `json:"middle_name"`
	LastName                null.String `json:"last_name"`
	Birthdate               null.Time   `json:"birthdate"`
	JoinDate                null.Time   `json:"join_date"`
	Nationality             null.String `json:"nationality"`
	Gender                  null.String `json:"gender"`
	TaxId                   null.String `json:"tax_id"`
	Id                      uint64      `json:"id"`
	IndexedText             null.String `json:"indexed_text"`
	// TypeOf                   int8            `json:"type_of"`
	IsOkToEmail              bool        `json:"is_ok_to_email"`
	IsOkToText               bool        `json:"is_ok_to_text"`
	CreatedFrom              null.String `json:"created_from"`
	CreatedFromIsPublic      bool        `json:"created_from_is_public"`
	LastModifiedFrom         null.String `json:"last_modified_from"`
	LastModifiedFromIsPublic bool        `json:"last_modified_from_is_public"`
	CreatedById              null.Int    `json:"created_by_id"`
	LastModifiedById         null.Int    `json:"last_modified_by_id"`
	OwnerId                  null.Int    `json:"owner_id"`
	IsArchived               bool        `json:"is_archived"`
	HowHearId                null.Int    `json:"how_hear_id"`
	HowHearOld               int8        `json:"how_hear_old"`
	OrganizationName         null.String `json:"organization_name"`
	OrganizationTypeOf       int8        `json:"organization_type_of"`
	AvatarImageId            null.Int    `json:"avatar_image_id"`
}

func ListAllPartners(db *sql.DB) ([]*OldUPartner, error) {
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
		nationality, gender, tax_id, id, indexed_text, is_ok_to_email,
		is_ok_to_text,created_from, created_from_is_public, last_modified_from,
		last_modified_from_is_public, is_archived, created_by_id, last_modified_by_id,
		owner_id,
		how_hear_id, how_hear_old, organization_name,
		organization_type_of, avatar_image_id
	FROM
	    workery_partners
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUPartner
	defer rows.Close()
	for rows.Next() {
		m := new(OldUPartner)
		err = rows.Scan(
			&m.Id, &m.Created, &m.LastModified, &m.AlternateName, &m.Description, &m.Name, &m.Url,
			&m.AreaServed, &m.AvailableLanguage, &m.ContactType, &m.Email, &m.FaxNumber,
			&m.ProductSupported, &m.Telephone, &m.TelephoneTypeOf, &m.TelephoneExtension,
			&m.OtherTelephone, &m.OtherTelephoneExtension, &m.OtherTelephoneTypeOf,
			&m.AddressCountry, &m.AddressLocality, &m.AddressRegion, &m.PostOfficeBoxNumber,
			&m.PostalCode, &m.StreetAddress, &m.StreetAddressExtra, &m.Elevation,
			&m.Latitude, &m.Longitude, &m.GivenName, &m.MiddleName, &m.LastName,
			&m.Birthdate, &m.JoinDate, &m.Nationality, &m.Gender, &m.TaxId, &m.Id, &m.IndexedText,
			&m.IsOkToEmail, &m.IsOkToText,
			&m.CreatedFrom, &m.CreatedFromIsPublic, &m.LastModifiedFrom,
			&m.LastModifiedFromIsPublic, &m.IsArchived, &m.CreatedById, &m.LastModifiedById,
			&m.OwnerId,
			&m.HowHearId, &m.HowHearOld, &m.OrganizationName,
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

func runPartnerETL(
	ctx context.Context,
	tenantId uint64,
	userRepo *repositories.UserRepo,
	partnerRepo *repositories.PartnerRepo,
	serviceFeeRepo *repositories.WorkOrderServiceFeeRepo,
	r *repositories.HowHearAboutUsItemRepo,
	oldDb *sql.DB,
) {
	partners, err := ListAllPartners(oldDb)
	if err != nil {
		log.Fatal("runPartnerETL | ListAllPartners | err:", err)
	}
	for _, oldPartner := range partners {
		insertPartnerETL(ctx, tenantId, userRepo, partnerRepo, r, oldPartner)
	}
}

func insertPartnerETL(
	ctx context.Context,
	tenantId uint64,
	userRepo *repositories.UserRepo,
	partnerRepo *repositories.PartnerRepo,
	r *repositories.HowHearAboutUsItemRepo,
	oldPartner *OldUPartner,
) {
	//
	// Set the `state`.
	//

	var state int8 = 1
	if oldPartner.IsArchived == true {
		state = 0
	}

	//
	// Get `UserId` value - Variable used to keep the ID of the user record in our database.
	//

	userId := uint64(oldPartner.OwnerId.Int64)

	//
	// Generate our full name / lexical full name.
	//

	var name string
	var lexicalName string
	if oldPartner.MiddleName.Valid {
		name = oldPartner.GivenName.String + " " + oldPartner.MiddleName.String + " " + oldPartner.LastName.String
		lexicalName = oldPartner.LastName.String + ", " + oldPartner.MiddleName.String + ", " + oldPartner.GivenName.String
	} else {
		name = oldPartner.GivenName.String + " " + oldPartner.LastName.String
		lexicalName = oldPartner.LastName.String + ", " + oldPartner.GivenName.String
	}
	lexicalName = strings.Replace(lexicalName, ", ,", ",", 0)
	lexicalName = strings.Replace(lexicalName, "  ", " ", 0)

	// CASE 1: User record exists in our database.
	if oldPartner.OwnerId.Valid {
		// log.Println(oldPartner.Id)
		user, err := userRepo.GetByOldId(ctx, userId)
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
		if oldPartner.Email.Valid {
			email = oldPartner.Email.String

			// CASE 2B: Email is not specified
		} else {
			partnerIdStr := strconv.FormatUint(oldPartner.Id, 10)
			email = "partner+" + partnerIdStr + "@workery.ca"
		}

		user, err := userRepo.GetByEmail(ctx, email)
		if err != nil {
			log.Panic("(C)", err)
		}
		if user == nil {
			um := &models.User{
				Uuid:        uuid.NewString(),
				FirstName:   oldPartner.GivenName.String,
				LastName:    oldPartner.LastName.String,
				Name:        name,
				LexicalName: lexicalName,
				Email:       email,
				// JoinedTime:        oldPartner.DateJoined,
				State:    state,
				Timezone: "America/Toronto",
				// CreatedTime:       oldPartner.DateJoined,
				// ModifiedTime:      oldPartner.LastModified,
				Salt:              "",
				WasEmailActivated: false,
				PrAccessCode:      "",
				PrExpiryTime:      time.Now(),
				TenantId:          tenantId,
				RoleId:            5, // Partner
			}
			err = userRepo.InsertOrUpdateByEmail(ctx, um)
			if err != nil {
				log.Panic("(D)", err)
			}
			user, err = userRepo.GetByEmail(ctx, email)
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
	if oldPartner.CreatedById.ValueOrZero() > 0 {
		userId, err := userRepo.GetIdByOldId(ctx, tenantId, uint64(oldPartner.CreatedById.ValueOrZero()))

		if err != nil {
			log.Panic("userRepo.GetIdByOldId", err)
		}
		user, err := userRepo.GetById(ctx, userId)
		if err != nil {
			log.Panic("userRepo.GetById", err)
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
	if oldPartner.LastModifiedById.ValueOrZero() > 0 {
		userId, err := userRepo.GetIdByOldId(ctx, tenantId, uint64(oldPartner.LastModifiedById.ValueOrZero()))
		if err != nil {
			log.Panic("userRepo.GetIdByOldId", err)
		}
		user, err := userRepo.GetById(ctx, userId)
		if err != nil {
			log.Panic("userRepo.GetById", err)
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
	if oldPartner.StreetAddress != "" && oldPartner.StreetAddress != "-" {
		address += oldPartner.StreetAddress
	}
	if oldPartner.StreetAddressExtra.IsZero() != false && oldPartner.StreetAddressExtra.ValueOrZero() != "" {
		address += oldPartner.StreetAddressExtra.ValueOrZero()
	}
	if oldPartner.StreetAddress != "" && oldPartner.StreetAddress != "-" {
		address += ", "
	}
	address += oldPartner.AddressLocality
	address += ", " + oldPartner.AddressRegion
	address += ", " + oldPartner.AddressCountry
	fullAddressWithoutPostalCode := address
	fullAddressWithPostalCode := "-"
	fullAddressUrl := ""
	if oldPartner.PostalCode.String != "" {
		fullAddressWithPostalCode = address + ", " + oldPartner.PostalCode.String
		fullAddressUrl = "https://www.google.com/maps/place/" + fullAddressWithPostalCode
	} else {
		fullAddressUrl = "https://www.google.com/maps/place/" + fullAddressWithoutPostalCode
	}

	//
	// Compile the `how hear` text.
	//

	howHearId := uint64(oldPartner.HowHearId.Int64)
	howHearText := ""
	howHear, err := r.GetById(ctx, howHearId)
	if err != nil {
		log.Fatal(err)
		return
	}
	if howHearId == 1 {
		howHearText = "-"
	} else {
		howHearText = howHear.Text
	}

	//
	// Insert our `Partner` data.
	//

	partner := &models.Partner{
		OldId:    oldPartner.Id,
		Uuid:     uuid.NewString(),
		TenantId: tenantId,
		UserId:   userId,
		// TypeOf:                  oldPartner.TypeOf,
		IndexedText: oldPartner.IndexedText.String,
		IsOkToEmail: oldPartner.IsOkToEmail,
		IsOkToText:  oldPartner.IsOkToText,
		HowHearId:   uint64(oldPartner.HowHearId.Int64),
		HowHearOld:  oldPartner.HowHearOld,
		HowHearText: howHearText,
		// HowHearOther:            oldPartner.HowHearOther,
		State:                        state,
		CreatedTime:                  oldPartner.Created,
		CreatedById:                  createdById,
		CreatedByName:                createdByName,
		CreatedFromIP:                oldPartner.CreatedFrom.String,
		LastModifiedTime:             oldPartner.LastModified,
		LastModifiedById:             lastModifiedById,
		LastModifiedByName:           lastModifiedByName,
		LastModifiedFromIP:           oldPartner.LastModifiedFrom.String,
		OrganizationName:             oldPartner.OrganizationName.String,
		AddressCountry:               oldPartner.AddressCountry,
		AddressRegion:                oldPartner.AddressRegion,
		AddressLocality:              oldPartner.AddressLocality,
		PostOfficeBoxNumber:          oldPartner.PostOfficeBoxNumber.String,
		PostalCode:                   oldPartner.PostalCode.String,
		StreetAddress:                oldPartner.StreetAddress,
		StreetAddressExtra:           oldPartner.StreetAddressExtra.String,
		FullAddressWithoutPostalCode: fullAddressWithoutPostalCode,
		FullAddressWithPostalCode:    fullAddressWithPostalCode,
		FullAddressUrl:               fullAddressUrl,
		GivenName:                    oldPartner.GivenName.String,
		MiddleName:                   oldPartner.MiddleName.String,
		LastName:                     oldPartner.LastName.String,
		Name:                         name,
		LexicalName:                  lexicalName,
		Birthdate:                    oldPartner.Birthdate,
		JoinDate:                     oldPartner.JoinDate,
		Nationality:                  oldPartner.Nationality.String,
		Gender:                       oldPartner.Gender.String,
		TaxId:                        oldPartner.TaxId.String,
		Elevation:                    oldPartner.Elevation.Float64,
		Latitude:                     oldPartner.Latitude.Float64,
		Longitude:                    oldPartner.Longitude.Float64,
		AreaServed:                   oldPartner.AreaServed.String,
		AvailableLanguage:            oldPartner.AvailableLanguage.String,
		ContactType:                  oldPartner.ContactType.String,
		Email:                        oldPartner.Email.String,
		FaxNumber:                    oldPartner.FaxNumber.String,
		Telephone:                    oldPartner.Telephone.String,
		TelephoneTypeOf:              oldPartner.TelephoneTypeOf,
		TelephoneExtension:           oldPartner.TelephoneExtension.String,
		OtherTelephone:               oldPartner.OtherTelephone.String,
		OtherTelephoneExtension:      oldPartner.OtherTelephoneExtension.String,
		OtherTelephoneTypeOf:         oldPartner.OtherTelephoneTypeOf,
	}

	err = partnerRepo.Insert(ctx, partner)
	if err != nil {
		log.Panic("omr.Insert:", err)
	}
	fmt.Println("Imported ID#", oldPartner.Id)
}
