package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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
	Short: "Import the work order from old workery",
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
	isfr := repositories.NewWorkOrderServiceFeeRepo(db)
	owor := repositories.NewOngoingWorkOrderRepo(db)

	runWorkOrderETL(ctx, uint64(woETLTenantId), asr, ar, cr, isfr, owor, oldDb)
}

func runWorkOrderETL(
	ctx context.Context,
	tenantId uint64,
	asr *repositories.WorkOrderRepo,
	ar *repositories.AssociateRepo,
	cr *repositories.CustomerRepo,
	isfr *repositories.WorkOrderServiceFeeRepo,
	owor *repositories.OngoingWorkOrderRepo,
	oldDb *sql.DB,
) {
	ass, err := ListAllWorkOrders(oldDb)
	if err != nil {
		log.Fatal("ListAllWorkOrders", err)
	}
	for _, oss := range ass {
		insertWorkOrderETL(ctx, tenantId, asr, ar, cr, isfr, owor, oss)
	}
}

type OldWorkOrder struct {
	Id                                        uint64      `json:"id"`
	AssociateId                               null.Int    `json:"associate_id"`
	CustomerId                                uint64      `json:"customer_id"`
	Description                               string      `json:"description"`
	AssignmentDate                            null.Time   `json:"assignment_date"`
	IsOngoing                                 bool        `json:"is_ongoing"`
	IsHomeSupportService                      bool        `json:"is_home_support_service"`
	StartDate                                 time.Time   `json:"start_date"`
	CompletionDate                            null.Time   `json:"completion_date"`
	Hours                                     float64     `json:"hours"`
	TypeOf                                    int8        `json:"type_of"`
	IndexedText                               string      `json:"indexed_text"`
	ClosingReason                             int8        `json:"closing_reason"`
	ClosingReasonOther                        null.String `json:"closing_reason_other"`
	State                                     string      `json:"state"`
	WasJobSatisfactory                        bool        `json:"was_job_satisfactory"`
	WasJobFinishedOnTimeAndOnBudget           bool        `json:"was_job_finished_on_time_and_on_budget"`
	WasAssociatePunctual                      bool        `json:"was_associate_punctual"`
	WasAssociateProfessional                  bool        `json:"was_associate_professional"`
	WouldCustomerReferOurOrganization         bool        `json:"would_customer_refer_our_organization"`
	Score                                     int8        `json:"score"`
	InvoiceDate                               null.Time   `json:"invoice_date"`
	InvoiceQuoteAmountCurrency                string      `json:"invoice_quote_amount_currency"`
	InvoiceQuoteAmount                        float64     `json:"invoice_quote_amount"`
	InvoiceLabourAmountCurrency               string      `json:"invoice_labour_amount_currency"`
	InvoiceLabourAmount                       float64     `json:"invoice_labour_amount"`
	InvoiceMaterialAmountCurrency             string      `json:"invoice_material_amount_currency"`
	InvoiceMaterialAmount                     float64     `json:"invoice_material_amount"`
	InvoiceTaxAmountCurrency                  string      `json:"invoice_tax_amount_currency"`
	InvoiceTaxAmount                          float64     `json:"invoice_tax_amount"`
	InvoiceTotalAmountCurrency                string      `json:"invoice_total_amount_currency"`
	InvoiceTotalAmount                        float64     `json:"invoice_total_amount"`
	InvoiceServiceFeeAmountCurrency           string      `json:"invoice_service_fee_amount_currency"`
	InvoiceServiceFeeAmount                   float64     `json:"invoice_service_fee_amount"`
	InvoiceServiceFeePaymentDate              null.Time   `json:"invoice_service_fee_payment_date"`
	Created                                   time.Time   `json:"created"`
	CreatedById                               null.Int    `json:"created_by_id"`
	CreatedFrom                               null.String `json:"created_from"`
	LastModified                              time.Time   `json:"last_modified"`
	LastModifiedById                          null.Int    `json:"last_modified_by_id"`
	LastModifiedFrom                          null.String `json:"last_modified_from"`
	InvoiceServiceFeeId                       null.Int    `json:"invoice_service_fee_id"`
	LatestPendingTaskId                       null.Int    `json:"latest_pending_task_id"`
	OngoingWorkOrderId                        null.Int    `json:"ongoing_work_order_id"`
	WasSurveyConducted                        bool        `json:"was_survey_conducted"`
	WasThereFinancialsInputted                bool        `json:"was_there_financials_inputted"`
	InvoiceActualServiceFeeAmountPaidCurrency string      `json:"invoice_actual_service_fee_amount_paid_currency"`
	InvoiceActualServiceFeeAmountPaid         float64     `json:"invoice_actual_service_fee_amount_paid"`
	InvoiceBalanceOwingAmountCurrency         string      `json:"invoice_balance_owing_amount_currency"`
	InvoiceBalanceOwingAmount                 float64     `json:"invoice_balance_owing_amount"`
	InvoiceQuotedLabourAmountCurrency         string      `json:"invoice_quoted_labour_amount_currency"`
	InvoiceQuotedLabourAmount                 float64     `json:"invoice_quoted_labour_amount"`
	InvoiceQuotedMaterialAmountCurrency       string      `json:"invoice_quoted_material_amount_currency"`
	InvoiceQuotedMaterialAmount               float64     `json:"invoice_quoted_material_amount"`
	InvoiceTotalQuoteAmountCurrency           string      `json:"invoice_total_quote_amount_currency"`
	InvoiceTotalQuoteAmount                   float64     `json:"invoice_total_quote_amount"`
	Visits                                    int8        `json:"visits"`
	InvoiceIds                                null.String `json:"invoice_ids"`
	NoSurveyConductedReason                   null.Int    `json:"no_survey_conducted_reason"`
	NoSurveyConductedReasonOther              null.String `json:"no_survey_conducted_reason_other"`
	ClonedFromId                              null.Int    `json:"cloned_from_id"`
	InvoiceDepositAmountCurrency              string      `json:"invoice_deposit_amount_currency"`
	InvoiceDepositAmount                      float64     `json:"invoice_deposit_amount"`
	InvoiceOtherCostsAmountCurrency           string      `json:"invoice_other_costs_amount_currency"`
	InvoiceOtherCostsAmount                   float64     `json:"invoice_other_costs_amount"`
	InvoiceQuotedOtherCostsAmountCurrency     string      `json:"invoice_quoted_other_costs_amount_currency"`
	InvoiceQuotedOtherCostsAmount             float64     `json:"invoice_quoted_other_costs_amount"`
	InvoicePaidTo                             null.Int    `json:"invoice_paid_to"`
	InvoiceAmountDueCurrency                  string      `json:"invoice_amount_due_currency"`
	InvoiceAmountDue                          float64     `json:"invoice_amount_due"`
	InvoiceSubTotalAmountCurrency             string      `json:"invoice_sub_total_amount_currency"`
	InvoiceSubTotalAmount                     float64     `json:"invoice_sub_total_amount"`
	ClosingReasonComment                      string      `json:"closing_reason_comment"`
}

