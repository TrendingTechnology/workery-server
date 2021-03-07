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
	ssirTenantId int
	ssirFilePath string
)

func init() {
	ssirETLCmd.Flags().IntVarP(&ssirTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	ssirETLCmd.MarkFlagRequired("tenant_id")
	ssirETLCmd.Flags().StringVarP(&ssirFilePath, "filepath", "f", "", "Path to the workery insurance requirement csv file.")
	ssirETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(ssirETLCmd)
}

var ssirETLCmd = &cobra.Command{
	Use:   "etl_skill_set_insurance_requirement",
	Short: "Import the insurance_requirement data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportSkillSetInsuranceRequirement()
	},
}

func doRunImportSkillSetInsuranceRequirement() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewSkillSetInsuranceRequirementRepo(db)

	f, err := os.Open(ssirFilePath)
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

		saveSkillSetInsuranceRequirementRowInDb(r, line)
	}
}

func saveSkillSetInsuranceRequirementRowInDb(r *repositories.SkillSetInsuranceRequirementRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	idString := col[0]
	skillSetIdString := col[1]
	insuranceRequirementIdString := col[2]

	id, _ := strconv.ParseUint(idString, 10, 64)
	skillSetId, _ := strconv.ParseUint(skillSetIdString, 10, 64)
	insuranceRequirementId, _ := strconv.ParseUint(insuranceRequirementIdString, 10, 64)

	if id != 0 {
		m := &models.SkillSetInsuranceRequirement{
			OldId: id,
			TenantId: uint64(ssirTenantId),
			Uuid: uuid.NewString(),
			SkillSetId: skillSetId,
			InsuranceRequirementId: insuranceRequirementId,
		}
		ctx := context.Background()
		err := r.Insert(ctx, m)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Imported ID#", id)
	}
}
