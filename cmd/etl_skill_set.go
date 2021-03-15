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
	etlssTenantId int
	etlssFilePath string
)

func init() {
	skillSetETLCmd.Flags().IntVarP(&etlssTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	skillSetETLCmd.MarkFlagRequired("tenant_id")
	skillSetETLCmd.Flags().StringVarP(&etlssFilePath, "filepath", "f", "", "Path to the workery insurance requirement csv file.")
	skillSetETLCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(skillSetETLCmd)
}

var skillSetETLCmd = &cobra.Command{
	Use:   "etl_skill_set",
	Short: "Import the skill_set data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportSkillSet()
	},
}

func doRunImportSkillSet() {
	// Load up our database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
	    log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	r := repositories.NewSkillSetRepo(db)

	f, err := os.Open(etlssFilePath)
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

		saveSkillSetRowInDb(r, line)
	}
}

func saveSkillSetRowInDb(r *repositories.SkillSetRepo, col []string) {
	// For debugging purposes only.
	// log.Println(col)

	// Extract the row.
	idString := col[0]
	category := col[1]
	subCategory := col[2]
	description := col[3]
	stateString := col[4]

	var state int8 = 1
	if stateString == "t" {
		state = 0
	}

	id, _ := strconv.ParseUint(idString, 10, 64)
	if id != 0 {
		m := &models.SkillSet{
			OldId: id,
			TenantId: uint64(etlssTenantId),
			Uuid: uuid.NewString(),
			Category: category,
			SubCategory: subCategory,
			Description: description,
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
