package repositories

import (
	// "log"
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type LiteTenantRepo struct {
	db *sql.DB
}

func NewLiteTenantRepo(db *sql.DB) *LiteTenantRepo {
	return &LiteTenantRepo{
		db: db,
	}
}

func (s *LiteTenantRepo) ListAllIds(ctx context.Context) ([]uint64, error) {
	query := `SELECT id FROM tenants ORDER BY (id) ASC`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var arr []uint64
	defer rows.Close()
	for rows.Next() {
		var sessionId uint64
		err = rows.Scan(
			&sessionId,
		)
		if err != nil {
			panic(err)
		}
		arr = append(arr, sessionId)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arr, err
}

func (s *LiteTenantRepo) queryRowsWithFilter(ctx context.Context, query string, filter *models.LiteTenantFilter) (*sql.Rows, error) {
	// Array will hold all the unique values we want to add into the query.
	var filterValues []interface{}

	// The SQL query statement we will be calling in the database, start
	// by setting the `tenant_id` placeholder and then append our value to
	// the array.
	filterValues = append(filterValues, filter.LastSeenId)
	query += ` WHERE id > $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our OPTIONAL filters
	//

	if !filter.State.IsZero() {
		filterValues = append(filterValues, filter.State)
		query += `AND state = $` + strconv.Itoa(len(filterValues))
	}

	//
	// The following code will add our REQUIRED filters
	//

	query += ` ORDER BY id`
	filterValues = append(filterValues, filter.Limit)
	query += ` DESC LIMIT $` + strconv.Itoa(len(filterValues))

	//
	// Execute our custom built SQL query to the database.
	// (Notice our usage of the `variadic function`?)
	//

	// log.Println("query:", query)
	// log.Println("filterValues:", filterValues)

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *LiteTenantRepo) ListByFilter(ctx context.Context, filter *models.LiteTenantFilter) ([]*models.LiteTenant, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        id, schema_name, name, state
    FROM
        tenants`

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.LiteTenant
	defer rows.Close()
	for rows.Next() {
		m := new(models.LiteTenant)
		err := rows.Scan(
			&m.Id,
			&m.SchemaName,
			&m.Name,
			&m.State,
		)
		if err != nil {
			return nil, err
		}
		arr = append(arr, m)
	}
	err = rows.Err()
	return arr, err
}

func (s *LiteTenantRepo) CountByFilter(ctx context.Context, filter *models.LiteTenantFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// The result we are looking for.
	var count uint64

	// Array will hold all the unique values we want to add into the query.
	var filterValues []interface{}

	// The SQL query statement we will be calling in the database, start
	// by setting the `tenant_id` placeholder and then append our value to
	// the array.
	filterValues = append(filterValues)
	query := `
    SELECT COUNT(id) FROM
        tenants`

	//
	// The following code will add our OPTIONAL filters
	//

	if !filter.State.IsZero() {
		filterValues = append(filterValues, filter.State)
		query += `AND state = $` + strconv.Itoa(len(filterValues))
	}

	//
	// Execute our custom built SQL query to the database.
	//

	err := s.db.QueryRowContext(ctx, query, filterValues...).Scan(&count)
	return count, err
}
