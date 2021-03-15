package cmd

import (
	"bufio"
	"context"
	"database/sql"
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
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

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

		saveSkillSetInsuranceRequirementRowInDb(db, line)
	}
}

func saveSkillSetInsuranceRequirementRowInDb(db *sql.DB, col []string) {
	ctx := context.Background()

	// Load up our repositories.
	ssirr := repositories.NewSkillSetInsuranceRequirementRepo(db)
	ssr := repositories.NewSkillSetRepo(db)
	irr := repositories.NewInsuranceRequirementRepo(db)

	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	tenantId := uint64(ssirTenantId)
	idString := col[0]
	skillSetIdString := col[1]
	insuranceRequirementIdString := col[2]

	id, _ := strconv.ParseUint(idString, 10, 64)
	skillSetId, _ := strconv.ParseUint(skillSetIdString, 10, 64)
	insuranceRequirementId, _ := strconv.ParseUint(insuranceRequirementIdString, 10, 64)

    // Lookup in the DB to return the NEW ID.
	skillSet, _ := ssr.GetByOld(ctx, tenantId, skillSetId)
	insuranceRequirement, _ := irr.GetByOld(ctx, tenantId, insuranceRequirementId)

	if id != 0 {
		m := &models.SkillSetInsuranceRequirement{
			OldId: id,
			TenantId: tenantId,
			Uuid: uuid.NewString(),
			SkillSetId: skillSet.Id,
			InsuranceRequirementId: insuranceRequirement.Id,
		}
		err := ssirr.Insert(ctx, m)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Imported ID#", id)
	}
}
