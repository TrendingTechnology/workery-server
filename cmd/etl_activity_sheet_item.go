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
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/repositories"
	"github.com/over55/workery-server/internal/utils"
)

var (
	activitySheetItemETLSchemaName string
)

func init() {
	activitySheetItemETLCmd.Flags().StringVarP(&activitySheetItemETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	activitySheetItemETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(activitySheetItemETLCmd)
}

var activitySheetItemETLCmd = &cobra.Command{
	Use:   "etl_activity_sheet_item",
	Short: "Import the activity_sheet_items data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportActivitySheetItem()
	},
}

func doRunImportActivitySheetItem() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, activitySheetItemETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	ur := repositories.NewUserRepo(db)
	asir := repositories.NewActivitySheetItemRepo(db)
	owor := repositories.NewOngoingWorkOrderRepo(db)
	wor := repositories.NewWorkOrderRepo(db)
	ar := repositories.NewAssociateRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, activitySheetItemETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runActivitySheetItemETL(ctx, tenant.Id, ur, asir, owor, wor, ar, oldDb)
}

type OldUActivitySheetItem struct {
	Id          uint64 `json:"id"`
	Comment     string `json:"comment"`
    CreatedAt   time.Time   `json:"created_at"`
	CreatedFrom null.String `json:"created_from"`
    CreatedById null.Int `json:"created_by_id"`
	AssociateId  uint64 `json:"associate_id"`
	JobId        null.Int  `json:"job_id"`
	State        string `json:"state"`
	OngoingJobId null.Int `json:"ongoing_job_id"`
}

func ListAllActivitySheetItems(db *sql.DB) ([]*OldUActivitySheetItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, comment, created_at, created_from, created_by_id, associate_id, job_id, state, ongoing_job_id
	FROM
	    workery_activity_sheet_items
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUActivitySheetItem
	defer rows.Close()
	for rows.Next() {
		m := new(OldUActivitySheetItem)
		err = rows.Scan(
			&m.Id,
			&m.Comment,
			&m.CreatedAt,
			&m.CreatedFrom,
			&m.CreatedById,
			&m.AssociateId,
			&m.JobId,
			&m.State,
			&m.OngoingJobId,
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

func runActivitySheetItemETL(
	ctx context.Context,
	tenantId uint64,
	ur *repositories.UserRepo,
	asir *repositories.ActivitySheetItemRepo,
	owor *repositories.OngoingWorkOrderRepo,
	wor *repositories.WorkOrderRepo,
	ar *repositories.AssociateRepo,
	oldDb *sql.DB,
) {
	activitySheetItems, err := ListAllActivitySheetItems(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oldRecord := range activitySheetItems {
		insertActivitySheetItemETL(ctx, tenantId, ur, asir, owor, wor, ar, oldRecord)
	}
}

func insertActivitySheetItemETL(
	ctx context.Context,
	tid uint64,
	ur *repositories.UserRepo,
	asir *repositories.ActivitySheetItemRepo,
	owor *repositories.OngoingWorkOrderRepo,
	wor *repositories.WorkOrderRepo,
	ar *repositories.AssociateRepo,
	oldRecord *OldUActivitySheetItem,
) {
	//
	// State
	//

	var state int8 = 1
	if oldRecord.State == "pending" {
		state = 2
	} else if oldRecord.State == "accepted" {
		state = 3
	} else if oldRecord.State == "declined" {
		state = 1
	}

    //
	// CreatedById
	//

	var createdById null.Int
	if oldRecord.CreatedById.Valid {
		createdByIdInt64 := oldRecord.CreatedById.ValueOrZero()
		createdByIdUint64, err := ur.GetIdByOldId(ctx, tid, uint64(createdByIdInt64))
		if err != nil {
			log.Panic("ur.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		createdById = null.NewInt(int64(createdByIdUint64), createdByIdUint64 != 0)
	}

	//
	// AssociateId
	//

	associateId, err := ar.GetIdByOldId(ctx, tid, oldRecord.AssociateId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	//
	// WorkOrderId
	//

	var workOrderId null.Int
	if oldRecord.JobId.Valid {
		workOrderIdInt64 := oldRecord.JobId.ValueOrZero()
		workOrderIdUint64, err := wor.GetIdByOldId(ctx, tid, uint64(workOrderIdInt64))
		if err != nil {
			log.Panic("wor.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		workOrderId = null.NewInt(int64(workOrderIdUint64), workOrderIdUint64 != 0)
    }

	//
	// OngoingOrderId
	//

	var ongoingWorkOrderId null.Int
	if oldRecord.OngoingJobId.Valid {
		ongoingWorkOrderIdInt64 := oldRecord.OngoingJobId.ValueOrZero()
		ongoingWorkOrderIdUint64, err := owor.GetIdByOldId(ctx, tid, uint64(ongoingWorkOrderIdInt64))
		if err != nil {
			log.Panic("owor.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		ongoingWorkOrderId = null.NewInt(int64(ongoingWorkOrderIdUint64), ongoingWorkOrderIdUint64 != 0)
    }

	//
	// Insert into database
	//

	m := &models.ActivitySheetItem{
		OldId:         oldRecord.Id,
		TenantId:      tid,
		Uuid:          uuid.NewString(),
		Comment:       oldRecord.Comment,
		CreatedTime:   oldRecord.CreatedAt,
		CreatedFromIP: oldRecord.CreatedFrom,
		CreatedById:   createdById,
		AssociateId:   associateId,
		State:         state,
		OrderId:         workOrderId,
		OngoingOrderId:  ongoingWorkOrderId,
	}

	err = asir.Insert(ctx, m)
	if err != nil {
		log.Fatal("Aborted on ID#", oldRecord.Id, "via err:", err)
	} else {
		fmt.Println("Imported ID#", oldRecord.Id)
	}
}
