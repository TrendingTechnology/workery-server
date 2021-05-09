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
	workOrderSkillSetSchemaName string
)

func init() {
	workOrderSkillSetCmd.Flags().StringVarP(&workOrderSkillSetSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	workOrderSkillSetCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(workOrderSkillSetCmd)
}

var workOrderSkillSetCmd = &cobra.Command{
	Use:   "etl_work_order_skill_set",
	Short: "Import the work order skill_sets from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportWorkOrderSkillSet()
	},
}

func doRunImportWorkOrderSkillSet() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, workOrderSkillSetSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	wossr := repositories.NewWorkOrderSkillSetRepo(db)
	ar := repositories.NewWorkOrderRepo(db)
	vtr := repositories.NewSkillSetRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, workOrderSkillSetSchemaName)
	if err != nil {
		log.Fatal(err)
	}
    if tenant != nil {
		runWorkOrderSkillSetETL(ctx, uint64(tenant.Id), wossr, ar, vtr, oldDb)
	}
}

func runWorkOrderSkillSetETL(
	ctx context.Context,
	tenantId uint64,
	wossr *repositories.WorkOrderSkillSetRepo,
	wor *repositories.WorkOrderRepo,
	vtr *repositories.SkillSetRepo,
	oldDb *sql.DB,
) {
	aats, err := ListAllWorkOrderSkillSets(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range aats {
		insertWorkOrderSkillSetETL(ctx, tenantId, wossr, wor, vtr, oss)
	}
}

type OldWorkOrderSkillSet struct {
	Id          uint64 `json:"id"`
	WorkOrderId uint64 `json:"workorder_id"`
	SkillSetId       uint64 `json:"skillset_id"`
}

func ListAllWorkOrderSkillSets(db *sql.DB) ([]*OldWorkOrderSkillSet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, workorder_id, skillset_id
	FROM
        workery_work_orders_skill_sets
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldWorkOrderSkillSet
	defer rows.Close()
	for rows.Next() {
		m := new(OldWorkOrderSkillSet)
		err = rows.Scan(
			&m.Id,
			&m.WorkOrderId,
			&m.SkillSetId,
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

func insertWorkOrderSkillSetETL(
	ctx context.Context,
	tid uint64,
	wossr *repositories.WorkOrderSkillSetRepo,
	wor *repositories.WorkOrderRepo,
	vtr *repositories.SkillSetRepo,
	oss *OldWorkOrderSkillSet,
) {
	//
	// OrderId
	//

	orderId, err := wor.GetIdByOldId(ctx, tid, oss.WorkOrderId)
	if err != nil {
		log.Panic("wor.GetIdByOldId | err", err)
	}

	//
	// SkillSetId
	//

	skillSetId, err := vtr.GetIdByOldId(ctx, tid, oss.SkillSetId)
	if err != nil {
		log.Panic("vtr.GetIdByOldId | err", err)
	}

	//
	// Insert into database.
	//

	m := &models.WorkOrderSkillSet{
		OldId:       oss.Id,
		TenantId:    tid,
		Uuid:        uuid.NewString(),
		OrderId: orderId,
		SkillSetId:       skillSetId,
	}
	err = wossr.Insert(ctx, m)
	if err != nil {
		log.Panic("wossr.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
