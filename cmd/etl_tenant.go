package cmd

import (
	"bufio"
	"context"
	"fmt"
	"encoding/csv"
	"os"
	"log"
	"io"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	tenantFilePath string
)

func init() {
	tenantETLCmd.Flags().StringVarP(&tenantFilePath, "filepath", "f", "", "Path to the workery tenant csv file.")
	tenantETLCmd.MarkFlagRequired("filepath")
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
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewTenantRepo(db)

	f, err := os.Open(tenantFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// defer the closing of our `f` so that we can parse it later on
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))

	for {
		// Read line by line until no more lines left.
		line, error := reader.Read()
		if error == io.EOF {
			break
        } else if error != nil {
			log.Fatal(error)
		}

		saveTenantRowInDb(r, line)
	}
}

func saveTenantRowInDb(r *repositories.TenantRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	idString := col[0]
	createdTimeString := col[2]
	modifiedTimeString := col[3]
	alternate_name := col[4]
	description := col[5]
	name := col[6]
	url := col[7]
	timezone := "America/Toronto"
	// email := col[11]
	addressCountry := col[20]
	addressRegion := col[22]
	addressLocality := col[21]
	postalCode := col[24]
	streetAdress := col[25]
	streetAdressExtra := col[25]

	ct, _ := utils.ConvertPGAdminTimeStringToTime(createdTimeString)
	mt, _ := utils.ConvertPGAdminTimeStringToTime(modifiedTimeString)

	id, _ := strconv.ParseUint(idString, 10, 64)
	if id != 0 {
		m := &models.Tenant{
			Id: id,
			Uuid: uuid.NewString(),
			AlternateName: alternate_name,
			Description: description,
			Name: name,
			Url: url,
			State: 1,
			Timezone: timezone,
			CreatedTime: ct,
			ModifiedTime: mt,
			AddressCountry: addressCountry,
			AddressRegion: addressRegion,
			AddressLocality: addressLocality,
			PostalCode: postalCode,
			StreetAddress: streetAdress,
			StreetAddressExtra: streetAdressExtra,
		}
		ctx := context.Background()
		err := r.InsertOrUpdate(ctx, m)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Imported ID#", id)
	}
}
