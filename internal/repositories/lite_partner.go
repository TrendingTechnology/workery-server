package repositories

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type LitePartnerRepo struct {
	db *sql.DB
}

func NewLitePartnerRepo(db *sql.DB) *LitePartnerRepo {
	return &LitePartnerRepo{
		db: db,
	}
}

func (s *LitePartnerRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.LitePartnerFilter) (*sql.Rows, error) {
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
	// log.Println("LitePartnerRepo | query:", query, "\n")
	// log.Println("LitePartnerRepo | filterValues:", filterValues, "\n")

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *LitePartnerRepo) ListByFilter(ctx context.Context, filter *models.LitePartnerFilter) ([]*models.LitePartner, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        id,
		tenant_id,
		state,
		given_name,
		last_name,
		telephone,
		email,
		join_date,
		type_of
    FROM
        partners
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.LitePartner
	defer rows.Close()
	for rows.Next() {
		m := new(models.LitePartner)
		err := rows.Scan(
			&m.Id,
			&m.TenantId,
			&m.State,
			&m.GivenName,
			&m.LastName,
			&m.Telephone,
			&m.Email,
			&m.JoinDate,
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

func (s *LitePartnerRepo) CountByFilter(ctx context.Context, f *models.LitePartnerFilter) (uint64, error) {
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
	    partners
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

	// For debugging purposes only.
	// log.Println("query:", query)
	// log.Println("filterValues:", filterValues)

	// Return our values.
	return count, err
}
