package repositories

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type LiteTaskItemRepo struct {
	db *sql.DB
}

func NewLiteTaskItemRepo(db *sql.DB) *LiteTaskItemRepo {
	return &LiteTaskItemRepo{
		db: db,
	}
}

func (s *LiteTaskItemRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.LiteTaskItemFilter) (*sql.Rows, error) {
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

	if !f.IsClosed.IsZero() {
		filterValues = append(filterValues, f.IsClosed.ValueOrZero())
		query += ` AND is_closed = $` + strconv.Itoa(len(filterValues))
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

	// log.Println("QUERY:", query)
	// log.Println("VALUES:", filterValues)
	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *LiteTaskItemRepo) ListByFilter(ctx context.Context, filter *models.LiteTaskItemFilter) ([]*models.LiteTaskItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        id,
		tenant_id,
		state,
		due_date,
		type_of,
		customer_id,
		customer_name,
		customer_lexical_name,
		associate_id,
		associate_name,
		associate_lexical_name,
		order_type_of
    FROM
        task_items
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.LiteTaskItem
	defer rows.Close()
	for rows.Next() {
		m := new(models.LiteTaskItem)
		err := rows.Scan(
			&m.Id,
			&m.TenantId,
			&m.State,
			&m.DueDate,
			&m.TypeOf,
			&m.CustomerId,
			&m.CustomerName,
			&m.CustomerLexicalName,
			&m.AssociateId,
			&m.AssociateName,
			&m.AssociateLexicalName,
			&m.OrderTypeOf,
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

func (s *LiteTaskItemRepo) CountByFilter(ctx context.Context, f *models.LiteTaskItemFilter) (uint64, error) {
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
	    task_items
	WHERE
		tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our filters
	//

	if !f.IsClosed.IsZero() {
		filterValues = append(filterValues, f.IsClosed.ValueOrZero())
		query += ` AND is_closed = $` + strconv.Itoa(len(filterValues))
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

	// log.Println("QUERY:", query)
	// log.Println("VALUES:", filterValues)
	return count, err
}