func ListAllWorkOrders(db *sql.DB) ([]*OldWorkOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        id, associate_id, customer_id, description, assignment_date, is_ongoing, is_home_support_service, start_date,
		completion_date, hours, type_of, indexed_text, closing_reason, closing_reason_other, state,
		was_job_satisfactory, was_job_finished_on_time_and_on_budget, was_associate_punctual, was_associate_professional, would_customer_refer_our_organization,
		score, invoice_date, invoice_quote_amount_currency, invoice_quote_amount, invoice_labour_amount_currency, invoice_labour_amount,
		invoice_material_amount_currency, invoice_material_amount, invoice_tax_amount_currency, invoice_tax_amount,
		invoice_total_amount_currency, invoice_total_amount, invoice_service_fee_amount_currency, invoice_service_fee_amount, invoice_service_fee_payment_date,
		created, created_by_id, created_from, last_modified, last_modified_by_id, last_modified_from, invoice_service_fee_id, latest_pending_task_id, ongoing_work_order_id,
		was_survey_conducted, was_there_financials_inputted, invoice_actual_service_fee_amount_paid_currency, invoice_actual_service_fee_amount_paid,
		invoice_balance_owing_amount_currency, invoice_balance_owing_amount, invoice_quoted_labour_amount_currency, invoice_quoted_labour_amount,
		invoice_quoted_material_amount_currency, invoice_quoted_material_amount, invoice_total_quote_amount_currency, invoice_total_quote_amount, visits, invoice_ids,
		no_survey_conducted_reason, no_survey_conducted_reason_other, cloned_from_id, invoice_deposit_amount_currency, invoice_deposit_amount,
		invoice_other_costs_amount_currency, invoice_other_costs_amount, invoice_quoted_other_costs_amount_currency, invoice_quoted_other_costs_amount, invoice_paid_to,
		invoice_amount_due_currency, invoice_amount_due, invoice_sub_total_amount_currency, invoice_sub_total_amount, closing_reason_comment
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
			&m.Id, &m.AssociateId, &m.CustomerId, &m.Description, &m.AssignmentDate, &m.IsOngoing, &m.IsHomeSupportService, &m.StartDate,
			&m.CompletionDate, &m.Hours, &m.TypeOf, &m.IndexedText, &m.ClosingReason, &m.ClosingReasonOther, &m.State,
			&m.WasJobSatisfactory, &m.WasJobFinishedOnTimeAndOnBudget, &m.WasAssociatePunctual, &m.WasAssociateProfessional, &m.WouldCustomerReferOurOrganization,
			&m.Score, &m.InvoiceDate, &m.InvoiceQuoteAmountCurrency, &m.InvoiceQuoteAmount, &m.InvoiceLabourAmountCurrency, &m.InvoiceLabourAmount,
			&m.InvoiceMaterialAmountCurrency, &m.InvoiceMaterialAmount, &m.InvoiceTaxAmountCurrency, &m.InvoiceTaxAmount,
			&m.InvoiceTotalAmountCurrency, &m.InvoiceTotalAmount, &m.InvoiceServiceFeeAmountCurrency, &m.InvoiceServiceFeeAmount, &m.InvoiceServiceFeePaymentDate,
			&m.Created, &m.CreatedById, &m.CreatedFrom, &m.LastModified, &m.LastModifiedById, &m.LastModifiedFrom, &m.InvoiceServiceFeeId, &m.LatestPendingTaskId, &m.OngoingWorkOrderId,
			&m.WasSurveyConducted, &m.WasThereFinancialsInputted, &m.InvoiceActualServiceFeeAmountPaidCurrency, &m.InvoiceActualServiceFeeAmountPaid,
			&m.InvoiceBalanceOwingAmountCurrency, &m.InvoiceBalanceOwingAmount, &m.InvoiceQuotedLabourAmountCurrency, &m.InvoiceQuotedLabourAmount,
			&m.InvoiceQuotedMaterialAmountCurrency, &m.InvoiceQuotedMaterialAmount, &m.InvoiceTotalQuoteAmountCurrency, &m.InvoiceTotalQuoteAmount, &m.Visits, &m.InvoiceIds,
			&m.NoSurveyConductedReason, &m.NoSurveyConductedReasonOther, &m.ClonedFromId, &m.InvoiceDepositAmountCurrency, &m.InvoiceDepositAmount,
			&m.InvoiceOtherCostsAmountCurrency, &m.InvoiceOtherCostsAmount, &m.InvoiceQuotedOtherCostsAmountCurrency, &m.InvoiceQuotedOtherCostsAmount, &m.InvoicePaidTo,
			&m.InvoiceAmountDueCurrency, &m.InvoiceAmountDue, &m.InvoiceSubTotalAmountCurrency, &m.InvoiceSubTotalAmount, &m.ClosingReasonComment,
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
	isfr *repositories.WorkOrderServiceFeeRepo,
	owor *repositories.OngoingWorkOrderRepo,
	oss *OldWorkOrder,
) {
	var associateId null.Int
	var associateName null.String
	var associateLexicalName null.String
	if oss.AssociateId.Valid {
		associateIdInt64 := oss.AssociateId.ValueOrZero()
		associateIdUint64, err := ar.GetIdByOldId(ctx, tid, uint64(associateIdInt64))
		if err != nil {
			log.Panic("ar.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		associateId = null.NewInt(int64(associateIdUint64), associateIdUint64 != 0)

		if associateIdUint64 != 0 {
			associate, err := ar.GetById(ctx, associateIdUint64)
			if err != nil {
				log.Panic("ar.GetById | err", err)
			}

			// Generate our full name / lexical full name.
			if associate.MiddleName != "" {
				associateNameStr := associate.GivenName + " " + associate.MiddleName + " " + associate.LastName
				associateLexicalNameStr := associate.LastName + ", " + associate.MiddleName + ", " + associate.GivenName
				associateLexicalNameStr = strings.Replace(associateLexicalNameStr, ", ,", ",", 0)
				associateLexicalNameStr = strings.Replace(associateLexicalNameStr, "  ", " ", 0)

				associateName = null.NewString(associateNameStr, true)
				associateLexicalName = null.NewString(associateLexicalNameStr, true)
			} else {
				associateNameStr := associate.GivenName + " " + associate.LastName
				associateLexicalNameStr := associate.LastName + ", " + associate.GivenName
				associateLexicalNameStr = strings.Replace(associateLexicalNameStr, ", ,", ",", 0)
				associateLexicalNameStr = strings.Replace(associateLexicalNameStr, "  ", " ", 0)

				associateName = null.NewString(associateNameStr, true)
				associateLexicalName = null.NewString(associateLexicalNameStr, true)
			}
		}
	}

	customerId, err := cr.GetIdByOldId(ctx, tid, oss.CustomerId)
	if err != nil {
		log.Panic("cr.GetIdByOldId | err", err)
	}

	// Lookup our customer record so we can generate the full name / lexical full name.
	customer, err := cr.GetById(ctx, customerId)

	// Generate our full name / lexical full name.
	var customerName string
	var customerLexicalName string
	if customer.MiddleName != "" {
		customerName = customer.GivenName + " " + customer.MiddleName + " " + customer.LastName
		customerLexicalName = customer.LastName + ", " + customer.MiddleName + ", " + customer.GivenName
	} else {
		customerName = customer.GivenName + " " + customer.LastName
		customerLexicalName = customer.LastName + ", " + customer.GivenName
	}
	customerLexicalName = strings.Replace(customerLexicalName, ", ,", ",", 0)
	customerLexicalName = strings.Replace(customerLexicalName, "  ", " ", 0)

	var state int8
	switch s := oss.State; s {
	case "new":
		state = models.WorkOrderNewState
	case "declined":
		state = models.WorkOrderDeclinedState
	case "pending":
		state = models.WorkOrderPendingState
	case "cancelled":
		state = models.WorkOrderCancelledState
	case "ongoing":
		state = models.WorkOrderOngoingState
	case "in_progress":
		state = models.WorkOrderInProgressState
	case "completed_and_unpaid":
		state = models.WorkOrderCompletedButUnpaidState
	case "completed_but_unpaid":
		state = models.WorkOrderCompletedButUnpaidState
	case "completed_and_paid":
		state = models.WorkOrderCompletedAndPaidState
	case "archived":
		state = models.WorkOrderArchivedState
	default:
		state = models.WorkOrderArchivedState
	}

	var invoiceServiceFeeId null.Int
	if oss.InvoiceServiceFeeId.Valid {
		invoiceServiceFeeIdInt64 := oss.InvoiceServiceFeeId.ValueOrZero()
		invoiceServiceFeeIdUint64, err := isfr.GetIdByOldId(ctx, tid, uint64(invoiceServiceFeeIdInt64))
		if err != nil {
			log.Panic("isfr.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		invoiceServiceFeeId = null.NewInt(int64(invoiceServiceFeeIdUint64), invoiceServiceFeeIdUint64 != 0)
	}

	var ongoingWorkOrderId null.Int
	if oss.OngoingWorkOrderId.Valid {
		ongoingWorkOrderIdInt64 := oss.OngoingWorkOrderId.ValueOrZero()
		ongoingWorkOrderIdUint64, err := owor.GetIdByOldId(ctx, tid, uint64(ongoingWorkOrderIdInt64))
		if err != nil {
			log.Panic("owor.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		ongoingWorkOrderId = null.NewInt(int64(ongoingWorkOrderIdUint64), ongoingWorkOrderIdUint64 != 0)
	}

	var clonedFromId null.Int
	if oss.ClonedFromId.Valid {
		clonedFromIdInt64 := oss.ClonedFromId.ValueOrZero()
		clonedFromIdUint64, err := asr.GetIdByOldId(ctx, tid, uint64(clonedFromIdInt64))
		if err != nil {
			log.Panic("asr.GetIdByOldId | err", err)
		}

		// Convert from null supported integer times.
		clonedFromId = null.NewInt(int64(clonedFromIdUint64), clonedFromIdUint64 != 0)
	}

	// InvoicePaidTo

	m := &models.WorkOrder{
		OldId:                             oss.Id,
		TenantId:                          tid,
		Uuid:                              uuid.NewString(),
		CustomerId:                        customerId,
		CustomerName:                      customerName,
		CustomerLexicalName:               customerLexicalName,
		AssociateId:                       associateId,
		AssociateName:                     associateName,
		AssociateLexicalName:              associateLexicalName,
		Description:                       oss.Description,
		AssignmentDate:                    oss.AssignmentDate,
		IsOngoing:                         oss.IsOngoing,
		IsHomeSupportService:              oss.IsHomeSupportService,
		StartDate:                         oss.StartDate,
		CompletionDate:                    oss.CompletionDate,
		Hours:                             oss.Hours,
		TypeOf:                            oss.TypeOf,
		IndexedText:                       oss.IndexedText,
		ClosingReason:                     oss.ClosingReason,
		ClosingReasonOther:                oss.ClosingReasonOther,
		State:                             state,
		Currency:                          "CAD",
		WasJobSatisfactory:                oss.WasJobSatisfactory,
		WasJobFinishedOnTimeAndOnBudget:   oss.WasJobFinishedOnTimeAndOnBudget,
		WasAssociatePunctual:              oss.WasAssociatePunctual,
		WasAssociateProfessional:          oss.WasAssociateProfessional,
		WouldCustomerReferOurOrganization: oss.WouldCustomerReferOurOrganization,
		Score:                             oss.Score,
		InvoiceDate:                       oss.InvoiceDate,
		InvoiceQuoteAmount:                oss.InvoiceQuoteAmount,
		InvoiceLabourAmount:               oss.InvoiceLabourAmount,
		InvoiceMaterialAmount:             oss.InvoiceMaterialAmount,
		InvoiceTaxAmount:                  oss.InvoiceTaxAmount,
		InvoiceTotalAmount:                oss.InvoiceTotalAmount,
		InvoiceServiceFeeAmount:           oss.InvoiceServiceFeeAmount,
		InvoiceServiceFeePaymentDate:      oss.InvoiceServiceFeePaymentDate,
		CreatedTime:                       oss.Created,
		CreatedById:                       oss.CreatedById,
		CreatedFromIP:                     oss.CreatedFrom,
		LastModifiedTime:                  oss.LastModified,
		LastModifiedById:                  oss.LastModifiedById,
		LastModifiedFromIP:                oss.LastModifiedFrom,
		InvoiceServiceFeeId:               invoiceServiceFeeId,
		LatestPendingTaskId:               oss.LatestPendingTaskId,
		OngoingWorkOrderId:                ongoingWorkOrderId,
		WasSurveyConducted:                oss.WasSurveyConducted,
		WasThereFinancialsInputted:        oss.WasThereFinancialsInputted,
		InvoiceActualServiceFeeAmountPaid: oss.InvoiceActualServiceFeeAmountPaid,
		InvoiceBalanceOwingAmount:         oss.InvoiceBalanceOwingAmount,
		InvoiceQuotedLabourAmount:         oss.InvoiceQuotedLabourAmount,
		InvoiceQuotedMaterialAmount:       oss.InvoiceQuotedMaterialAmount,
		InvoiceTotalQuoteAmount:           oss.InvoiceTotalQuoteAmount,
		Visits:                            oss.Visits,
		InvoiceIds:                        oss.InvoiceIds,
		NoSurveyConductedReason:           oss.NoSurveyConductedReason,
		NoSurveyConductedReasonOther:      oss.NoSurveyConductedReasonOther,
		ClonedFromId:                      clonedFromId,
		InvoiceDepositAmount:              oss.InvoiceDepositAmount,
		InvoiceOtherCostsAmount:           oss.InvoiceOtherCostsAmount,
		InvoiceQuotedOtherCostsAmount:     oss.InvoiceQuotedOtherCostsAmount,
		InvoicePaidTo:                     oss.InvoicePaidTo,
		InvoiceAmountDue:                  oss.InvoiceAmountDue,
		InvoiceSubTotalAmount:             oss.InvoiceSubTotalAmount,
		ClosingReasonComment:              oss.ClosingReasonComment,
	}

	// // For debugging purposes only.
	// log.Println("associateId -->", associateId)
	// log.Println("customerId  -->", customerId)
	// log.Println("State       -->", state)
	// log.Println("Model       -->", m)

	err = asr.Insert(ctx, m)
	if err != nil {
		// log.Print("associateId", associateId)
		// log.Print("customerId", customerId)
		// log.Panic("asr.Insert | err", err, "\n\n", m, oss)
		log.Panic("asr.Insert | err", err, "\n\n")
	}
	fmt.Println("Imported ID#", oss.Id)
}
