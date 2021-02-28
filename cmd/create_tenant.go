package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	repo "github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	ctName string
	ctState int
	ctTimezone string
	ctAddressCountry string
	ctAddressRegion string
	ctAddressLocality string
	ctPostOfficeBoxNumber string
	ctPostalCode string
	ctStreetAddress string
	ctStreetAddressExtra string
	// elevation,
	// latitude,
	// longitude,
	ctAreaServed string
	ctAvailableLanguage string
	// contact_type,
	ctEmail string
	// fax_number,
	ctTelephone string
	// telephone_type_of,
	// telephone_extension,
	// other_telephone,
	// other_telephone_extension,
	// other_telephone_type_of
)

func init() {
	createTenantCmd.Flags().StringVarP(&ctName, "name", "a", "", "Name of the tenant")
	createTenantCmd.MarkFlagRequired("name")
	createTenantCmd.Flags().IntVarP(&ctState, "state", "b", 0, "State of the tenant")
	createTenantCmd.MarkFlagRequired("state")
	createTenantCmd.Flags().StringVarP(&ctName, "timezone", "c", "", "Timezone of the tenant")
	createTenantCmd.MarkFlagRequired("timezone")
	createTenantCmd.Flags().StringVarP(&ctAddressCountry, "address_country", "d", "", "")
	createTenantCmd.MarkFlagRequired("address_country")
	createTenantCmd.Flags().StringVarP(&ctAddressRegion, "address_region", "e", "", "")
	createTenantCmd.MarkFlagRequired("address_region")
	createTenantCmd.Flags().StringVarP(&ctAddressLocality, "address_locality", "f", "", "")
	createTenantCmd.MarkFlagRequired("address_locality")
	createTenantCmd.Flags().StringVarP(&ctPostOfficeBoxNumber, "post_office_box_number", "g", "", "")
	createTenantCmd.MarkFlagRequired("post_office_box_number")
	createTenantCmd.Flags().StringVarP(&ctPostalCode, "postal_code", "i", "", "")
	createTenantCmd.MarkFlagRequired("postal_code")
	createTenantCmd.Flags().StringVarP(&ctStreetAddress, "street_address", "j", "", "")
	createTenantCmd.MarkFlagRequired("street_address")
	createTenantCmd.Flags().StringVarP(&ctStreetAddressExtra, "street_address_extra", "k", "", "")
	createTenantCmd.MarkFlagRequired("street_address_extra")
	// elevation,
	// latitude,
	// longitude,
	createTenantCmd.Flags().StringVarP(&ctStreetAddressExtra, "area_served", "l", "", "")
	createTenantCmd.MarkFlagRequired("area_served")
	createTenantCmd.Flags().StringVarP(&ctAvailableLanguage, "available_language", "m", "", "")
	createTenantCmd.MarkFlagRequired("available_language")
	// contact_type,
	createTenantCmd.Flags().StringVarP(&ctEmail, "email", "n", "", "")
	createTenantCmd.MarkFlagRequired("email")
	// fax_number,
	createTenantCmd.Flags().StringVarP(&ctTelephone, "telephone", "o", "", "")
	createTenantCmd.MarkFlagRequired("telephone")
	// telephone_type_of,
	// telephone_extension,
	// other_telephone,
	// other_telephone_extension,
	// other_telephone_type_of

	rootCmd.AddCommand(createTenantCmd)
}

var createTenantCmd = &cobra.Command{
	Use:   "create_tenant",
	Short: "Create a tenant",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runCreateTenant()
	},
}

func runCreateTenant() {
	ctx := context.Background()

	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repo.NewTenantRepo(db)

	// Check to see the tenant account already exists.
	tenantFound, _ := r.CheckIfExistsByName(ctx, ctName)
	if tenantFound {
		log.Fatal("Email already exists.")
	}

	m := &models.Tenant{
		Uuid: uuid.NewString(),
		Name: ctName,
		State: int8(ctState),
		Timezone: ctTimezone,
		CreatedTime: time.Now(),
		ModifiedTime: time.Now(),
		AddressCountry: ctAddressCountry,
		AddressRegion: ctAddressRegion,
		AddressLocality: ctAddressLocality,
		PostOfficeBoxNumber: ctPostOfficeBoxNumber,
		PostalCode: ctPostalCode,
		StreetAddress: ctStreetAddress,
		StreetAddressExtra: ctStreetAddressExtra,
		// elevation,
		// latitude,
		// longitude,
		// area_served,
		AvailableLanguage: ctAvailableLanguage,
		// contact_type,
		Email: ctEmail,
		// fax_number,
		Telephone: ctTelephone,
		// telephone_type_of,
		// telephone_extension,
		// other_telephone,
		// other_telephone_extension,
		// other_telephone_type_of
	}

	err = r.InsertOrUpdate(ctx, m)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("Tenant created with UUID:", m.Uuid)
}
