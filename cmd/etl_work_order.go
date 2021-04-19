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
	woETLSchemaName string
	woETLTenantId   int
)

func init() {
	woETLCmd.Flags().StringVarP(&woETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	woETLCmd.MarkFlagRequired("schema_name")
	woETLCmd.Flags().IntVarP(&woETLTenantId, "tenant_id", "t", 0, "Tenant Id that this data belongs to")
	woETLCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(woETLCmd)
}

var woETLCmd = &cobra.Command{
	Use:   "etl_work_order",
	Short: "Import the associate vehicle types from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportWorkOrder()
	},
}

func doRunImportWorkOrder() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, woETLSchemaName)
	if err != nil {
		log.Fatal("utils.ConnectDB", err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	asr := repositories.NewWorkOrderRepo(db)
	ar := repositories.NewAssociateRepo(db)
	cr := repositories.NewCustomerRepo(db)

	runWorkOrderETL(ctx, uint64(woETLTenantId), asr, ar, cr, oldDb)
}

func runWorkOrderETL(
	ctx context.Context,
	tenantId uint64,
	asr *repositories.WorkOrderRepo,
	ar *repositories.AssociateRepo,
	cr *repositories.CustomerRepo,
	oldDb *sql.DB,
) {
	ass, err := ListAllWorkOrders(oldDb)
	if err != nil {
		log.Fatal("ListAllWorkOrders", err)
	}
	for _, oss := range ass {
		insertWorkOrderETL(ctx, tenantId, asr, ar, cr, oss)
	}
}

type OldWorkOrder struct {
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

func ListAllWorkOrders(db *sql.DB) ([]*OldWorkOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, state, associate_id, customer_id, created_at, created_by_id, created_from, last_modified_at, last_modified_by_id, last_modified_from
	FROM
        workery_work_orders
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldWorkOrder
	defer rows.Close()
	for rows.Next() {
		m := new(OldWorkOrder)
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
			log.Fatal("ListAllWorkOrders | rows.Scan", err)
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("ListAllWorkOrders | rows.Err", err)
	}
	return arr, err
}

func insertWorkOrderETL(
	ctx context.Context,
	tid uint64,
	asr *repositories.WorkOrderRepo,
	ar *repositories.AssociateRepo,
	cr *repositories.CustomerRepo,
	oss *OldWorkOrder,
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

	m := &models.WorkOrder{
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
