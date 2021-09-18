package repositories

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type CustomerTagRepo struct {
	db *sql.DB
}

func NewCustomerTagRepo(db *sql.DB) *CustomerTagRepo {
	return &CustomerTagRepo{
		db: db,
	}
}

func (r *CustomerTagRepo) Insert(ctx context.Context, m *models.CustomerTag) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO customer_tags (
        uuid, tenant_id, customer_id, tag_id, old_id
    ) VALUES (
        $1, $2, $3, $4, $5
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.CustomerId, m.TagId, m.OldId,
	)
	return err
}

func (r *CustomerTagRepo) UpdateById(ctx context.Context, m *models.CustomerTag) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        customer_tags
    SET
        tenant_id = $1, customer_id = $2, tag_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.CustomerId, m.TagId, m.Id,
	)
	return err
}

func (r *CustomerTagRepo) GetById(ctx context.Context, id uint64) (*models.CustomerTag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.CustomerTag)

	query := `
    SELECT
        id, uuid, tenant_id, customer_id, tag_id
	FROM
        customer_tags
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.CustomerId, &m.TagId,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *CustomerTagRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        customer_tags
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *CustomerTagRepo) InsertOrUpdateById(ctx context.Context, m *models.CustomerTag) error {
	if m.Id == 0 {
		return r.Insert(ctx, m)
	}

	doesExist, err := r.CheckIfExistsById(ctx, m.Id)
	if err != nil {
		return err
	}

	if doesExist == false {
		return r.Insert(ctx, m)
	}
	return r.UpdateById(ctx, m)
}

func (s *CustomerTagRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.CustomerTagFilter) (*sql.Rows, error) {
	// Array will hold all the unique values we want to add into the query.
	var filterValues []interface{}

	// The SQL query statement we will be calling in the database, start
	// by setting the `tenant_id` placeholder and then append our value to
	// the array.
	filterValues = append(filterValues, f.TenantId)
	query += ` WHERE customer_tags.tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our filters
	//

	if !f.Search.IsZero() {
		log.Fatal("TODO: PLEASE IMPLEMENT")
		// filterValues = append(filterValues, f.Search)
		// query += `AND state = $` + strconv.Itoa(len(filterValues))
	}

	if !f.CustomerId.IsZero() {
		filterValues = append(filterValues, f.CustomerId.ValueOrZero())
		query += ` AND customer_id = $` + strconv.Itoa(len(filterValues))
	}

	if len(f.States) > 0 {
		query += ` AND (`
		for i, v := range f.States {
			s := strconv.Itoa(int(v))
			filterValues = append(filterValues, s)
			if i != 0 {
				query += ` OR`
			}
			query += ` tags.state = $` + strconv.Itoa(len(filterValues))
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

	// // For debugging purposes only.
	// log.Println("CustomerTagRepo | query:", query, "\n")
	// log.Println("CustomerTagRepo | filterValues:", filterValues, "\n")

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *CustomerTagRepo) ListByFilter(ctx context.Context, filter *models.CustomerTagFilter) ([]*models.CustomerTag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        customer_tags.id,
		customer_tags.tenant_id,
		customer_tags.customer_id,
		customer_tags.tag_id,
		tags.text,
		tags.description,
		tags.state
    FROM
        customer_tags
	INNER JOIN tags ON customer_tags.tag_id = tags.id
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.CustomerTag
	defer rows.Close()
	for rows.Next() {
		m := new(models.CustomerTag)
		err := rows.Scan(
			&m.Id,
			&m.TenantId,
			&m.CustomerId,
			&m.TagId,
			&m.Text,
			&m.Description,
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

func (s *CustomerTagRepo) CountByFilter(ctx context.Context, f *models.CustomerTagFilter) (uint64, error) {
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
	    customer_tags
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

	if !f.CustomerId.IsZero() {
		filterValues = append(filterValues, f.CustomerId.ValueOrZero())
		query += ` AND customer_id = $` + strconv.Itoa(len(filterValues))
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
