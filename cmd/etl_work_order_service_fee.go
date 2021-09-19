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
	ur := repositories.NewUserRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, workOrderServiceFeeETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runWorkOrderServiceFeeETL(ctx, tenant.Id, irr, ur, oldDb)
}

type OldUWorkOrderServiceFee struct {
	Id               uint64    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Percentage       float64   `json:"percentage"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedById      null.Int  `json:"created_by_id,omitempty"`
	LastModifiedAt   time.Time `json:"last_modified_at"`
	LastModifiedById null.Int  `json:"last_modified_by_id,omitempty"`
	IsArchived       bool      `json:"is_archived"`
}

func ListAllWorkOrderServiceFees(db *sql.DB) ([]*OldUWorkOrderServiceFee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, title, description, percentage, created_at, created_by_id, last_modified_at, last_modified_by_id, is_archived
	FROM
	    workery_work_order_service_fees
	ORDER BY
	    id
	ASC
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

func runWorkOrderServiceFeeETL(ctx context.Context, tenantId uint64, irr *repositories.WorkOrderServiceFeeRepo, ur *repositories.UserRepo, oldDb *sql.DB) {
	workOrderServiceFees, err := ListAllWorkOrderServiceFees(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range workOrderServiceFees {
		insertWorkOrderServiceFeeETL(ctx, tenantId, irr, ur, oir)
	}
}

func insertWorkOrderServiceFeeETL(ctx context.Context, tid uint64, irr *repositories.WorkOrderServiceFeeRepo, ur *repositories.UserRepo, oir *OldUWorkOrderServiceFee) {
	//
	// Set the `state`.
	//

	var state int8 = 1
	if oir.IsArchived == true {
		state = 0
	}

	//
	// Get `createdById` and `createdByName` values.
	//

	var createdById null.Int
	var createdByName null.String
	if oir.CreatedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(oir.CreatedById.ValueOrZero()))

		if err != nil {
			log.Panic("ur.GetIdByOldId", err)
		}
		user, err := ur.GetById(ctx, userId)
		if err != nil {
			log.Panic("ur.GetById", err)
		}

		if user != nil {
			createdById = null.IntFrom(int64(userId))
			createdByName = null.StringFrom(user.Name)
		} else {
			log.Println("WARNING: D.N.E.")
		}

		// // For debugging purposes only.
		// log.Println("createdById:", createdById)
		// log.Println("createdByName:", createdByName)
	}

	//
	// Get `lastModifiedById` and `lastModifiedByName` values.
	//

	var lastModifiedById null.Int
	var lastModifiedByName null.String
	if oir.LastModifiedById.ValueOrZero() > 0 {
		userId, err := ur.GetIdByOldId(ctx, tid, uint64(oir.LastModifiedById.ValueOrZero()))
		if err != nil {
			log.Panic("ur.GetIdByOldId", err)
		}
		user, err := ur.GetById(ctx, userId)
		if err != nil {
			log.Panic("ur.GetById", err)
		}

		if user != nil {
			lastModifiedById = null.IntFrom(int64(userId))
			lastModifiedByName = null.StringFrom(user.Name)
		} else {
			log.Println("WARNING: D.N.E.")
		}

		// // For debugging purposes only.
		// log.Println("lastModifiedById:", lastModifiedById)
		// log.Println("lastModifiedByName:", lastModifiedByName)
	}

	//
	// Insert the `WorkOrderServiceFee`.
	//

	m := &models.WorkOrderServiceFee{
		OldId:              oir.Id,
		TenantId:           tid,
		Uuid:               uuid.NewString(),
		Title:              oir.Title,
		Description:        oir.Description,
		Percentage:         oir.Percentage,
		CreatedTime:        oir.CreatedAt,
		CreatedById:        createdById,
		CreatedByName:      createdByName,
		LastModifiedTime:   oir.LastModifiedAt,
		LastModifiedById:   lastModifiedById,
		LastModifiedByName: lastModifiedByName,
		State:              state,
	}
	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oir.Id)
}
