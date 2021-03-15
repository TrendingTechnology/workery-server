package cmd

import (
	"context"
	"fmt"
	"os"
	"log"
	"database/sql"
	"time"

	"github.com/spf13/cobra"
	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	insuranceRequirementETLSchemaName string
)

func init() {
	insuranceRequirementETLCmd.Flags().StringVarP(&insuranceRequirementETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	insuranceRequirementETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(insuranceRequirementETLCmd)
}

var insuranceRequirementETLCmd = &cobra.Command{
	Use:   "etl_insurance_requirement",
	Short: "Import the insurance requirement data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportInsuranceRequirement()
	},
}

func doRunImportInsuranceRequirement() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, insuranceRequirementETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

    // Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewInsuranceRequirementRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, insuranceRequirementETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}

	runInsuranceRequirementETL(ctx, tenant.Id, irr, oldDb)
}

type OldUInsuranceRequirement struct {
	Id                      uint64 `json:"id"`
	Text                    string `json:"text"`
	Description             string `json:"description"`
	IsArchived              bool   `json:"is_archived"`
}

func ListAllInsuranceRequirements(db *sql.DB) ([]*OldUInsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, text, description, is_archived
	FROM
	    workery_insurance_requirements
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUInsuranceRequirement
	defer rows.Close()
	for rows.Next() {
		m := new(OldUInsuranceRequirement)
		err = rows.Scan(
			&m.Id,
			&m.Text,
			&m.Description,
			&m.IsArchived,
		)
		if err != nil {
			panic(err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func runInsuranceRequirementETL(ctx context.Context, tenantId uint64, irr *repositories.InsuranceRequirementRepo, oldDb *sql.DB) {
	insuranceRequirements, err := ListAllInsuranceRequirements(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range insuranceRequirements {
		insertInsuranceRequirementETL(ctx, tenantId, irr, oir)
	}
}

func insertInsuranceRequirementETL(ctx context.Context, tid uint64, irr *repositories.InsuranceRequirementRepo, oir *OldUInsuranceRequirement) {
	var state int8 = 1
	if oir.IsArchived == true {
		state = 0
	}

	m := &models.InsuranceRequirement{
		OldId: oir.Id,
		TenantId: tid,
		Uuid: uuid.NewString(),
		Text: oir.Text,
		Description: oir.Description,
		State: state,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oir.Id)
}
