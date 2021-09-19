package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/over55/workery-server/internal/models"
)

type BulletinBoardItemRepo struct {
	db *sql.DB
}

func NewBulletinBoardItemRepo(db *sql.DB) *BulletinBoardItemRepo {
	return &BulletinBoardItemRepo{
		db: db,
	}
}

func (r *BulletinBoardItemRepo) Insert(ctx context.Context, m *models.BulletinBoardItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO bulletin_board_items (
        uuid, tenant_id, text, created_time, created_by_id, created_by_name, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_by_name, last_modified_from_ip, state,
		old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.Text, m.CreatedTime, m.CreatedById, m.CreatedByName,
		m.CreatedFromIP, m.LastModifiedTime, m.LastModifiedById, m.LastModifiedByName,
		m.LastModifiedFromIP, m.State, m.OldId,
	)
	return err
}

func (r *BulletinBoardItemRepo) UpdateById(ctx context.Context, m *models.BulletinBoardItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        bulletin_board_items
    SET
        tenant_id = $1, text = $2, created_time = $3, created_by_id = $4,
		created_from_ip = $5, last_modified_time = $6, last_modified_by_id = $7,
		last_modified_from_ip = $8, state = $9
    WHERE
        id = $10`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.Text, m.CreatedTime, m.CreatedById,
		m.CreatedFromIP, m.LastModifiedTime, m.LastModifiedById,
		m.LastModifiedFromIP, m.State, m.Id,
	)
	return err
}

func (r *BulletinBoardItemRepo) GetById(ctx context.Context, id uint64) (*models.BulletinBoardItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.BulletinBoardItem)

	query := `
    SELECT
        id, uuid, tenant_id, text, created_time, created_by_id, created_from_ip,
		last_modified_time, last_modified_by_id, last_modified_from_ip, state
	FROM
        bulletin_board_items
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.CreatedTime, &m.CreatedById,
		&m.CreatedFromIP, &m.LastModifiedTime, &m.LastModifiedById,
		&m.LastModifiedFromIP, &m.State,
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

func (r *BulletinBoardItemRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        bulletin_board_items
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

func (r *BulletinBoardItemRepo) InsertOrUpdateById(ctx context.Context, m *models.BulletinBoardItem) error {
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

func (r *BulletinBoardItemRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.BulletinBoardItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.BulletinBoardItem)

	query := `
    SELECT
	    id, uuid, tenant_id, text, created_time, created_by_id, created_from_ip,
    	last_modified_time, last_modified_by_id, last_modified_from_ip, state
    FROM
	    bulletin_board_items
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.Text, &m.CreatedTime, &m.CreatedById,
		&m.CreatedFromIP, &m.LastModifiedTime, &m.LastModifiedById,
		&m.LastModifiedFromIP, &m.State,
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

func (s *BulletinBoardItemRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.BulletinBoardItemFilter) (*sql.Rows, error) {
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
	query += ` ORDER BY created_time `
	filterValues = append(filterValues, f.Limit)
	query += ` DESC LIMIT $` + strconv.Itoa(len(filterValues))

	//
	// Execute our custom built SQL query to the database.
	//

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *BulletinBoardItemRepo) ListByFilter(ctx context.Context, filter *models.BulletinBoardItemFilter) ([]*models.BulletinBoardItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
		id,
		uuid,
		tenant_id,
		text,
		created_time,
		created_by_id,
		created_from_ip,
		last_modified_time,
		last_modified_by_id,
		last_modified_from_ip,
		state
    FROM
        bulletin_board_items
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.BulletinBoardItem
	defer rows.Close()
	for rows.Next() {
		m := new(models.BulletinBoardItem)
		err := rows.Scan(
			&m.Id,
			&m.Uuid,
			&m.TenantId,
			&m.Text,
			&m.CreatedTime,
			&m.CreatedById,
			&m.CreatedFromIP,
			&m.LastModifiedTime,
			&m.LastModifiedById,
			&m.LastModifiedFromIP,
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

func (s *BulletinBoardItemRepo) CountByFilter(ctx context.Context, f *models.BulletinBoardItemFilter) (uint64, error) {
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
	    bulletin_board_items
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

	// Return our values.
	return count, err
}
