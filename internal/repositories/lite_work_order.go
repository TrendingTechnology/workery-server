package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type LiteWorkOrderRepo struct {
	db *sql.DB
}

func NewLiteWorkOrderRepo(db *sql.DB) *LiteWorkOrderRepo {
	return &LiteWorkOrderRepo{
		db: db,
	}
}

func (s *LiteWorkOrderRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.LiteWorkOrderFilter) (*sql.Rows, error) {
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
		query += `AND last_modified_by_id = $` + strconv.Itoa(len(filterValues))
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

	if f.LastSeenId > 0 {
		filterValues = append(filterValues, f.LastSeenId)
		query += ` AND id < $` + strconv.Itoa(len(filterValues))
	}
	query += ` ORDER BY last_modified_time `
	filterValues = append(filterValues, f.Limit)
	query += ` DESC LIMIT $` + strconv.Itoa(len(filterValues))

	//
	// Execute our custom built SQL query to the database.
	//

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *LiteWorkOrderRepo) ListByFilter(ctx context.Context, filter *models.LiteWorkOrderFilter) ([]*models.LiteWorkOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        id,
		tenant_id,
		state,
		customer_id,
		associate_id
    FROM
        work_orders
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.LiteWorkOrder
	defer rows.Close()
	for rows.Next() {
		m := new(models.LiteWorkOrder)
		err := rows.Scan(
			&m.Id,
			&m.TenantId,
			&m.State,
			&m.CustomerId,
			&m.AssociateId,
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

func (s *LiteWorkOrderRepo) CountByFilter(ctx context.Context, f *models.LiteWorkOrderFilter) (uint64, error) {
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
		query += `AND last_modified_by_id = $` + strconv.Itoa(len(filterValues))
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
