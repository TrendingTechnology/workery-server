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
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	staffETLSchemaName string
)

func init() {
	staffETLCmd.Flags().StringVarP(&staffETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	staffETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(staffETLCmd)
}

var staffETLCmd = &cobra.Command{
	Use:   "etl_staff",
	Short: "Import the staff data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportStaff()
	},
}

func doRunImportStaff() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, staffETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	ur := repositories.NewUserRepo(db)
	sr := repositories.NewStaffRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, staffETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runStaffETL(ctx, tenant.Id, ur, sr, oldDb)
}

type OldUStaff struct {
	Created                 time.Time       `json:"created"`
	LastModified            time.Time       `json:"last_modified"`
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
	AddressLocality         string          `json:"address_locality"`
	AddressRegion           string          `json:"address_region"`
	PostOfficeBoxNumber     sql.NullString  `json:"post_office_box_number"`
	PostalCode              sql.NullString  `json:"postal_code"`
	StreetAddress           string          `json:"street_address"`
	StreetAddressExtra      sql.NullString  `json:"street_address_extra"`
	Elevation               sql.NullFloat64 `json:"elevation"`
	Latitude                sql.NullFloat64 `json:"latitude"`
	Longitude               sql.NullFloat64 `json:"longitude"`
	GivenName               sql.NullString  `json:"given_name"`
	MiddleName              sql.NullString  `json:"middle_name"`
	LastName                sql.NullString  `json:"last_name"`
	Birthdate               sql.NullTime    `json:"birthdate"`
	JoinDate                sql.NullTime    `json:"join_date"`
	Nationality             sql.NullString  `json:"nationality"`
	Gender                  sql.NullString  `json:"gender"`
	TaxId                   sql.NullString  `json:"tax_id"`
	Id                      uint64          `json:"id"`
	IndexedText             sql.NullString  `json:"indexed_text"`
	CreatedFrom              sql.NullString `json:"created_from"`
	CreatedFromIsPublic      bool           `json:"created_from_is_public"`
	LastModifiedFrom         sql.NullString `json:"last_modified_from"`
	LastModifiedFromIsPublic bool           `json:"last_modified_from_is_public"`
	IsArchived               bool           `json:"is_archived"`
	CreatedById              sql.NullInt64  `json:"created_by_id"`
	LastModifiedById         sql.NullInt64  `json:"last_modified_by_id"`
	OwnerId      sql.NullInt64  `json:"owner_id"`
	HowHearOther sql.NullString `json:"how_hear_other"`
	HowHearId sql.NullInt64 `json:"how_hear_id"`
	AvatarImageId                        sql.NullInt64  `json:"avatar_image_id"`
	PersonalEmail                        sql.NullString `json:"personal_email"`
	EmergencyContactAlternativeTelephone sql.NullString `json:"emergency_contact_alternative_telephone"`
	EmergencyContactName                 sql.NullString `json:"emergency_contact_name"`
	EmergencyContactRelationship         sql.NullString `json:"emergency_contact_relationship"`
	EmergencyContactTelephone            sql.NullString `json:"emergency_contact_telephone"`
	PoliceCheck                          sql.NullTime   `json:"police_check"`
	// TypeOf                   int8            `json:"type_of"`
	// IsOkToEmail              bool            `json:"is_ok_to_email"`
	// IsOkToText               bool            `json:"is_ok_to_text"`
	// IsBusiness               bool            `json:"is_business"`
	// IsSenior                 bool            `json:"is_senior"`
	// IsSupport                bool            `json:"is_support"`
	// JobInfoRead              sql.NullString  `json:"job_info_read"`
	// OrganizationId           sql.NullInt64   `json:"organization_id"`
	// IsBlacklisted            bool            `json:"is_blacklisted"`
	// DeactivationReason       int8            `json:"deactivation_reason"`
	// DeactivationReasonOther  string          `json:"deactivation_reason_other"`
	// State                    string          `json:"state"`
	// HowHearOld               int8            `json:"how_hear_old"`
	// OrganizationName         sql.NullString  `json:"organization_name"`
	// OrganizationTypeOf       int8            `json:"organization_type_of"`
}

