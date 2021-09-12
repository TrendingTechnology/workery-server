package repositories

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type LiteFinancialRepo struct {
	db *sql.DB
}

func NewLiteFinancialRepo(db *sql.DB) *LiteFinancialRepo {
	return &LiteFinancialRepo{
		db: db,
	}
}

func (s *LiteFinancialRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.LiteFinancialFilter) (*sql.Rows, error) {
	// Array will hold all the unique values we want to add into the query.
	var filterValues []interface{}

	// The SQL query statement we will be calling in the database, start
	// by setting the `tenant_id` placeholder and then append our value to
	// the array.
	filterValues = append(filterValues, f.TenantId)
	query += ` WHERE tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our OPTIONAL filters
	//

	if !f.LastModifiedById.IsZero() {
		filterValues = append(filterValues, f.LastModifiedById)
		query += ` AND last_modified_by_id = $` + strconv.Itoa(len(filterValues))
	}

	if !f.AssociateName.IsZero() {
		filterValues = append(filterValues, f.AssociateLexicalName)
		query += ` AND associate_name = $` + strconv.Itoa(len(filterValues))
	}

	if !f.AssociateLexicalName.IsZero() {
		filterValues = append(filterValues, f.AssociateLexicalName)
		query += ` AND associate_lexical_name = $` + strconv.Itoa(len(filterValues))
	}

	if !f.CustomerName.IsZero() {
		filterValues = append(filterValues, f.CustomerLexicalName)
		query += ` AND customer_name = $` + strconv.Itoa(len(filterValues))
	}

	if !f.CustomerLexicalName.IsZero() {
		filterValues = append(filterValues, f.CustomerLexicalName)
		query += ` AND customer_lexical_name = $` + strconv.Itoa(len(filterValues))
	}

	if !f.Search.IsZero() {
		log.Fatal("TODO: PLEASE IMPLEMENT")
		// filterValues = append(filterValues, f.Search)
		// query += `AND state = $` + strconv.Itoa(len(filterValues))
	}

	if len(f.States) > 0 {
		query += ` AND (`
		for i, v := range f.States {
			s := strconv.Itoa(int(v))
			filterValues = append(filterValues, s)
			if i != 0 {
				query += ` OR`
			}
			query += ` state = $` + strconv.Itoa(len(filterValues))
		}
		query += ` )`
	}

	//
	// The following code will add our pagination.
	//

	query += ` ORDER BY ` + f.SortField + ` ` + f.SortOrder
	filterValues = append(filterValues, f.Limit)
	query += ` LIMIT $` + strconv.Itoa(len(filterValues))
	filterValues = append(filterValues, f.Offset)
	query += ` OFFSET $` + strconv.Itoa(len(filterValues))

	//
	// Execute our custom built SQL query to the database.
	//

	// For debugging purposes only.
	// log.Println("LiteCustomerRepo | query:", query, "\n")
	// log.Println("LiteCustomerRepo | SortField:", f.SortField, "SortField" + f.SortOrder + "\n")
	// log.Println("LiteCustomerRepo | filterValues:", filterValues, "\n")

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *LiteFinancialRepo) ListByFilter(ctx context.Context, filter *models.LiteFinancialFilter) ([]*models.LiteFinancial, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        id,
		tenant_id,
		state,
		customer_id,
		customer_name,
		associate_id,
		associate_name,
		invoice_service_fee_payment_date,
		type_of
    FROM
        work_orders
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.LiteFinancial
	defer rows.Close()
	for rows.Next() {
		m := new(models.LiteFinancial)
		err := rows.Scan(
			&m.Id,
			&m.TenantId,
			&m.State,
			&m.CustomerId,
			&m.CustomerName,
			&m.AssociateId,
			&m.AssociateName,
			&m.InvoiceServiceFeePaymentDate,
			&m.TypeOf,
		)
		if err != nil {
			return nil, err
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return arr, err
}

func (s *LiteFinancialRepo) CountByFilter(ctx context.Context, f *models.LiteFinancialFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// The result we are looking for.
	var count uint64

	// Array will hold all the unique values we want to add into the query.
	var filterValues []interface{}

	// The SQL query statement we will be calling in the database, start
	// by setting the `tenant_id` placeholder and then append our value to
	// the array.
	filterValues = append(filterValues, f.TenantId)
	query := `
	SELECT COUNT(id) FROM
	    work_orders
	WHERE
		tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our filters
	//

	if !f.LastModifiedById.IsZero() {
		filterValues = append(filterValues, f.LastModifiedById)
		query += ` AND last_modified_by_id = $` + strconv.Itoa(len(filterValues))
	}

	if !f.Search.IsZero() {
		log.Fatal("TODO: PLEASE IMPLEMENT")
		// filterValues = append(filterValues, f.Search)
		// query += `AND state = $` + strconv.Itoa(len(filterValues))
	}

	if len(f.States) > 0 {
		query += ` AND (`
		for i, v := range f.States {
			s := strconv.Itoa(int(v))
			filterValues = append(filterValues, s)
			if i != 0 {
				query += ` OR`
			}
			query += ` state = $` + strconv.Itoa(len(filterValues))
		}
		query += ` )`
	}

	//
	// Execute our custom built SQL query to the database.
	//

	err := s.db.QueryRowContext(ctx, query, filterValues...).Scan(&count)

	// Return our values.
	return count, err
}
