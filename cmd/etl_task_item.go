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
	taskItemETLSchemaName string
)

func init() {
	taskItemETLCmd.Flags().StringVarP(&taskItemETLSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	taskItemETLCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(taskItemETLCmd)
}

var taskItemETLCmd = &cobra.Command{
	Use:   "etl_task_item",
	Short: "Import the taskItem data from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportTaskItem()
	},
}

func doRunImportTaskItem() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, taskItemETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	tir := repositories.NewTaskItemRepo(db)
	wor := repositories.NewWorkOrderRepo(db)
	owor := repositories.NewOngoingWorkOrderRepo(db)
	ur := repositories.NewUserRepo(db)
	ar := repositories.NewAssociateRepo(db)
	cr := repositories.NewCustomerRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, taskItemETLSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	if tenant == nil {
		log.Fatal("Tenant does not exist!")
	}

	runTaskItemETL(ctx, tenant.Id, tir, wor, owor, ur, ar, cr, oldDb)
}

type OldUTaskItem struct {
	Id                       uint64      `json:"id"`
	TypeOf                   string      `json:"type_of"`
	Title                    string      `json:"title"`
	Description              string      `json:"description"`
	DueDate                  time.Time   `json:"due_date"`
	IsClosed                 bool        `json:"is_closed"`
	WasPostponed             string      `json:"was_postponed"`
	ClosingReason            int8        `json:"closing_reason"`
	ClosingReasonOther       string      `json:"closing_reason_other"`
	CreatedAt                time.Time   `json:"created_at"`
	CreatedFrom              null.String `json:"created_from"`
	CreatedFromIsPublic      null.Bool   `json:"created_from_is_public"`
	CreatedById              null.Int    `json:"created_by_id"`
	LastModifiedAt           time.Time   `json:"last_modified_at"`
	LastModifiedFrom         null.String `json:"last_modified_from"`
	LastModifiedFromIsPublic null.Bool   `json:"last_modified_from_is_public"`
	JobId                    uint64      `json:"job_id"`
	LastModifiedById         null.Int    `json:"last_modified_by_id"`
	OngoingJobId             null.Int    `json:"ongoing_job_id"`
}

/*
    integer,
   bigint,
*/

func ListAllTaskItems(db *sql.DB) ([]*OldUTaskItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
	    id, type_of, title, description, due_date, is_closed, was_postponed,
		closing_reason, closing_reason_other, created_at, created_from,
		created_from_is_public, last_modified_at, last_modified_from,
		last_modified_from_is_public, created_by_id, job_id, last_modified_by_id,
		ongoing_job_id
	FROM
	    workery_task_items
	ORDER BY
	    id
	ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldUTaskItem
	defer rows.Close()
	for rows.Next() {
		m := new(OldUTaskItem)
		err = rows.Scan(
			&m.Id,
			&m.TypeOf,
			&m.Title,
			&m.Description,
			&m.DueDate,
			&m.IsClosed,
			&m.WasPostponed,
			&m.ClosingReason,
			&m.ClosingReasonOther,
			&m.CreatedAt,
			&m.CreatedFrom,
			&m.CreatedFromIsPublic,
			&m.LastModifiedAt,
			&m.LastModifiedFrom,
			&m.LastModifiedFromIsPublic,
			&m.CreatedById,
			&m.JobId,
			&m.LastModifiedById,
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

func runTaskItemETL(
	ctx context.Context,
	tenantId uint64,
	tir *repositories.TaskItemRepo,
	wor *repositories.WorkOrderRepo,
	owor *repositories.OngoingWorkOrderRepo,
	ur *repositories.UserRepo,
	ar *repositories.AssociateRepo,
	cr *repositories.CustomerRepo,
	oldDb *sql.DB,
) {
	oldTaskItems, err := ListAllTaskItems(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oti := range oldTaskItems {
		insertTaskItemETL(ctx, tenantId, tir, wor, owor, ur, ar, cr, oti)
	}
}

func insertTaskItemETL(
	ctx context.Context,
	tid uint64,
	tir *repositories.TaskItemRepo,
	wor *repositories.WorkOrderRepo,
	owor *repositories.OngoingWorkOrderRepo,
	ur *repositories.UserRepo,
	ar *repositories.AssociateRepo,
	cr *repositories.CustomerRepo,
	oti *OldUTaskItem,
) {
	var state int8 = 1

	// log.Println("OngoingJobId --->>", oti.OngoingJobId)
	// log.Println("OrderId --->>", oti.JobId)
	// log.Println("CreatedById --->>", oti.CreatedById)
	// log.Println("LastModifiedById --->>", oti.LastModifiedById)

	orderId, _ := wor.GetIdByOldId(ctx, tid, oti.JobId)
	order, _ := wor.GetById(ctx, orderId)

	var createdById null.Int
	if oti.CreatedById.Valid {
		val := oti.CreatedById.ValueOrZero()
		id, _ := ur.GetIdByOldId(ctx, tid, uint64(val))
		createdById = null.IntFrom(int64(id))
		// log.Println("ID:", oir.Id, "|User|IN:", oir.CreatedById, "OUT:", createdById, "\tTenantId:", tid)
	}
	var lastModifiedById null.Int
	if oti.LastModifiedById.Valid {
		val := oti.LastModifiedById.ValueOrZero()
		id, _ := ur.GetIdByOldId(ctx, tid, uint64(val))
		lastModifiedById = null.IntFrom(int64(id))
	}

	var ongoingJobId null.Int
	if oti.OngoingJobId.Valid {
		val := oti.OngoingJobId.ValueOrZero()
		id, _ := owor.GetIdByOldId(ctx, tid, uint64(val))
		ongoingJobId = null.IntFrom(int64(id))
	}

	m := &models.TaskItem{
		TenantId:             tid,
		Uuid:                 uuid.NewString(),
		TypeOf:               oti.TypeOf,
		Title:                oti.Title,
		Description:          oti.Description,
		DueDate:              oti.DueDate,
		IsClosed:             oti.IsClosed,
		WasPostponed:         oti.WasPostponed,
		ClosingReason:        oti.ClosingReason,
		ClosingReasonOther:   oti.ClosingReasonOther,
		OrderId:              orderId,
		OrderTypeOf:          order.TypeOf,
		OngoingOrderId:       ongoingJobId,
		CreatedTime:          oti.CreatedAt,
		CreatedFromIP:        oti.CreatedFrom,
		CreatedById:          createdById,
		LastModifiedTime:     oti.LastModifiedAt,
		LastModifiedFromIP:   oti.LastModifiedFrom,
		LastModifiedById:     lastModifiedById,
		State:                state,
		OldId:                oti.Id,
		AssociateId:          order.AssociateId,
		AssociateName:        order.AssociateName,
		AssociateLexicalName: order.AssociateLexicalName,
		CustomerId:           null.NewInt(int64(order.CustomerId), order.CustomerId != 0),
		CustomerName:         null.NewString(order.CustomerName, order.CustomerName != ""),
		CustomerLexicalName:  null.NewString(order.CustomerLexicalName, order.CustomerLexicalName != ""),
	}
	err := tir.Insert(ctx, m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Imported ID#", oti.Id)
}
