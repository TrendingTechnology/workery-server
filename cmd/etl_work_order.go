package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/google/uuid"
	"github.com/spf13/cobra"
	null "gopkg.in/guregu/null.v4"

	// "github.com/over55/workery-server/internal/models"
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
		fmt.Println(oss, "\n")
		// insertWorkOrderETL(ctx, tenantId, asr, ar, cr, oss)
	}
}

// '''
//     was_job_satisfactory boolean NOT NULL,
//     was_job_finished_on_time_and_on_budget boolean NOT NULL,
//     was_associate_punctual boolean NOT NULL,
//     was_associate_professional boolean NOT NULL,
//     would_customer_refer_our_organization boolean NOT NULL,
//     score smallint NOT NULL,
//     invoice_date date,
//     invoice_quote_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quote_amount numeric(10,2) NOT NULL,
//     invoice_labour_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_labour_amount numeric(10,2) NOT NULL,
//     invoice_material_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_material_amount numeric(10,2) NOT NULL,
//     invoice_tax_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_tax_amount numeric(10,2) NOT NULL,
//     invoice_total_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_total_amount numeric(10,2) NOT NULL,
//     invoice_service_fee_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_service_fee_amount numeric(10,2) NOT NULL,
//     invoice_service_fee_payment_date date,
//     created timestamp with time zone NOT NULL,
//     created_from inet,
//     created_from_is_public boolean NOT NULL,
//     last_modified timestamp with time zone NOT NULL,
//     last_modified_from inet,
//     last_modified_from_is_public boolean NOT NULL,
//     associate_id bigint,
//     created_by_id integer,
//     customer_id bigint NOT NULL,
//     invoice_service_fee_id bigint,
//     last_modified_by_id integer,
//     latest_pending_task_id bigint,
//     ongoing_work_order_id bigint,
//     was_survey_conducted boolean NOT NULL,
//     was_there_financials_inputted boolean NOT NULL,
//     invoice_actual_service_fee_amount_paid numeric(10,2) NOT NULL,
//     invoice_actual_service_fee_amount_paid_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_balance_owing_amount numeric(10,2) NOT NULL,
//     invoice_balance_owing_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quoted_labour_amount numeric(10,2) NOT NULL,
//     invoice_quoted_labour_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quoted_material_amount numeric(10,2) NOT NULL,
//     invoice_quoted_material_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_total_quote_amount numeric(10,2) NOT NULL,
//     invoice_total_quote_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     visits smallint NOT NULL,
//     invoice_ids character varying(127) COLLATE pg_catalog."default",
//     no_survey_conducted_reason smallint,
//     no_survey_conducted_reason_other character varying(1024) COLLATE pg_catalog."default",
//     cloned_from_id bigint,
//     invoice_deposit_amount numeric(10,2) NOT NULL,
//     invoice_deposit_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_other_costs_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_quoted_other_costs_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_paid_to smallint,
//     invoice_amount_due numeric(10,2) NOT NULL,
//     invoice_amount_due_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_sub_total_amount numeric(10,2) NOT NULL,
//     invoice_sub_total_amount_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
//     invoice_other_costs_amount numeric(10,2) NOT NULL,
//     invoice_quoted_other_costs_amount numeric(10,2) NOT NULL,
//     closing_reason_comment character varying(1024) COLLATE pg_catalog."default",
// '''

type OldWorkOrder struct {
	Id               uint64      `json:"id"`
	AssociateId      null.Int    `json:"associate_id"`
	CustomerId       uint64      `json:"customer_id"`

	Description      string      `json:"description"`
    AssignmentDate   null.Time `json:"assignment_date"`
	IsOngoing        bool      `json:"is_ongoing"`
	IsHomeSupportService bool      `json:"is_home_support_service"`
	StartDate time.Time      `json:"start_date"`
	CompletionDate null.Time      `json:"completion_date"`
	Hours      string      `json:"hours"`
	TypeOf      int8      `json:"type_of"`
	IndexedText      string      `json:"indexed_text"`
	ClosingReason      int8      `json:"closing_reason"`
	ClosingReasonOther null.String     `json:"closing_reason_other"`
	State              string      `json:"state"`

	Created          time.Time   `json:"created"`
	CreatedById      null.Int    `json:"created_by_id"`
	CreatedFrom      null.String `json:"created_from"`
	LastModified     time.Time   `json:"last_modified"`
	LastModifiedById null.Int    `json:"last_modified_by_id"`
	LastModifiedFrom null.String `json:"last_modified_from"`
}

// description text COLLATE pg_catalog."default",
//  date,

func ListAllWorkOrders(db *sql.DB) ([]*OldWorkOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, associate_id, customer_id, description, assignment_date, is_ongoing, is_home_support_service, start_date, completion_date,
		hours, type_of, indexed_text, closing_reason, closing_reason_other, state,
		created, created_by_id, created_from, last_modified, last_modified_by_id, last_modified_from
	FROM
        workery_work_orders
	ORDER BY
	    id
	ASC
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
			&m.Id, &m.AssociateId, &m.CustomerId, &m.Description, &m.AssignmentDate, &m.IsOngoing, &m.IsHomeSupportService, &m.StartDate, &m.CompletionDate,
			&m.Hours, &m.TypeOf, &m.IndexedText, &m.ClosingReason, &m.ClosingReasonOther, &m.State,
			&m.Created, &m.CreatedById, &m.CreatedFrom, &m.LastModified, &m.LastModifiedById, &m.LastModifiedFrom,
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

// func insertWorkOrderETL(
// 	ctx context.Context,
// 	tid uint64,
// 	asr *repositories.WorkOrderRepo,
// 	ar *repositories.AssociateRepo,
// 	cr *repositories.CustomerRepo,
// 	oss *OldWorkOrder,
// ) {
// 	var associateId null.Int
// 	if oss.AssociateId.Valid {
// 		associateIdInt64 := oss.AssociateId.ValueOrZero()
// 		associateIdUint64, err := ar.GetIdByOldId(ctx, tid, uint64(associateIdInt64))
// 		if err != nil {
// 			log.Panic("ar.GetIdByOldId | err", err)
// 		}
//
// 		// Convert from null supported integer times.
// 		associateId = null.NewInt(int64(associateIdUint64), associateIdUint64 != 0)
// 	}
//
// 	customerId, err := cr.GetIdByOldId(ctx, tid, oss.CustomerId)
//
// 	var state int8 = 1 // Running
// 	if oss.State == "terminated" {
// 		state = 2
// 	}
//
// 	m := &models.WorkOrder{
// 		OldId:              oss.Id,
// 		TenantId:           tid,
// 		Uuid:               uuid.NewString(),
// 		CustomerId:         customerId,
// 		AssociateId:        associateId,
// 		State:              state,
// 		CreatedTime:        oss.CreatedAt,
// 		CreatedById:        oss.CreatedById,
// 		CreatedFromIP:      oss.CreatedFrom,
// 		LastModifiedTime:   oss.LastModifiedAt,
// 		LastModifiedById:   oss.LastModifiedById,
// 		LastModifiedFromIP: oss.LastModifiedFrom,
// 	}
// 	err = asr.Insert(ctx, m)
// 	if err != nil {
// 		log.Print("associateId", associateId)
// 		log.Print("customerId", customerId)
// 		log.Panic("asr.Insert | err", err, "\n\n", m, oss)
// 	}
// 	fmt.Println("Imported ID#", oss.Id)
// }
