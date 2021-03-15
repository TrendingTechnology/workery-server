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
	etlhhauiTenantId int
	etlhhauiFilePath string
)

func init() {
	howHearAboutUsItemETLCmd.Flags().IntVarP(&etlhhauiTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	howHearAboutUsItemETLCmd.MarkFlagRequired("tenant_id")
	howHearAboutUsItemETLCmd.Flags().StringVarP(&etlhhauiFilePath, "filepath", "f", "", "Path to the workery insurance requirement csv file.")
	howHearAboutUsItemETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(howHearAboutUsItemETLCmd)
}

var howHearAboutUsItemETLCmd = &cobra.Command{
	Use:   "etl_how_hear_about_us_item",
	Short: "Import the how_hear_about_us_item data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportHowHearAboutUsItem()
	},
}

func doRunImportHowHearAboutUsItem() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewHowHearAboutUsItemRepo(db)

	f, err := os.Open(etlhhauiFilePath)
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

		saveHowHearAboutUsItemRowInDb(r, line)
	}
}

func saveHowHearAboutUsItemRowInDb(r *repositories.HowHearAboutUsItemRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	idString := col[0]
	text := col[1]
	IsForAssociateStr := col[2]
	IsForCustomerStr := col[3]
	IsForStaffStr := col[4]
	IsForPartneStrr := col[5]
	// SortNumbeStrr := col[6]
	StateStr := col[7]

	var IsForAssociate bool = false
	if IsForAssociateStr == "t" {
		IsForAssociate = true
	}

	var IsForCustomer bool = false
	if IsForCustomerStr == "t" {
		IsForCustomer = true
	}

	var IsForStaff bool = false
	if IsForStaffStr == "t" {
		IsForStaff = true
	}

	var IsForPartner bool = false
	if IsForPartneStrr == "t" {
		IsForPartner = true
	}

	var state int8 = 1
	if StateStr == "f" {
		state = 1
	}

	id, _ := strconv.ParseUint(idString, 10, 64)
	if id != 0 {
		m := &models.HowHearAboutUsItem{
			OldId: id,
			TenantId: uint64(etlhhauiTenantId),
			Uuid: uuid.NewString(),
			Text: text,
			IsForAssociate: IsForAssociate,
			IsForCustomer: IsForCustomer,
			IsForStaff: IsForStaff,
			IsForPartner: IsForPartner,
			SortNumber: 1,
			State: state,
		}
		ctx := context.Background()
		err := r.Insert(ctx, m)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Imported ID#", id)
	}
}
