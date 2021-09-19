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
	associateETLSchemaName string
)

func init() {
	associateETLCmd.Flags().StringVarP(&associateETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	associateETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(associateETLCmd)
}

var associateETLCmd = &cobra.Command{
	Use:   "etl_associate",
	Short: "Import the associate data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportAssociate()
	},
}

func doRunImportAssociate() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, associateETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tenantRepo := repositories.NewTenantRepo(db)
	userRepo := repositories.NewUserRepo(db)
	associateRepo := repositories.NewAssociateRepo(db)
	serviceFeeRepo := repositories.NewWorkOrderServiceFeeRepo(db)
	r := repositories.NewHowHearAboutUsItemRepo(db)

	// Lookup the tenant.
	tenant, err := tenantRepo.GetBySchemaName(ctx, associateETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runAssociateETL(ctx, tenant.Id, userRepo, associateRepo, serviceFeeRepo, r, oldDb)
}

type OldUAssociate struct {
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
	CreatedFrom              null.String `json:"created_from"`
	CreatedFromIsPublic      bool        `json:"created_from_is_public"`
	LastModifiedFrom         null.String `json:"last_modified_from"`
	LastModifiedFromIsPublic bool        `json:"last_modified_from_is_public"`
	CreatedById              null.Int    `json:"created_by_id"`
	LastModifiedById         null.Int    `json:"last_modified_by_id"`
	OwnerId                  null.Int    `json:"owner_id"`
	HowHearOther             string      `json:"how_hear_other"`
	IsArchived               bool        `json:"is_archived"`
	HowHearId                null.Int    `json:"how_hear_id"`
	HowHearOld               int8        `json:"how_hear_old"`
	OrganizationName         null.String `json:"organization_name"`
	OrganizationTypeOf       int8        `json:"organization_type_of"`
	AvatarImageId            null.Int    `json:"avatar_image_id"`
	ServiceFeeId             null.Int    `json:"service_fee_id"`
}