func ListAllStaffs(db *sql.DB) ([]*OldUStaff, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, created, last_modified, available_language, contact_type, email, fax_number,
		telephone, telephone_type_of, telephone_extension,
		other_telephone, other_telephone_extension, other_telephone_type_of,
		address_country, address_locality, address_region, post_office_box_number,
		postal_code, street_address, street_address_extra, elevation, latitude,
		longitude, given_name, middle_name, last_name, birthdate, join_date,
		nationality, gender, tax_id, id, indexed_text,
		created_from, created_from_is_public, last_modified_from,
		last_modified_from_is_public, is_archived, created_by_id, last_modified_by_id,
		owner_id, how_hear_other, how_hear_id, avatar_image_id, personal_email,
		emergency_contact_alternative_telephone, emergency_contact_name,
		emergency_contact_relationship, emergency_contact_telephone, police_check
	FROM
	    workery_staff
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUStaff
	defer rows.Close()
	for rows.Next() {
		m := new(OldUStaff)
		err = rows.Scan(
			&m.Id, &m.Created, &m.LastModified, &m.AvailableLanguage, &m.ContactType, &m.Email, &m.FaxNumber,
			&m.Telephone, &m.TelephoneTypeOf, &m.TelephoneExtension,
			&m.OtherTelephone, &m.OtherTelephoneExtension, &m.OtherTelephoneTypeOf,
			&m.AddressCountry, &m.AddressLocality, &m.AddressRegion, &m.PostOfficeBoxNumber,
			&m.PostalCode, &m.StreetAddress, &m.StreetAddressExtra, &m.Elevation,
			&m.Latitude, &m.Longitude, &m.GivenName, &m.MiddleName, &m.LastName,
			&m.Birthdate, &m.JoinDate, &m.Nationality, &m.Gender, &m.TaxId, &m.Id, &m.IndexedText,
			&m.CreatedFrom, &m.CreatedFromIsPublic, &m.LastModifiedFrom,
			&m.LastModifiedFromIsPublic, &m.IsArchived, &m.CreatedById, &m.LastModifiedById,
			&m.OwnerId, &m.HowHearOther, &m.HowHearId, &m.AvatarImageId, &m.PersonalEmail,
			&m.EmergencyContactAlternativeTelephone, &m.EmergencyContactName,
			&m.EmergencyContactRelationship, &m.EmergencyContactTelephone, &m.PoliceCheck,
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

func runStaffETL(ctx context.Context, tenantId uint64, ur *repositories.UserRepo, sr *repositories.StaffRepo, oldDb *sql.DB) {
	staffs, err := ListAllStaffs(oldDb)
	if err != nil {
		log.Fatal("(0000)", err)
	}
	for _, om := range staffs {
		insertStaffETL(ctx, tenantId, ur, sr, om)
	}
}

func insertStaffETL(ctx context.Context, tid uint64, ur *repositories.UserRepo, sr *repositories.StaffRepo, om *OldUStaff) {
	var state int8 = 1
	if om.IsArchived == true {
		state = 0
	}

	// Variable used to keep the ID of the user record in our database.
	userId := uint64(om.OwnerId.Int64)

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
			staffIdStr := strconv.FormatUint(om.Id, 10)
			email = "staff+" + staffIdStr + "@workery.ca"
		}

		user, err := ur.GetByEmail(ctx, email)
		if err != nil {
			log.Panic("(C)", err)
		}
		if user == nil {
			um := &models.User{
				Uuid:      uuid.NewString(),
				FirstName: om.GivenName.String,
				LastName:  om.LastName.String,
				Email:     email,
				// JoinedTime:        om.DateJoined,
				State:    state,
				Timezone: "America/Toronto",
				// CreatedTime:       om.DateJoined,
				// ModifiedTime:      om.LastModified,
				Salt:              "",
				WasEmailActivated: false,
				PrAccessCode:      "",
				PrExpiryTime:      time.Now(),
				TenantId:          tid,
				Role:              5, // Staff
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

	var createdById uint64
	if om.CreatedById.Valid {
		createdById = uint64(om.CreatedById.Int64)
		user, err := ur.GetByOldId(ctx, createdById)
		if err != nil {
			log.Fatal("(F)", err)
		}
		if user == nil {
			log.Fatal("(G) User is null")
		}
	} else {
		createdById = userId
	}
	etlCreatedById := null.NewInt(int64(createdById), createdById != 0)

	var lastModifiedById uint64
	if om.LastModifiedById.Valid {
		lastModifiedById = uint64(om.LastModifiedById.Int64)
		user, err := ur.GetByOldId(ctx, lastModifiedById)
		if err != nil {
			log.Fatal("(F)", err)
		}
		if user == nil {
			log.Fatal("(G) User is null")
		}
	} else {
		lastModifiedById = userId
	}
	etlLastModifiedById := null.NewInt(int64(lastModifiedById), lastModifiedById != 0)

	m := &models.Staff{
		OldId:                                om.Id,
		Uuid:                                 uuid.NewString(),
		TenantId:                             tid,
		UserId:                               userId,
		CreatedTime:                          om.Created,
		LastModifiedTime:                     om.LastModified,
		AvailableLanguage:                    om.AvailableLanguage,
		ContactType:                          om.ContactType,
		Email:                                om.Email,
		FaxNumber:                            om.FaxNumber,
		Telephone:                            om.Telephone,
		TelephoneTypeOf:                      om.TelephoneTypeOf,
		TelephoneExtension:                   om.TelephoneExtension,
		OtherTelephone:                       om.OtherTelephone,
		OtherTelephoneExtension:              om.OtherTelephoneExtension,
		OtherTelephoneTypeOf:                 om.OtherTelephoneTypeOf,
		AddressCountry:                       om.AddressCountry,
		AddressLocality:                      om.AddressLocality,
		AddressRegion:                        om.AddressRegion,
		PostOfficeBoxNumber:                  om.PostOfficeBoxNumber,
		PostalCode:                           om.PostalCode,
		StreetAddress:                        om.StreetAddress,
		StreetAddressExtra:                   om.StreetAddressExtra,
		Elevation:                            om.Elevation,
		Latitude:                             om.Latitude,
		Longitude:                            om.Longitude,
		GivenName:                            om.GivenName,
		MiddleName:                           om.MiddleName,
		LastName:                             om.LastName,
		Birthdate:                            om.Birthdate,
		JoinDate:                             om.JoinDate,
		Nationality:                          om.Nationality,
		Gender:                               om.Gender,
		TaxId:                                om.TaxId,
		IndexedText:                          om.IndexedText,
		CreatedFromIP:                        om.CreatedFrom,
		LastModifiedFromIP:                   om.LastModifiedFrom,
		CreatedById:                          etlCreatedById,
		LastModifiedById:                     etlLastModifiedById,
		HowHearOther:                         om.HowHearOther,
		HowHearId:                            om.HowHearId,
		AvatarImageId:                        om.AvatarImageId,
		PersonalEmail:                        om.PersonalEmail,
		EmergencyContactAlternativeTelephone: om.EmergencyContactAlternativeTelephone,
		EmergencyContactName:                 om.EmergencyContactName,
		EmergencyContactRelationship:         om.EmergencyContactRelationship,
		EmergencyContactTelephone:            om.EmergencyContactTelephone,
		PoliceCheck:                          om.PoliceCheck,
		State:                                state,
	}

	// fmt.Println(m) // For debugging purposes only.

	err := sr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", om.Id)
}
