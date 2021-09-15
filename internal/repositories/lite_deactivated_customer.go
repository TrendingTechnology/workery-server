package repositories

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type LiteDeactivatedCustomerRepo struct {
	db *sql.DB
}

func NewLiteDeactivatedCustomerRepo(db *sql.DB) *LiteDeactivatedCustomerRepo {
	return &LiteDeactivatedCustomerRepo{
		db: db,
	}
}

func (s *LiteDeactivatedCustomerRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.LiteDeactivatedCustomerFilter) (*sql.Rows, error) {
	// Array will hold all the unique values we want to add into the query.
	var filterValues []interface{}

	// The SQL query statement we will be calling in the database, start
	// by setting the `tenant_id` placeholder and then append our value to
	// the array.
	filterValues = append(filterValues, f.TenantId)
	query += ` WHERE tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our filters
	//

	if !f.Search.IsZero() {
		log.Fatal("TODO: PLEASE IMPLEMENT")
		// filterValues = append(filterValues, f.Search)
		// query += `AND state = $` + strconv.Itoa(len(filterValues))
	}

	// if len(f.States) > 0 {
	// 	query += ` AND (`
	// 	for i, v := range f.States {
	// 		s := strconv.Itoa(int(v))
	// 		filterValues = append(filterValues, s)
	// 		if i != 0 {
	// 			query += ` OR`
	// 		}
	// 		query += ` state = $` + strconv.Itoa(len(filterValues))
	// 	}
	// 	query += ` )`
	// }

	// Deactivated
	query += ` AND state = 0` // query += ` state = $` + strconv.Itoa(len(filterValues))

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
	// log.Println("LiteDeactivatedCustomerRepo | query:", query, "\n")
	// log.Println("LiteDeactivatedCustomerRepo | filterValues:", filterValues, "\n")

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *LiteDeactivatedCustomerRepo) ListByFilter(ctx context.Context, filter *models.LiteDeactivatedCustomerFilter) ([]*models.LiteDeactivatedCustomer, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        id,
		tenant_id,
		name,
		lexical_name,
		deactivation_reason,
		deactivation_reason_other,
		state
    FROM
        customers
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.LiteDeactivatedCustomer
	defer rows.Close()
	for rows.Next() {
		m := new(models.LiteDeactivatedCustomer)
		err := rows.Scan(
			&m.Id,
			&m.TenantId,
			&m.Name,
			&m.LexicalName,
			&m.DeactivationReason,
			&m.DeactivationReasonOther,
			&m.State,
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

func (s *LiteDeactivatedCustomerRepo) CountByFilter(ctx context.Context, f *models.LiteDeactivatedCustomerFilter) (uint64, error) {
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
	    customers
	WHERE
		tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our filters
	//

	if !f.Search.IsZero() {
		log.Fatal("TODO: PLEASE IMPLEMENT")
		// filterValues = append(filterValues, f.Search)
		// query += `AND state = $` + strconv.Itoa(len(filterValues))
	}

	// if len(f.States) > 0 {
	// 	query += ` AND (`
	// 	for i, v := range f.States {
	// 		s := strconv.Itoa(int(v))
	// 		filterValues = append(filterValues, s)
	// 		if i != 0 {
	// 			query += ` OR`
	// 		}
	// 		query += ` state = $` + strconv.Itoa(len(filterValues))
	// 	}
	// 	query += ` )`
	// }

	// Deactivated
	query += ` AND state = 0` // query += ` state = $` + strconv.Itoa(len(filterValues))

	//
	// Execute our custom built SQL query to the database.
	//

	err := s.db.QueryRowContext(ctx, query, filterValues...).Scan(&count)

	// For debugging purposes only.
	// log.Println("query:", query)
	// log.Println("filterValues:", filterValues)

	// Return our values.
	return count, err
}
