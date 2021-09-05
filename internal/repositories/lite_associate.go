package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"log"

	"github.com/over55/workery-server/internal/models"
)

type LiteAssociateRepo struct {
	db *sql.DB
}

func NewLiteAssociateRepo(db *sql.DB) *LiteAssociateRepo {
	return &LiteAssociateRepo{
		db: db,
	}
}

func (s *LiteAssociateRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.LiteAssociateFilter) (*sql.Rows, error) {
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

    // log.Println("QUERY:", query)
	// log.Println("VALUES:", filterValues)
	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *LiteAssociateRepo) ListByFilter(ctx context.Context, filter *models.LiteAssociateFilter) ([]*models.LiteAssociate, error) {
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
		telephone_type_of,
		telephone_extension,
		email,
		join_date
    FROM
        associates
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.LiteAssociate
	defer rows.Close()
	for rows.Next() {
		m := new(models.LiteAssociate)
		err := rows.Scan(
			&m.Id,
			&m.TenantId,
			&m.State,
			&m.GivenName,
			&m.LastName,
			&m.Telephone,
			&m.TelephoneTypeOf,
			&m.TelephoneExtension,
			&m.Email,
			&m.JoinDate,
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

func (s *LiteAssociateRepo) CountByFilter(ctx context.Context, f *models.LiteAssociateFilter) (uint64, error) {
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
	    associates
	WHERE
		tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our filters
	//

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
	return count, err
}
