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
	workOrderServiceFeeETLSchemaName string
)

func init() {
	workOrderServiceFeeETLCmd.Flags().StringVarP(&workOrderServiceFeeETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	workOrderServiceFeeETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(workOrderServiceFeeETLCmd)
}

var workOrderServiceFeeETLCmd = &cobra.Command{
	Use:   "etl_work_order_service_fee",
	Short: "Import the workOrderServiceFee data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportWorkOrderServiceFee()
	},
}

func doRunImportWorkOrderServiceFee() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, workOrderServiceFeeETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

    // Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	irr := repositories.NewWorkOrderServiceFeeRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, workOrderServiceFeeETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runWorkOrderServiceFeeETL(ctx, tenant.Id, irr, oldDb)
}

type OldUWorkOrderServiceFee struct {
	Id                uint64    `json:"id"`
	Title              string    `json:"title"`
	Description       string    `json:"description"`
    Percentage        float64   `json:"percentage"`
	CreatedAt         time.Time     `json:"created_at"`
	CreatedById       sql.NullInt64         `json:"created_by_id,omitempty"`
	// CreatedFrom       sql.NullString        `json:"created_from"`
	LastModifiedAt    time.Time             `json:"last_modified_at"`
	LastModifiedById  sql.NullInt64 `json:"last_modified_by_id,omitempty"`
	// LastModifiedFrom  sql.NullString        `json:"last_modified_from"`
	IsArchived        bool      `json:"is_archived"`
}

func ListAllWorkOrderServiceFees(db *sql.DB) ([]*OldUWorkOrderServiceFee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, title, description, percentage, created_at, created_by_id, last_modified_at, last_modified_by_id, is_archived
	FROM
	    workery_work_order_service_fees
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUWorkOrderServiceFee
	defer rows.Close()
	for rows.Next() {
		m := new(OldUWorkOrderServiceFee)
		err = rows.Scan(
			&m.Id,
			&m.Title,
			&m.Description,
			&m.Percentage,
			&m.CreatedAt,
			&m.CreatedById,
			&m.LastModifiedAt,
			&m.LastModifiedById,
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

func runWorkOrderServiceFeeETL(ctx context.Context, tenantId uint64, irr *repositories.WorkOrderServiceFeeRepo, oldDb *sql.DB) {
	workOrderServiceFees, err := ListAllWorkOrderServiceFees(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range workOrderServiceFees {
		insertWorkOrderServiceFeeETL(ctx, tenantId, irr, oir)
	}
}

func insertWorkOrderServiceFeeETL(ctx context.Context, tid uint64, irr *repositories.WorkOrderServiceFeeRepo, oir *OldUWorkOrderServiceFee) {
	var state int8 = 1
	if oir.IsArchived == true {
		state = 0
	}

	m := &models.WorkOrderServiceFee{
		OldId: oir.Id,
		TenantId: tid,
		Uuid: uuid.NewString(),
		Title: oir.Title,
		Description: oir.Description,
		Percentage: oir.Percentage,
		CreatedTime: oir.CreatedAt,
		// CreatedById: oir.CreatedById,
		LastModifiedTime: oir.LastModifiedAt,
		LastModifiedById: oir.LastModifiedById,
		State: state,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oir.Id)
}
