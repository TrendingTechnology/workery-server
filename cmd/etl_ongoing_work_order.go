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
	owoETLSchemaName string
	owoETLTenantId   int
)

func init() {
	owoETLCmd.Flags().StringVarP(&owoETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	owoETLCmd.MarkFlagRequired("schema_name")
	owoETLCmd.Flags().IntVarP(&owoETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	owoETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(owoETLCmd)
}

var owoETLCmd = &cobra.Command{
	Use:   "etl_ongoing_work_order",
	Short: "Import the associate vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportOngoingWorkOrder()
	},
}

func doRunImportOngoingWorkOrder() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, owoETLSchemaName)
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	asr := repositories.NewOngoingWorkOrderRepo(db)
	ar := repositories.NewAssociateRepo(db)
	cr := repositories.NewCustomerRepo(db)

	runOngoingWorkOrderETL(ctx, uint64(owoETLTenantId), asr, ar, cr, oldDb)
}

func runOngoingWorkOrderETL(
	ctx context.Context,
	tenantId uint64,
	asr *repositories.OngoingWorkOrderRepo,
	ar *repositories.AssociateRepo,
	cr *repositories.CustomerRepo,
	oldDb *sql.DB,
) {
	ass, err := ListAllOngoingWorkOrders(oldDb)
	if err != nil {
		log.Fatal("ListAllOngoingWorkOrders", err)
	}
	for _, oss := range ass {
		insertOngoingWorkOrderETL(ctx, tenantId, asr, ar, cr, oss)
	}
}

type OldOngoingWorkOrder struct {
	Id               uint64      `json:"id"`
	State            string      `json:"state"`
	AssociateId      null.Int    `json:"associate_id"`
	CustomerId       uint64      `json:"customer_id"`
	CreatedAt        time.Time   `json:"created_at"`
	CreatedById      null.Int    `json:"created_by_id"`
	CreatedFrom      null.String `json:"created_from"`
	LastModifiedAt   time.Time   `json:"last_modified_at"`
	LastModifiedById null.Int    `json:"last_modified_by_id"`
	LastModifiedFrom null.String `json:"last_modified_from"`
}

func ListAllOngoingWorkOrders(db *sql.DB) ([]*OldOngoingWorkOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, state, associate_id, customer_id, created_at, created_by_id, created_from, last_modified_at, last_modified_by_id, last_modified_from
	FROM
        workery_ongoing_work_orders
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldOngoingWorkOrder
	defer rows.Close()
	for rows.Next() {
		m := new(OldOngoingWorkOrder)
		err = rows.Scan(
			&m.Id,
			&m.State,
			&m.AssociateId,
			&m.CustomerId,
			&m.CreatedAt,
			&m.CreatedById,
			&m.CreatedFrom,
			&m.LastModifiedAt,
			&m.LastModifiedById,
			&m.LastModifiedFrom,
		)
		if err != nil {
			log.Fatal("ListAllOngoingWorkOrders | rows.Scan", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("ListAllOngoingWorkOrders | rows.Err", err)
	}
	return arr, err
}

func insertOngoingWorkOrderETL(
	ctx context.Context,
	tid uint64,
	asr *repositories.OngoingWorkOrderRepo,
	ar *repositories.AssociateRepo,
	cr *repositories.CustomerRepo,
	oss *OldOngoingWorkOrder,
) {
	var associateId null.Int
	if oss.AssociateId.Valid {
		associateIdInt64 := oss.AssociateId.ValueOrZero()
		associateIdUint64, err := ar.GetIdByOldId(ctx, tid, uint64(associateIdInt64))
		if err != nil {
			log.Panic("ar.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		associateId = null.NewInt(int64(associateIdUint64), associateIdUint64 != 0)
	}

	customerId, err := cr.GetIdByOldId(ctx, tid, oss.CustomerId)

	var state int8 = 1 // Running
	if oss.State == "terminated" {
		state = 2
	}

	m := &models.OngoingWorkOrder{
		OldId:              oss.Id,
		TenantId:           tid,
		Uuid:               uuid.NewString(),
		CustomerId:         customerId,
		AssociateId:        associateId,
		State:              state,
		CreatedTime:        oss.CreatedAt,
		CreatedById:        oss.CreatedById,
		CreatedFromIP:      oss.CreatedFrom,
		LastModifiedTime:   oss.LastModifiedAt,
		LastModifiedById:   oss.LastModifiedById,
		LastModifiedFromIP: oss.LastModifiedFrom,
	}
	err = asr.Insert(ctx, m)
	if err != nil {
		log.Print("associateId", associateId)
		log.Print("customerId", customerId)
		log.Panic("asr.Insert | err", err, "\n\n", m, oss)
	}
	fmt.Println("Imported ID#", oss.Id)
}
