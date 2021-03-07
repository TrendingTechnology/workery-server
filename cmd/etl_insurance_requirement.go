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
	etlirTenantId int
	etlirFilePath string
)

func init() {
	insuranceRequirementETLCmd.Flags().IntVarP(&etlirTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	insuranceRequirementETLCmd.MarkFlagRequired("tenant_id")
	insuranceRequirementETLCmd.Flags().StringVarP(&etlirFilePath, "filepath", "f", "", "Path to the workery insurance requirement csv file.")
	insuranceRequirementETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(insuranceRequirementETLCmd)
}

var insuranceRequirementETLCmd = &cobra.Command{
	Use:   "etl_insurance_requirement",
	Short: "Import the insurance_requirement data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportInsuranceRequirement()
	},
}

func doRunImportInsuranceRequirement() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewInsuranceRequirementRepo(db)

	f, err := os.Open(etlirFilePath)
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

		saveInsuranceRequirementRowInDb(r, line)
	}
}

func saveInsuranceRequirementRowInDb(r *repositories.InsuranceRequirementRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	idString := col[0]
	text := col[1]
	description := col[2]
	stateString := col[3]

	var state int8
	if stateString == "f" {
		state = 1
	} else {
		state = 0
	}

	id, _ := strconv.ParseUint(idString, 10, 64)
	if id != 0 {
		m := &models.InsuranceRequirement{
			Id: id,
			TenantId: uint64(etlirTenantId),
			Uuid: uuid.NewString(),
			Text: text,
			Description: description,
			State: state,
		}
		ctx := context.Background()
		err := r.InsertOrUpdateByText(ctx, m)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Imported ID#", id)
	}
}
