package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	ssirETLSchemaName string
	ssirETLTenantId   int
)

func init() {
	ssirETLCmd.Flags().StringVarP(&ssirETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	ssirETLCmd.MarkFlagRequired("schema_name")
	ssirETLCmd.Flags().IntVarP(&ssirETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	ssirETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(ssirETLCmd)
}

var ssirETLCmd = &cobra.Command{
	Use:   "etl_skill_set_insurance_requirement",
	Short: "Import the insurance requirement and skill set data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportSkillSetInsuranceRequirement()
	},
}

func doRunImportSkillSetInsuranceRequirement() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, ssirETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	ssir := repositories.NewSkillSetInsuranceRequirementRepo(db)

	runSkillSetInsuranceRequirementETL(ctx, uint64(ssirETLTenantId), ssir, oldDb)
}

func runSkillSetInsuranceRequirementETL(ctx context.Context, tenantId uint64, ssir *repositories.SkillSetInsuranceRequirementRepo, oldDb *sql.DB) {
	ssirs, err := ListAllSkillSetInsuranceRequirements(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range ssirs {
		insertSkillSetInsuranceRequirementETL(ctx, tenantId, ssir, oss)
	}
}

type OldSkillSetInsuranceRequirement struct {
	Id                     uint64 `json:"id"`
	SkillSetId             uint64 `json:"skill_set_id"`
	InsuranceRequirementId uint64 `json:"insurance_requirement_id"`
}

func ListAllSkillSetInsuranceRequirements(db *sql.DB) ([]*OldSkillSetInsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, skillset_id, insurancerequirement_id
	FROM
        workery_skill_sets_insurance_requirements
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldSkillSetInsuranceRequirement
	defer rows.Close()
	for rows.Next() {
		m := new(OldSkillSetInsuranceRequirement)
		err = rows.Scan(
			&m.Id,
			&m.SkillSetId,
			&m.InsuranceRequirementId,
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

func insertSkillSetInsuranceRequirementETL(ctx context.Context, tid uint64, ssir *repositories.SkillSetInsuranceRequirementRepo, oss *OldSkillSetInsuranceRequirement) {
	m := &models.SkillSetInsuranceRequirement{
		OldId:                  oss.Id,
		TenantId:               tid,
		Uuid:                   uuid.NewString(),
		SkillSetId:             oss.SkillSetId,
		InsuranceRequirementId: oss.InsuranceRequirementId,
	}
	err := ssir.Insert(ctx, m)
	if err != nil {
		log.Panic("ssir.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
