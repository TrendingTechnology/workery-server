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
	airETLSchemaName string
	airETLTenantId   int
)

func init() {
	airETLCmd.Flags().StringVarP(&airETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	airETLCmd.MarkFlagRequired("schema_name")
	airETLCmd.Flags().IntVarP(&airETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	airETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(airETLCmd)
}

var airETLCmd = &cobra.Command{
	Use:   "etl_associate_insurance_requirement",
	Short: "Import the associate vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportAssociateInsuranceRequirement()
	},
}

func doRunImportAssociateInsuranceRequirement() {
	// Load up our NEW database.
	db, err := utils.ConnectDB(databaseHost, databasePort, databaseUser, databasePassword, databaseName, "public")
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer db.Close()

	// Load up our OLD database.
	oldDBHost := os.Getenv("WORKERY_OLD_DB_HOST")
	oldDBPort := os.Getenv("WORKERY_OLD_DB_PORT")
	oldDBUser := os.Getenv("WORKERY_OLD_DB_USER")
	oldDBPassword := os.Getenv("WORKERY_OLD_DB_PASSWORD")
	oldDBName := os.Getenv("WORKERY_OLD_DB_NAME")
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, airETLSchemaName)
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	asr := repositories.NewAssociateInsuranceRequirementRepo(db)
	ar := repositories.NewAssociateRepo(db)
	vtr := repositories.NewInsuranceRequirementRepo(db)

	runAssociateInsuranceRequirementETL(ctx, uint64(airETLTenantId), asr, ar, vtr, oldDb)
}

func runAssociateInsuranceRequirementETL(
	ctx context.Context,
	tenantId uint64,
	asr *repositories.AssociateInsuranceRequirementRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.InsuranceRequirementRepo,
	oldDb *sql.DB,
) {
	ass, err := ListAllAssociateInsuranceRequirements(oldDb)
	if err != nil {
		log.Fatal("ListAllAssociateInsuranceRequirements", err)
	}
	for _, oss := range ass {
		insertAssociateInsuranceRequirementETL(ctx, tenantId, asr, ar, vtr, oss)
	}
}

type OldAssociateInsuranceRequirement struct {
	Id                     uint64 `json:"id"`
	AssociateId            uint64 `json:"associate_id"`
	InsuranceRequirementId              uint64 `json:"insurancerequirement_id"`
}

func ListAllAssociateInsuranceRequirements(db *sql.DB) ([]*OldAssociateInsuranceRequirement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, associate_id, insurancerequirement_id
	FROM
        workery_associates_insurance_requirements
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldAssociateInsuranceRequirement
	defer rows.Close()
	for rows.Next() {
		m := new(OldAssociateInsuranceRequirement)
		err = rows.Scan(
			&m.Id,
			&m.AssociateId,
			&m.InsuranceRequirementId,
		)
		if err != nil {
			log.Fatal("rows.Scan", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("rows.Err", err)
	}
	return arr, err
}

func insertAssociateInsuranceRequirementETL(
	ctx context.Context,
	tid uint64,
	asr *repositories.AssociateInsuranceRequirementRepo,
	ar *repositories.AssociateRepo,
	vtr *repositories.InsuranceRequirementRepo,
	oss *OldAssociateInsuranceRequirement,
) {
	associateId, err := ar.GetIdByOldId(ctx, tid, oss.AssociateId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	commentId, err := vtr.GetIdByOldId(ctx, tid, oss.InsuranceRequirementId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	if associateId != 0 && commentId != 0 {
		m := &models.AssociateInsuranceRequirement{
			OldId:                  oss.Id,
			TenantId:               tid,
			Uuid:                   uuid.NewString(),
			AssociateId:            associateId,
			InsuranceRequirementId:              commentId,
		}
		err = asr.Insert(ctx, m)
		if err != nil {
			log.Panic("asr.Insert | err", err, "\n\n", m, oss)
		}
		fmt.Println("Imported ID#", oss.Id)
	} else {
		fmt.Println("-------------------\nSkipped ID#", oss.Id, "\n-------------------\nassociateId #", associateId, "\ncommentId #", commentId, "\n\noss", oss, "\n\n")
	}
}
