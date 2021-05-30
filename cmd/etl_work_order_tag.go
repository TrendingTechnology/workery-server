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
	workOrderTagSchemaName string
)

func init() {
	workOrderTagCmd.Flags().StringVarP(&workOrderTagSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	workOrderTagCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(workOrderTagCmd)
}

var workOrderTagCmd = &cobra.Command{
	Use:   "etl_work_order_tag",
	Short: "Import the work order tags from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportWorkOrderTag()
	},
}

func doRunImportWorkOrderTag() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, workOrderTagSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	wotp := repositories.NewWorkOrderTagRepo(db)
	ar := repositories.NewWorkOrderRepo(db)
	vtr := repositories.NewTagRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, workOrderTagSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant != nil {
		runWorkOrderTagETL(ctx, uint64(tenant.Id), wotp, ar, vtr, oldDb)
	}
}

func runWorkOrderTagETL(
	ctx context.Context,
	tenantId uint64,
	wotp *repositories.WorkOrderTagRepo,
	ar *repositories.WorkOrderRepo,
	vtr *repositories.TagRepo,
	oldDb *sql.DB,
) {
	aats, err := ListAllWorkOrderTags(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range aats {
		insertWorkOrderTagETL(ctx, tenantId, wotp, ar, vtr, oss)
	}
}

type OldWorkOrderTag struct {
	Id          uint64 `json:"id"`
	WorkOrderId uint64 `json:"workorder_id"`
	TagId       uint64 `json:"tag_id"`
}

func ListAllWorkOrderTags(db *sql.DB) ([]*OldWorkOrderTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, workorder_id, tag_id
	FROM
        workery_work_orders_tags
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldWorkOrderTag
	defer rows.Close()
	for rows.Next() {
		m := new(OldWorkOrderTag)
		err = rows.Scan(
			&m.Id,
			&m.WorkOrderId,
			&m.TagId,
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

func insertWorkOrderTagETL(
	ctx context.Context,
	tid uint64,
	wotp *repositories.WorkOrderTagRepo,
	ar *repositories.WorkOrderRepo,
	vtr *repositories.TagRepo,
	oss *OldWorkOrderTag,
) {
	//
	// OrderId
	//

	orderId, err := ar.GetIdByOldId(ctx, tid, oss.WorkOrderId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	//
	// TagId
	//

	tagId, err := vtr.GetIdByOldId(ctx, tid, oss.TagId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	//
	// Insert into database.
	//

	m := &models.WorkOrderTag{
		OldId:    oss.Id,
		TenantId: tid,
		Uuid:     uuid.NewString(),
		OrderId:  orderId,
		TagId:    tagId,
	}
	err = wotp.Insert(ctx, m)
	if err != nil {
		log.Panic("wotp.Insert | err", err)
	}
	fmt.Println("Imported ID#", oss.Id)
}