func ListAllAssociates(db *sql.DB) ([]*OldUAssociate, error) {
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
		is_ok_to_text,created_from, created_from_is_public, last_modified_from,
		last_modified_from_is_public, is_archived, created_by_id, last_modified_by_id,
		owner_id, how_hear_other,
		how_hear_id, how_hear_old, organization_name,
		organization_type_of, avatar_image_id, service_fee_id
	FROM
	    workery_associates
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUAssociate
	defer rows.Close()
	for rows.Next() {
		m := new(OldUAssociate)
		err = rows.Scan(
			&m.Id, &m.Created, &m.LastModified, &m.AlternateName, &m.Description, &m.Name, &m.Url,
			&m.AreaServed, &m.AvailableLanguage, &m.ContactType, &m.Email, &m.FaxNumber,
			&m.ProductSupported, &m.Telephone, &m.TelephoneTypeOf, &m.TelephoneExtension,
			&m.OtherTelephone, &m.OtherTelephoneExtension, &m.OtherTelephoneTypeOf,
			&m.AddressCountry, &m.AddressLocality, &m.AddressRegion, &m.PostOfficeBoxNumber,
			&m.PostalCode, &m.StreetAddress, &m.StreetAddressExtra, &m.Elevation,
			&m.Latitude, &m.Longitude, &m.GivenName, &m.MiddleName, &m.LastName,
			&m.Birthdate, &m.JoinDate, &m.Nationality, &m.Gender, &m.TaxId, &m.Id, &m.IndexedText, &m.TypeOf,
			&m.IsOkToEmail, &m.IsOkToText,
			&m.CreatedFrom, &m.CreatedFromIsPublic, &m.LastModifiedFrom,
			&m.LastModifiedFromIsPublic, &m.IsArchived, &m.CreatedById, &m.LastModifiedById,
			&m.OwnerId, &m.HowHearOther,
			&m.HowHearId, &m.HowHearOld, &m.OrganizationName,
			&m.OrganizationTypeOf, &m.AvatarImageId, &m.ServiceFeeId,
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

func runAssociateETL(
	ctx context.Context,
	tenantId uint64,
	userRepo *repositories.UserRepo,
	associateRepo *repositories.AssociateRepo,
	serviceFeeRepo *repositories.WorkOrderServiceFeeRepo,
	r *repositories.HowHearAboutUsItemRepo,
	oldDb *sql.DB,
) {
	associates, err := ListAllAssociates(oldDb)
	if err != nil {
		log.Fatal("runAssociateETL | ListAllAssociates | err:", err)
	}
	for _, oldAssociate := range associates {
		var serviceFeeId uint64

		oldServiceFeeId, err := oldAssociate.ServiceFeeId.Value()

		if err != nil {
			log.Panic("runAssociateETL | oldAssociate.ServiceFeeId | err", err)
		}

		// DEVELOPERS NOTE:
		// THIS IS TECH DEBT! BUST SINCE WE IMPORT THE DATA, THIS ETL WILL
		// BE DELETED SO WE ACCEPT THIS TECH DEBT.
		// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
		if oldServiceFeeId == nil {
			log.Println("runAssociateETL | oldServiceFeeId == nil | oldAssociate.ServiceFeeId", oldAssociate.Id, oldAssociate.ServiceFeeId)
			if tenantId == 2 {
				serviceFeeId = 3
			} else if tenantId == 4 {
				serviceFeeId = 11
			} else {

			}
		} else {
			serviceFeeId, err = serviceFeeRepo.GetIdByOldId(ctx, tenantId, serviceFeeId)
			if err != nil {
				log.Panic("runAssociateETL | serviceFeeRepo.GetIdByOldId | err", err)
			}
		}
		if serviceFeeId == 0 {
			if tenantId == 2 {
				serviceFeeId = 3
			} else if tenantId == 4 {
				serviceFeeId = 11
			} else {
				log.Println("runAssociateETL | serviceFeeId", serviceFeeId)
			}
		}
		// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

		insertAssociateETL(ctx, tenantId, serviceFeeId, userRepo, associateRepo, r, oldAssociate)
	}
}

func insertAssociateETL(
	ctx context.Context,
	tenantId uint64,
	serviceFeeId uint64,
	userRepo *repositories.UserRepo,
	associateRepo *repositories.AssociateRepo,
	r *repositories.HowHearAboutUsItemRepo,
	oldAssociate *OldUAssociate,
) {
	var state int8 = 1
	if oldAssociate.IsArchived == true {
		state = 0
	}

	// Variable used to keep the ID of the user record in our database.
	userId := uint64(oldAssociate.OwnerId.Int64)

	// Generate our full name / lexical full name.
	var name string
	var lexicalName string
	if oldAssociate.MiddleName.Valid {
		name = oldAssociate.GivenName.String + " " + oldAssociate.MiddleName.String + " " + oldAssociate.LastName.String
		lexicalName = oldAssociate.LastName.String + ", " + oldAssociate.MiddleName.String + ", " + oldAssociate.GivenName.String
	} else {
		name = oldAssociate.GivenName.String + " " + oldAssociate.LastName.String
		lexicalName = oldAssociate.LastName.String + ", " + oldAssociate.GivenName.String
	}

	lexicalName = strings.Replace(lexicalName, ", ,", ",", 0)
	lexicalName = strings.Replace(lexicalName, "  ", " ", 0)
	lexicalName = strings.Replace(lexicalName, ", , ", ", ", 0)
	name = strings.Replace(name, "  ", " ", 0)

	// CASE 1: User record exists in our database.
	if oldAssociate.OwnerId.Valid {
		// log.Println(oldAssociate.Id)
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
		if oldAssociate.Email.Valid {
			email = oldAssociate.Email.String

			// CASE 2B: Email is not specified
		} else {
			associateIdStr := strconv.FormatUint(oldAssociate.Id, 10)
			email = "associate+" + associateIdStr + "@workery.ca"
		}

		user, err := userRepo.GetByEmail(ctx, email)
		if err != nil {
			log.Panic("(C)", err)
		}
		if user == nil {
			um := &models.User{
				Uuid:        uuid.NewString(),
				FirstName:   oldAssociate.GivenName.String,
				LastName:    oldAssociate.LastName.String,
				Name:        name,
				LexicalName: lexicalName,
				Email:       email,
				// JoinedTime:        oldAssociate.DateJoined,
				State:    state,
				Timezone: "America/Toronto",
				// CreatedTime:       oldAssociate.DateJoined,
				// ModifiedTime:      oldAssociate.LastModified,
				Salt:              "",
				WasEmailActivated: false,
				PrAccessCode:      "",
				PrExpiryTime:      time.Now(),
				TenantId:          tenantId,
				RoleId:            5, // Associate
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

	var createdById uint64
	if oldAssociate.CreatedById.Valid {
		createdById = uint64(oldAssociate.CreatedById.Int64)
		user, err := userRepo.GetByOldId(ctx, createdById)
		if err != nil {
			log.Fatal("(F)", err)
		}
		if user == nil {
			log.Fatal("(G) User is null")
		}
	} else {
		createdById = userId
	}

	var lastModifiedById uint64
	if oldAssociate.LastModifiedById.Valid {
		lastModifiedById = uint64(oldAssociate.LastModifiedById.Int64)
		user, err := userRepo.GetByOldId(ctx, lastModifiedById)
		if err != nil {
			log.Fatal("(F)", err)
		}
		if user == nil {
			log.Fatal("(G) User is null")
		}
	} else {
		lastModifiedById = userId
	}

	//
	// Compile the `full address` and `address url`.
	//

	address := ""
	if oldAssociate.StreetAddress != "" && oldAssociate.StreetAddress != "-" {
		address += oldAssociate.StreetAddress
	}
	if oldAssociate.StreetAddressExtra.IsZero() != false && oldAssociate.StreetAddressExtra.ValueOrZero() != "" {
		address += oldAssociate.StreetAddressExtra.ValueOrZero()
	}
	if oldAssociate.StreetAddress != "" && oldAssociate.StreetAddress != "-" {
		address += ", "
	}
	address += oldAssociate.AddressLocality
	address += ", " + oldAssociate.AddressRegion
	address += ", " + oldAssociate.AddressCountry
	fullAddressWithoutPostalCode := address
	fullAddressWithPostalCode := "-"
	fullAddressUrl := ""
	if oldAssociate.PostalCode.String != "" {
		fullAddressWithPostalCode = address + ", " + oldAssociate.PostalCode.String
		fullAddressUrl = "https://www.google.com/maps/place/" + fullAddressWithPostalCode
	} else {
		fullAddressUrl = "https://www.google.com/maps/place/" + fullAddressWithoutPostalCode
	}

	//
	// Compile the `how hear` text.
	//

	howHearId := uint64(oldAssociate.HowHearId.Int64)
	howHearText := ""
	howHear, err := r.GetById(ctx, howHearId)
	if err != nil {
		log.Fatal(err)
		return
	}
	if howHearId == 1 {
		if oldAssociate.HowHearOther == "" {
			howHearText = "-"
		} else {
			howHearText = oldAssociate.HowHearOther
		}
	} else {
		howHearText = howHear.Text
	}

	//
	// Insert our `Customer` data.
	//

	associate := &models.Associate{
		OldId:                        oldAssociate.Id,
		ServiceFeeId:                 serviceFeeId,
		Uuid:                         uuid.NewString(),
		TenantId:                     tenantId,
		UserId:                       userId,
		TypeOf:                       oldAssociate.TypeOf,
		IndexedText:                  oldAssociate.IndexedText.String,
		IsOkToEmail:                  oldAssociate.IsOkToEmail,
		IsOkToText:                   oldAssociate.IsOkToText,
		HowHearId:                    uint64(oldAssociate.HowHearId.Int64),
		HowHearOld:                   oldAssociate.HowHearOld,
		HowHearOther:                 oldAssociate.HowHearOther,
		HowHearText:                  howHearText,
		State:                        state,
		CreatedTime:                  oldAssociate.Created,
		CreatedById:                  createdById,
		CreatedFromIP:                oldAssociate.CreatedFrom.String,
		LastModifiedTime:             oldAssociate.LastModified,
		LastModifiedById:             lastModifiedById,
		LastModifiedFromIP:           oldAssociate.LastModifiedFrom.String,
		OrganizationName:             oldAssociate.OrganizationName.String,
		AddressCountry:               oldAssociate.AddressCountry,
		AddressRegion:                oldAssociate.AddressRegion,
		AddressLocality:              oldAssociate.AddressLocality,
		PostOfficeBoxNumber:          oldAssociate.PostOfficeBoxNumber.String,
		PostalCode:                   oldAssociate.PostalCode.String,
		StreetAddress:                oldAssociate.StreetAddress,
		StreetAddressExtra:           oldAssociate.StreetAddressExtra.String,
		FullAddressWithoutPostalCode: fullAddressWithoutPostalCode,
		FullAddressWithPostalCode:    fullAddressWithPostalCode,
		FullAddressUrl:               fullAddressUrl,
		GivenName:                    oldAssociate.GivenName.String,
		MiddleName:                   oldAssociate.MiddleName.String,
		Name:                         name,
		LexicalName:                  lexicalName,
		LastName:                     oldAssociate.LastName.String,
		Birthdate:                    oldAssociate.Birthdate,
		JoinDate:                     oldAssociate.JoinDate,
		Nationality:                  oldAssociate.Nationality.String,
		Gender:                       oldAssociate.Gender.String,
		TaxId:                        oldAssociate.TaxId.String,
		Elevation:                    oldAssociate.Elevation.Float64,
		Latitude:                     oldAssociate.Latitude.Float64,
		Longitude:                    oldAssociate.Longitude.Float64,
		AreaServed:                   oldAssociate.AreaServed.String,
		AvailableLanguage:            oldAssociate.AvailableLanguage.String,
		ContactType:                  oldAssociate.ContactType.String,
		Email:                        oldAssociate.Email.String,
		FaxNumber:                    oldAssociate.FaxNumber.String,
		Telephone:                    oldAssociate.Telephone.String,
		TelephoneTypeOf:              oldAssociate.TelephoneTypeOf,
		TelephoneExtension:           oldAssociate.TelephoneExtension.String,
		OtherTelephone:               oldAssociate.OtherTelephone.String,
		OtherTelephoneExtension:      oldAssociate.OtherTelephoneExtension.String,
		OtherTelephoneTypeOf:         oldAssociate.OtherTelephoneTypeOf,
	}

	err = associateRepo.Insert(ctx, associate)
	if err != nil {
		log.Panic("omr.Insert:", err)
	}
	fmt.Println("Imported ID#", oldAssociate.Id)
}
