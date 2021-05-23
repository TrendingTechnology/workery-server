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
	workOrderInvoiceSchemaName string
)

func init() {
	workOrderInvoiceCmd.Flags().StringVarP(&workOrderInvoiceSchemaName, "schema_name", "s", "", "The schema name in the postgres.")
	workOrderInvoiceCmd.MarkFlagRequired("schema_name")
	rootCmd.AddCommand(workOrderInvoiceCmd)
}

var workOrderInvoiceCmd = &cobra.Command{
	Use:   "etl_work_order_invoice",
	Short: "Import the work order invoices from old workery",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doRunImportWorkOrderInvoice()
	},
}

func doRunImportWorkOrderInvoice() {
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
	oldDb, err := utils.ConnectDB(oldDBHost, oldDBPort, oldDBUser, oldDBPassword, oldDBName, workOrderInvoiceSchemaName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldDb.Close()

	// Load up our background context.
	ctx := context.Background()

	// Load up our repositories.
	tr := repositories.NewTenantRepo(db)
	wotp := repositories.NewWorkOrderInvoiceRepo(db)
	ar := repositories.NewWorkOrderRepo(db)

	// Lookup the tenant.
	tenant, err := tr.GetBySchemaName(ctx, workOrderInvoiceSchemaName)
	if err != nil {
		log.Fatal(err)
	}
    if tenant != nil {
		runWorkOrderInvoiceETL(ctx, uint64(tenant.Id), wotp, ar, oldDb)
	}
}

func runWorkOrderInvoiceETL(
	ctx context.Context,
	tenantId uint64,
	wotp *repositories.WorkOrderInvoiceRepo,
	ar *repositories.WorkOrderRepo,
	// vtr *repositories.InvoiceRepo,
	oldDb *sql.DB,
) {
	aats, err := ListAllWorkOrderInvoices(oldDb)
	if err != nil {
		log.Fatal(err)
	}
	for _, oss := range aats {
		insertWorkOrderInvoiceETL(ctx, tenantId, wotp, ar, oss)
	}
}

type OldWorkOrderInvoice struct {
	OrderId   uint64 `json:"order_id"`
	IsArchived bool `json:"is_archived"`
	InvoiceId string `json:"invoice_id"`
	InvoiceDate time.Time `json:"invoice_date"`
	AssociateName string `json:"associate_name"`
	AssociateTelephone string `json:"associate_telephone"`
	ClientName string `json:"client_name"`
	ClientTelephone string `json:"client_telephone"`
	ClientEmail null.String `json:"client_email"`
	Line01Qty int8 `json:"line_01_qty"`
	Line01Desc string `json:"line_01_desc"`
	Line01PriceCurrency string `json:"line_01_price_currency"`
	Line01Price float64 `json:"line_01_price"`
	Line01AmountCurrency string `json:"line_01_amount_currency"`
	Line01Amount float64 `json:"line_01_amount"`
	Line02Qty null.Int `json:"line_02_qty"` // Make `int8`
	Line02Desc string `json:"line_02_desc"`
	Line02PriceCurrency string `json:"line_02_price_currency"`
	Line02Price float64 `json:"line_02_price"`
	Line02AmountCurrency string `json:"line_02_amount_currency"`
	Line02Amount float64 `json:"line_02_amount"`
	Line03Qty null.Int `json:"line_03_qty"` // Make `int8`
	Line03Desc string `json:"line_03_desc"`
	Line03PriceCurrency string `json:"line_03_price_currency"`
	Line03Price float64 `json:"line_03_price"`
	Line03AmountCurrency string `json:"line_03_amount_currency"`
	Line03Amount float64 `json:"line_03_amount"`
	Line04Qty null.Int `json:"line_04_qty"` // Make `int8`
	Line04Desc string `json:"line_04_desc"`
	Line04PriceCurrency string `json:"line_04_price_currency"`
	Line04Price float64 `json:"line_04_price"`
	Line04AmountCurrency string `json:"line_04_amount_currency"`
	Line04Amount float64 `json:"line_04_amount"`
	Line05Qty null.Int `json:"line_05_qty"` // Make `int8`
	Line05Desc string `json:"line_05_desc"`
	Line05PriceCurrency string `json:"line_05_price_currency"`
	Line05Price float64 `json:"line_05_price"`
	Line05AmountCurrency string `json:"line_05_amount_currency"`
	Line05Amount float64 `json:"line_05_amount"`
	Line06Qty null.Int `json:"line_06_qty"` // Make `int8`
	Line06Desc string `json:"line_06_desc"`
	Line06PriceCurrency string `json:"line_06_price_currency"`
	Line06Price float64 `json:"line_06_price"`
	Line06AmountCurrency string `json:"line_06_amount_currency"`
	Line06Amount float64 `json:"line_06_amount"`
	Line07Qty null.Int `json:"line_07_qty"` // Make `int8`
	Line07Desc string `json:"line_07_desc"`
	Line07PriceCurrency string `json:"line_07_price_currency"`
	Line07Price float64 `json:"line_07_price"`
	Line07AmountCurrency string `json:"line_07_amount_currency"`
	Line07Amount float64 `json:"line_07_amount"`
	Line08Qty null.Int `json:"line_08_qty"` // Make `int8`
	Line08Desc string `json:"line_08_desc"`
	Line08PriceCurrency string `json:"line_08_price_currency"`
	Line08Price float64 `json:"line_08_price"`
	Line08AmountCurrency string `json:"line_08_amount_currency"`
	Line08Amount float64 `json:"line_08_amount"`
	Line09Qty null.Int `json:"line_09_qty"` // Make `int8`
	Line09Desc string `json:"line_09_desc"`
	Line09PriceCurrency string `json:"line_09_price_currency"`
	Line09Price float64 `json:"line_09_price"`
	Line09AmountCurrency string `json:"line_09_amount_currency"`
	Line09Amount float64 `json:"line_09_amount"`
	Line10Qty null.Int `json:"line_10_qty"` // Make `int8`
	Line10Desc string `json:"line_10_desc"`
	Line10PriceCurrency string `json:"line_10_price_currency"`
	Line10Price float64 `json:"line_10_price"`
	Line10AmountCurrency string `json:"line_10_amount_currency"`
	Line10Amount float64 `json:"line_10_amount"`
	Line11Qty null.Int `json:"line_11_qty"` // Make `int8`
	Line11Desc string `json:"line_11_desc"`
	Line11PriceCurrency string `json:"line_11_price_currency"`
	Line11Price float64 `json:"line_11_price"`
	Line11AmountCurrency string `json:"line_11_amount_currency"`
	Line11Amount float64 `json:"line_11_amount"`
	Line12Qty null.Int `json:"line_12_qty"` // Make `int8`
	Line12Desc string `json:"line_12_desc"`
	Line12PriceCurrency string `json:"line_12_price_currency"`
	Line12Price float64 `json:"line_12_price"`
	Line12AmountCurrency string `json:"line_12_amount_currency"`
	Line12Amount float64 `json:"line_12_amount"`
	Line13Qty null.Int `json:"line_13_qty"` // Make `int8`
	Line13Desc string `json:"line_13_desc"`
	Line13PriceCurrency string `json:"line_13_price_currency"`
	Line13Price float64 `json:"line_13_price"`
	Line13AmountCurrency string `json:"line_13_amount_currency"`
	Line13Amount float64 `json:"line_13_amount"`
	Line14Qty null.Int `json:"line_14_qty"` // Make `int8`
	Line14Desc string `json:"line_14_desc"`
	Line14PriceCurrency string `json:"line_14_price_currency"`
	Line14Price float64 `json:"line_14_price"`
	Line14AmountCurrency string `json:"line_14_amount_currency"`
	Line14Amount float64 `json:"line_14_amount"`
	Line15Qty null.Int `json:"line_15_qty"` // Make `int8`
	Line15Desc string `json:"line_15_desc"`
	Line15PriceCurrency string `json:"line_15_price_currency"`
	Line15Price float64 `json:"line_15_price"`
	Line15AmountCurrency string `json:"line_15_amount_currency"`
	Line15Amount float64 `json:"line_15_amount"`
	InvoiceQuoteDays int8 `json:"invoice_quote_days"`
	InvoiceAssociateTax null.String `json:"invoice_associate_tax"`
	InvoiceQuoteDate time.Time `json:"invoice_quote_date"`
	InvoiceCustomersApproval string `json:"invoice_customers_approval"`
	Line01Notes null.String `json:"line_01_notes"`
	Line02Notes null.String `json:"line_02_notes"`
	TotalLabourCurrency string `json:"total_labour_currency"`
	TotalLabour float64 `json:"total_labour"`
	TotalMaterialsCurrency string `json:"total_materials_currency"`
	TotalMaterials float64 `json:"total_materials"`
	OtherCostsCurrency string `json:"other_costs_currency"`
	OtherCosts float64 `json:"other_costs"`
	AmountDueCurrency string `json:"amount_due_currency"`
	TaxCurrency string `json:"tax_currency"`
	Tax float64 `json:"tax"`
	TotalCurrency string `json:"total_currency"`
	Total float64 `json:"total"`
	DepositCurrency string `json:"deposit_currency"`
	PaymentAmountCurrency string `json:"payment_amount_currency"`
	PaymentAmount float64 `json:"payment_amount"`
	PaymentDate time.Time `json:"payment_date"`
	IsCash bool `json:"is_cash"`
	IsCheque bool `json:"is_cheque"`
	IsDebit bool `json:"is_debit"`
	IsCredit bool `json:"is_credit"`
	IsOther bool `json:"is_other"`
	ClientSignature string `json:"client_signature"`
	AssociateSignDate time.Time `json:"associate_sign_date"`
	AssociateSignature string `json:"associate_signature"`
	WorkOrderId   uint64 `json:"work_order_id"`
	CreatedAt time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	CreatedById   uint64 `json:"created_by_id"`
	LastModifiedById   uint64 `json:"last_modified_by_id"`
	CreatedFrom string `json:"created_from"`
	CreatedFromIsPublic bool `json:"created_from_is_public"`
	LastModifiedFrom string `json:"last_modified_from"`
	LastModifiedFromIsPublic bool `json:"last_modified_from_is_public"`
	ClientAddress string `json:"client_address"`
	RevisionVersion int8 `json:"revision_version"`
	Deposit float64 `json:"deposit"`
	AmountDue float64 `json:"amount_due"`
	SubTotal float64 `json:"sub_total"`
	SubTotalCurrency string `json:"sub_total_currency"`
}


func ListAllWorkOrderInvoices(db *sql.DB) ([]*OldWorkOrderInvoice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT
        order_id, is_archived, invoice_id, invoice_date, associate_name,
		associate_telephone, client_name, client_telephone, client_email,
		line_01_qty, line_01_desc, line_01_price_currency, line_01_price, line_01_amount_currency, line_01_amount,
		line_02_qty, line_02_desc, line_02_price_currency, line_02_price, line_02_amount_currency, line_02_amount,
		line_03_qty, line_03_desc, line_03_price_currency, line_03_price, line_03_amount_currency, line_03_amount,
		line_04_qty, line_04_desc, line_04_price_currency, line_04_price, line_04_amount_currency, line_04_amount,
		line_05_qty, line_05_desc, line_05_price_currency, line_05_price, line_05_amount_currency, line_05_amount,
		line_06_qty, line_06_desc, line_06_price_currency, line_06_price, line_06_amount_currency, line_06_amount,
		line_07_qty, line_07_desc, line_07_price_currency, line_07_price, line_07_amount_currency, line_07_amount,
		line_08_qty, line_08_desc, line_08_price_currency, line_08_price, line_08_amount_currency, line_08_amount,
		line_09_qty, line_09_desc, line_09_price_currency, line_09_price, line_09_amount_currency, line_09_amount,
		line_10_qty, line_10_desc, line_10_price_currency, line_10_price, line_10_amount_currency, line_10_amount,
		line_11_qty, line_11_desc, line_11_price_currency, line_11_price, line_11_amount_currency, line_11_amount,
		line_12_qty, line_12_desc, line_12_price_currency, line_12_price, line_12_amount_currency, line_12_amount,
		line_13_qty, line_13_desc, line_13_price_currency, line_13_price, line_13_amount_currency, line_13_amount,
		line_14_qty, line_14_desc, line_14_price_currency, line_14_price, line_14_amount_currency, line_14_amount,
		line_15_qty, line_15_desc, line_15_price_currency, line_15_price, line_15_amount_currency, line_15_amount,
		invoice_quote_days, invoice_associate_tax, invoice_quote_date, invoice_customers_approval, line_01_notes,
		line_02_notes, total_labour_currency, total_labour, total_materials_currency, total_materials,
		other_costs_currency, other_costs, amount_due_currency, tax_currency, tax, total_currency, total,
		deposit_currency, payment_amount_currency, payment_amount, payment_date, is_cash, is_cheque, is_debit,
		is_credit, is_other, client_signature, associate_sign_date, associate_signature, work_order_id, created_at,
		last_modified_at, created_by_id, last_modified_by_id, created_from, created_from_is_public, last_modified_from,
		last_modified_from_is_public, client_address, revision_version, deposit, amount_due, sub_total, sub_total_currency
	FROM
        workery_work_order_invoices
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []*OldWorkOrderInvoice
	defer rows.Close()
	for rows.Next() {
		m := new(OldWorkOrderInvoice)
		err = rows.Scan(
		    &m.OrderId, &m.IsArchived, &m.InvoiceId, &m.InvoiceDate, &m.AssociateName,
			&m.AssociateTelephone, &m.ClientName, &m.ClientTelephone, &m.ClientEmail,
			&m.Line01Qty, &m.Line01Desc, &m.Line01PriceCurrency, &m.Line01Price, &m.Line01AmountCurrency, &m.Line01Amount,
			&m.Line02Qty, &m.Line02Desc, &m.Line02PriceCurrency, &m.Line02Price, &m.Line02AmountCurrency, &m.Line02Amount,
			&m.Line03Qty, &m.Line03Desc, &m.Line03PriceCurrency, &m.Line03Price, &m.Line03AmountCurrency, &m.Line03Amount,
			&m.Line04Qty, &m.Line04Desc, &m.Line04PriceCurrency, &m.Line04Price, &m.Line04AmountCurrency, &m.Line04Amount,
			&m.Line05Qty, &m.Line05Desc, &m.Line05PriceCurrency, &m.Line05Price, &m.Line05AmountCurrency, &m.Line05Amount,
			&m.Line06Qty, &m.Line06Desc, &m.Line06PriceCurrency, &m.Line06Price, &m.Line06AmountCurrency, &m.Line06Amount,
			&m.Line07Qty, &m.Line07Desc, &m.Line07PriceCurrency, &m.Line07Price, &m.Line07AmountCurrency, &m.Line07Amount,
			&m.Line08Qty, &m.Line08Desc, &m.Line08PriceCurrency, &m.Line08Price, &m.Line08AmountCurrency, &m.Line08Amount,
			&m.Line09Qty, &m.Line09Desc, &m.Line09PriceCurrency, &m.Line09Price, &m.Line09AmountCurrency, &m.Line09Amount,
			&m.Line10Qty, &m.Line10Desc, &m.Line10PriceCurrency, &m.Line10Price, &m.Line10AmountCurrency, &m.Line10Amount,
			&m.Line11Qty, &m.Line11Desc, &m.Line11PriceCurrency, &m.Line11Price, &m.Line11AmountCurrency, &m.Line11Amount,
			&m.Line12Qty, &m.Line12Desc, &m.Line12PriceCurrency, &m.Line12Price, &m.Line12AmountCurrency, &m.Line12Amount,
			&m.Line13Qty, &m.Line13Desc, &m.Line13PriceCurrency, &m.Line13Price, &m.Line13AmountCurrency, &m.Line13Amount,
			&m.Line14Qty, &m.Line14Desc, &m.Line14PriceCurrency, &m.Line14Price, &m.Line14AmountCurrency, &m.Line14Amount,
			&m.Line15Qty, &m.Line15Desc, &m.Line15PriceCurrency, &m.Line15Price, &m.Line15AmountCurrency, &m.Line15Amount,
			&m.InvoiceQuoteDays, &m.InvoiceAssociateTax, &m.InvoiceQuoteDate, &m.InvoiceCustomersApproval, &m.Line01Notes,
			&m.Line02Notes, &m.TotalLabourCurrency, &m.TotalLabour, &m.TotalMaterialsCurrency, &m.TotalMaterials,
			&m.OtherCostsCurrency, &m.OtherCosts, &m.AmountDueCurrency, &m.TaxCurrency, &m.Tax, &m.TotalCurrency, &m.Total,
			&m.DepositCurrency, &m.PaymentAmountCurrency, &m.PaymentAmount, &m.PaymentDate, &m.IsCash, &m.IsCheque,
			&m.IsDebit, &m.IsCredit, &m.IsOther, &m.ClientSignature, &m.AssociateSignDate, &m.AssociateSignature,
			&m.WorkOrderId, &m.CreatedAt, &m.LastModifiedAt, &m.CreatedById, &m.LastModifiedById, &m.CreatedFrom,
			&m.CreatedFromIsPublic, &m.LastModifiedFrom, &m.LastModifiedFromIsPublic, &m.ClientAddress, &m.RevisionVersion,
			&m.Deposit, &m.AmountDue, &m.SubTotal, &m.SubTotalCurrency,
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

func insertWorkOrderInvoiceETL(
	ctx context.Context,
	tid uint64,
	wotp *repositories.WorkOrderInvoiceRepo,
	wor *repositories.WorkOrderRepo,
	oss *OldWorkOrderInvoice,
) {
	//
	// OrderId
	//

	orderId, err := wor.GetIdByOldId(ctx, tid, oss.OrderId)
	if err != nil {
		log.Panic("ar.GetIdByOldId | err", err)
	}

	// //
	// // InvoiceId
	// //
	//
	// invoiceId, err := vtr.GetIdByOldId(ctx, tid, oss.InvoiceId)
	// if err != nil {
	// 	log.Panic("ar.GetIdByOldId | err", err)
	// }

	//
	// Insert into database.
	//

	m := &models.WorkOrderInvoice{
		Uuid:                uuid.NewString(),        // 1
		TenantId:            tid,                     // 2
		OldId:               orderId,                 // 3
		InvoiceId:           oss.InvoiceId,           // 4
		OrderId:             orderId,                 // 5
        InvoiceDate:         oss.InvoiceDate,         // 6
		AssociateName:       oss.AssociateName,       // 7
		AssociateTelephone:  oss.AssociateTelephone,  // 8
		ClientName:          oss.ClientName,          // 9
        ClientTelephone:     oss.ClientTelephone,     // 10
		ClientEmail:         oss.ClientEmail,         // 11

		Line01Qty:           oss.Line01Qty,           // 12
		Line01Desc:          oss.Line01Desc,          // 13
		Line01Price:         oss.Line01Price,         // 14
		Line01Amount:        oss.Line01Amount,        // 15

		Line02Qty:           oss.Line02Qty,           // 16
		Line02Desc:          oss.Line02Desc,          // 17
		Line02Price:         oss.Line02Price,         // 18
        Line02Amount:        oss.Line02Amount,        // 19

		Line03Qty:           oss.Line03Qty,           // 20
		Line03Desc:          oss.Line03Desc,          // 21
		Line03Price:         oss.Line03Price,         // 22
        Line03Amount:        oss.Line03Amount,        // 23

		Line04Qty:           oss.Line04Qty,           // 24
		Line04Desc:          oss.Line04Desc,          // 25
		Line04Price:         oss.Line04Price,         // 26
        Line04Amount:        oss.Line04Amount,        // 27

		Line05Qty:           oss.Line05Qty,           // 28
		Line05Desc:          oss.Line05Desc,          // 29
		Line05Price:         oss.Line05Price,         // 30
        Line05Amount:        oss.Line05Amount,        // 31

		Line06Qty:           oss.Line06Qty,           // 32
		Line06Desc:          oss.Line06Desc,          // 33
		Line06Price:         oss.Line06Price,         // 34
        Line06Amount:        oss.Line06Amount,        // 35

		Line07Qty:           oss.Line07Qty,           // 36
		Line07Desc:          oss.Line07Desc,          // 37
		Line07Price:         oss.Line07Price,         // 38
        Line07Amount:        oss.Line07Amount,        // 39

		Line08Qty:           oss.Line08Qty,           // 40
		Line08Desc:          oss.Line08Desc,          // 41
		Line08Price:         oss.Line08Price,         // 42
        Line08Amount:        oss.Line08Amount,        // 43

		Line09Qty:           oss.Line09Qty,           // 44
		Line09Desc:          oss.Line09Desc,          // 45
		Line09Price:         oss.Line09Price,         // 46
        Line09Amount:        oss.Line09Amount,        // 47

		Line10Qty:           oss.Line10Qty,           // 48
		Line10Desc:          oss.Line10Desc,          // 49
		Line10Price:         oss.Line10Price,         // 50
        Line10Amount:        oss.Line10Amount,        // 51

		// Line15Qty null.Int `json:"line_15_qty"` // Make `int8`
		// Line15Desc string `json:"line_15_desc"`
		// Line15Price float64 `json:"line_15_price"`
		// Line15Amount float64 `json:"line_15_amount"`
		// InvoiceQuoteDays int8 `json:"invoice_quote_days"`
		// InvoiceAssociateTax null.String `json:"invoice_associate_tax"`
		// InvoiceQuoteDate time.Time `json:"invoice_quote_date"`
		// InvoiceCustomersApproval string `json:"invoice_customers_approval"`
		// Line01Notes null.String `json:"line_01_notes"`
		// Line02Notes null.String `json:"line_02_notes"`
		// TotalLabour float64 `json:"total_labour"`
		// TotalMaterials float64 `json:"total_materials"`
		// OtherCosts float64 `json:"other_costs"`
		// Tax float64 `json:"tax"`
		// Total float64 `json:"total"`
		// PaymentAmount float64 `json:"payment_amount"`
		// PaymentDate time.Time `json:"payment_date"`
		// IsCash bool `json:"is_cash"`
		// IsCheque bool `json:"is_cheque"`
		// IsDebit bool `json:"is_debit"`
		// IsCredit bool `json:"is_credit"`
		// IsOther bool `json:"is_other"`
		// ClientSignature string `json:"client_signature"`
		// AssociateSignDate time.Time `json:"associate_sign_date"`
		// AssociateSignature string `json:"associate_signature"`
		// WorkOrderId   uint64 `json:"work_order_id"`
		// CreatedAt time.Time `json:"created_at"`
		// LastModifiedAt time.Time `json:"last_modified_at"`
		// CreatedById   uint64 `json:"created_by_id"`
		// LastModifiedById   uint64 `json:"last_modified_by_id"`
		// CreatedFrom string `json:"created_from"`
		// CreatedFromIsPublic bool `json:"created_from_is_public"`
		// LastModifiedFrom string `json:"last_modified_from"`
		// LastModifiedFromIsPublic bool `json:"last_modified_from_is_public"`
		// ClientAddress string `json:"client_address"`
		// RevisionVersion int8 `json:"revision_version"`
		// Deposit float64 `json:"deposit"`
		// AmountDue float64 `json:"amount_due"`
		// SubTotal float64 `json:"sub_total"`
	    // State       int8 `json:"state"`	 // IsArchived bool `json:"is_archived"`
	}

	fmt.Println("OrderId:", orderId)

	err = wotp.Insert(ctx, m)
	if err != nil {
		log.Panic("wotp.Insert | err", err)
	}
	fmt.Println("Imported ID#", orderId)
}
