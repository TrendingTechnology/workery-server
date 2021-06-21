package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

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

	// Lookup the tenant.
	tenant, err := tenantRepo.GetBySchemaName(ctx, partnerETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runPartnerETL(ctx, tenant.Id, userRepo, partnerRepo, serviceFeeRepo, oldDb)
}

type OldUPartner struct {
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
	// TypeOf                   int8            `json:"type_of"`
	IsOkToEmail              bool            `json:"is_ok_to_email"`
	IsOkToText               bool            `json:"is_ok_to_text"`
	CreatedFrom              sql.NullString  `json:"created_from"`
	CreatedFromIsPublic      bool            `json:"created_from_is_public"`
	LastModifiedFrom         sql.NullString  `json:"last_modified_from"`
	LastModifiedFromIsPublic bool            `json:"last_modified_from_is_public"`
	CreatedById              sql.NullInt64   `json:"created_by_id"`
	LastModifiedById         sql.NullInt64   `json:"last_modified_by_id"`
	OwnerId                  sql.NullInt64   `json:"owner_id"`
	IsArchived               bool            `json:"is_archived"`
	HowHearId                sql.NullInt64   `json:"how_hear_id"`
	HowHearOld               int8            `json:"how_hear_old"`
	OrganizationName         sql.NullString  `json:"organization_name"`
	OrganizationTypeOf       int8            `json:"organization_type_of"`
	AvatarImageId            sql.NullInt64   `json:"avatar_image_id"`
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
	oldDb *sql.DB,
) {
	partners, err := ListAllPartners(oldDb)
	if err != nil {
		log.Fatal("runPartnerETL | ListAllPartners | err:", err)
	}
	for _, oldPartner := range partners {
		insertPartnerETL(ctx, tenantId, userRepo, partnerRepo, oldPartner)
	}
}

func insertPartnerETL(
	ctx context.Context,
	tenantId uint64,
	userRepo *repositories.UserRepo,
	partnerRepo *repositories.PartnerRepo,
	oldPartner *OldUPartner,
) {
	var state int8 = 1
	if oldPartner.IsArchived == true {
		state = 0
	}

	// Variable used to keep the ID of the user record in our database.
	userId := uint64(oldPartner.OwnerId.Int64)

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
				Uuid:      uuid.NewString(),
				FirstName: oldPartner.GivenName.String,
				LastName:  oldPartner.LastName.String,
				Email:     email,
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
				Role:              5, // Partner
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
	if oldPartner.CreatedById.Valid {
		createdById = uint64(oldPartner.CreatedById.Int64)
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
	if oldPartner.LastModifiedById.Valid {
		lastModifiedById = uint64(oldPartner.LastModifiedById.Int64)
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

	partner := &models.Partner{
		OldId:                   oldPartner.Id,
		Uuid:                    uuid.NewString(),
		TenantId:                tenantId,
		UserId:                  userId,
		// TypeOf:                  oldPartner.TypeOf,
		IndexedText:             oldPartner.IndexedText.String,
		IsOkToEmail:             oldPartner.IsOkToEmail,
		IsOkToText:              oldPartner.IsOkToText,
		HowHearId:               uint64(oldPartner.HowHearId.Int64),
		HowHearOld:              oldPartner.HowHearOld,
		// HowHearOther:            oldPartner.HowHearOther,
		State:                   state,
		CreatedTime:             oldPartner.Created,
		CreatedById:             createdById,
		CreatedFromIP:           oldPartner.CreatedFrom.String,
		LastModifiedTime:        oldPartner.LastModified,
		LastModifiedById:        lastModifiedById,
		LastModifiedFromIP:      oldPartner.LastModifiedFrom.String,
		OrganizationName:        oldPartner.OrganizationName.String,
		AddressCountry:          oldPartner.AddressCountry,
		AddressRegion:           oldPartner.AddressRegion,
		AddressLocality:         oldPartner.AddressLocality,
		PostOfficeBoxNumber:     oldPartner.PostOfficeBoxNumber.String,
		PostalCode:              oldPartner.PostalCode.String,
		StreetAddress:           oldPartner.StreetAddress,
		StreetAddressExtra:      oldPartner.StreetAddressExtra.String,
		GivenName:               oldPartner.GivenName.String,
		MiddleName:              oldPartner.MiddleName.String,
		LastName:                oldPartner.LastName.String,
		Birthdate:               oldPartner.Birthdate,
		JoinDate:                oldPartner.JoinDate,
		Nationality:             oldPartner.Nationality.String,
		Gender:                  oldPartner.Gender.String,
		TaxId:                   oldPartner.TaxId.String,
		Elevation:               oldPartner.Elevation.Float64,
		Latitude:                oldPartner.Latitude.Float64,
		Longitude:               oldPartner.Longitude.Float64,
		AreaServed:              oldPartner.AreaServed.String,
		AvailableLanguage:       oldPartner.AvailableLanguage.String,
		ContactType:             oldPartner.ContactType.String,
		Email:                   oldPartner.Email.String,
		FaxNumber:               oldPartner.FaxNumber.String,
		Telephone:               oldPartner.Telephone.String,
		TelephoneTypeOf:         oldPartner.TelephoneTypeOf,
		TelephoneExtension:      oldPartner.TelephoneExtension.String,
		OtherTelephone:          oldPartner.OtherTelephone.String,
		OtherTelephoneExtension: oldPartner.OtherTelephoneExtension.String,
		OtherTelephoneTypeOf:    oldPartner.OtherTelephoneTypeOf,
	}

	err := partnerRepo.Insert(ctx, partner)
	if err != nil {
		log.Panic("omr.Insert:", err)
	}
	fmt.Println("Imported ID#", oldPartner.Id)
}
