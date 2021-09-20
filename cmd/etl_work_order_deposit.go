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
	workOrderDepositETLSchemaName string
)

func init() {
	workOrderDepositETLCmd.Flags().StringVarP(&workOrderDepositETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	workOrderDepositETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(workOrderDepositETLCmd)
}

var workOrderDepositETLCmd = &cobra.Command{
	Use:   "etl_work_order_deposit",
	Short: "Import the workOrderDeposit data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportWorkOrderDeposit()
	},
}

func doRunImportWorkOrderDeposit() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, workOrderDepositETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	or := repositories.NewWorkOrderRepo(db)
	irr := repositories.NewWorkOrderDepositRepo(db)
	ur := repositories.NewUserRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, workOrderDepositETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runWorkOrderDepositETL(ctx, tenant.Id, or, irr, ur, oldDb)
}

type OldUWorkOrderDeposit struct {
	Id                       uint64      `json:"id"`
	PaidAt                   null.Time   `json:"paid_at"`
	DepositMethod            int8        `json:"deposit_method"`
	PaidTo                   null.Int    `json:"paid_to"`
	AmountCurrency           string      `json:"amount_currency"`
	Amount                   float64     `json:"amount"`
	PaidFor                  int8        `json:"paid_for"`
	IsArchived               bool        `json:"is_archived"`
	CreatedAt                time.Time   `json:"created_at"`
	LastModifiedAt           time.Time   `json:"last_modified_at"`
	CreatedById              null.Int    `json:"created_by_id"`
	LastModifiedById         null.Int    `json:"last_modified_by_id"`
	OrderId                  uint64      `json:"order_id"`
	CreatedFrom              null.String `json:"created_from"`
	CreatedFromIsPublic      bool        `json:"created_from_is_public"`
	LastModifiedFrom         null.String `json:"last_modified_from"`
	LastModifiedFromIsPublic bool        `json:"last_modified_from_is_public"`
}

/*
 inet,
 boolean NOT NULL,
*/

func ListAllWorkOrderDeposits(db *sql.DB) ([]*OldUWorkOrderDeposit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, paid_at, deposit_method, paid_to, amount_currency, amount, paid_for,
		is_archived, created_at, last_modified_at, created_by_id, last_modified_by_id,
		order_id, created_from, created_from_is_public, last_modified_from,
		last_modified_from_is_public
	FROM
	    workery_work_order_deposits
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUWorkOrderDeposit
	defer rows.Close()
	for rows.Next() {
		m := new(OldUWorkOrderDeposit)
		err = rows.Scan(
			&m.Id,
			&m.PaidAt,
			&m.DepositMethod,
			&m.PaidTo,
			&m.AmountCurrency,
			&m.Amount,
			&m.PaidFor,
			&m.IsArchived,
			&m.CreatedAt,
			&m.LastModifiedAt,
			&m.CreatedById,
			&m.LastModifiedById,
			&m.OrderId,
			&m.CreatedFrom,
			&m.CreatedFromIsPublic,
			&m.LastModifiedFrom,
			&m.LastModifiedFromIsPublic,
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

func runWorkOrderDepositETL(
	ctx context.Context,
	tenantId uint64,
	or *repositories.WorkOrderRepo,
	irr *repositories.WorkOrderDepositRepo,
	ur *repositories.UserRepo,
	oldDb *sql.DB,
) {
	workOrderDeposits, err := ListAllWorkOrderDeposits(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oir := range workOrderDeposits {
		insertWorkOrderDepositETL(ctx, tenantId, or, irr, ur, oir)
	}
}

func insertWorkOrderDepositETL(
	ctx context.Context,
	tid uint64,
	or *repositories.WorkOrderRepo,
	irr *repositories.WorkOrderDepositRepo,
	ur *repositories.UserRepo,
	oir *OldUWorkOrderDeposit,
) {

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

	orderId, _ := or.GetIdByOldId(ctx, tid, oir.OrderId)

	m := &models.WorkOrderDeposit{
		OldId:              oir.Id,
		TenantId:           tid,
		Uuid:               uuid.NewString(),
		PaidAt:             oir.PaidAt,
		DepositMethod:      oir.DepositMethod,
		PaidTo:             oir.PaidTo,
		Currency:           oir.AmountCurrency,
		Amount:             oir.Amount,
		PaidFor:            oir.PaidFor,
		CreatedTime:        oir.CreatedAt,
		LastModifiedTime:   oir.LastModifiedAt,
		CreatedById:        createdById,
		CreatedByName:      createdByName,
		LastModifiedById:   lastModifiedById,
		LastModifiedByName: lastModifiedByName,
		OrderId:            orderId,
		CreatedFromIP:      oir.CreatedFrom,
		LastModifiedFromIP: oir.LastModifiedFrom,
		State:              state,
	}

	err := irr.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", m.OldId)
}
