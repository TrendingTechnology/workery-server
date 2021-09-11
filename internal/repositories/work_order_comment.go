package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	// "log"

	"github.com/over55/workery-server/internal/models"
)

type WorkOrderCommentRepo struct {
	db *sql.DB
}

func NewWorkOrderCommentRepo(db *sql.DB) *WorkOrderCommentRepo {
	return &WorkOrderCommentRepo{
		db: db,
	}
}

func (r *WorkOrderCommentRepo) Insert(ctx context.Context, m *models.WorkOrderComment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO work_order_comments (
        uuid, tenant_id, order_id, comment_id, created_time, old_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.Uuid, m.TenantId, m.OrderId, m.CommentId, m.CreatedTime, m.OldId,
	)
	return err
}

func (r *WorkOrderCommentRepo) UpdateById(ctx context.Context, m *models.WorkOrderComment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        work_order_comments
    SET
        tenant_id = $1, order_id = $2, comment_id = $3
    WHERE
        id = $4`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantId, m.OrderId, m.CommentId, m.Id,
	)
	return err
}

func (r *WorkOrderCommentRepo) GetById(ctx context.Context, id uint64) (*models.WorkOrderComment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderComment)

	query := `
    SELECT
        id, uuid, tenant_id, order_id, comment_id, created_time
	FROM
        work_order_comments
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.OrderId, &m.CommentId, &m.CreatedTime,
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

func (r *WorkOrderCommentRepo) GetByOld(ctx context.Context, tenantId uint64, oldId uint64) (*models.WorkOrderComment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.WorkOrderComment)

	query := `
    SELECT
        id, uuid, tenant_id, order_id, comment_id, created_time
	FROM
        work_order_comments
    WHERE
        old_id = $1 AND tenant_id = $2`
	err := r.db.QueryRowContext(ctx, query, oldId, tenantId).Scan(
		&m.Id, &m.Uuid, &m.TenantId, &m.OrderId, &m.CommentId, &m.CreatedTime,
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

func (r *WorkOrderCommentRepo) CheckIfExistsById(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        work_order_comments
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

func (r *WorkOrderCommentRepo) InsertOrUpdateById(ctx context.Context, m *models.WorkOrderComment) error {
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

func (s *WorkOrderCommentRepo) queryRowsWithFilter(ctx context.Context, query string, f *models.WorkOrderCommentFilter) (*sql.Rows, error) {
	// Array will hold all the unique values we want to add into the query.
	var filterValues []interface{}

	// The SQL query statement we will be calling in the database, start
	// by setting the `tenant_id` placeholder and then append our value to
	// the array.
	filterValues = append(filterValues, f.TenantId)
	query += ` WHERE a.tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our OPTIONAL filters
	//

	if !f.CreatedTime.IsZero() {
		filterValues = append(filterValues, f.CreatedTime)
		query += ` AND a.created_time >= $` + strconv.Itoa(len(filterValues))
	}

	if len(f.States) > 0 {
		query += ` AND (`
		for i, v := range f.States {
			s := strconv.Itoa(int(v))
			filterValues = append(filterValues, s)
			if i != 0 {
				query += ` OR`
			}
			query += ` a.state = $` + strconv.Itoa(len(filterValues))
		}
		query += ` )`
	}

	//
	// The following code will add our pagination.
	//

	if f.LastSeenId > 0 {
		filterValues = append(filterValues, f.LastSeenId)
		query += ` AND a.id < $` + strconv.Itoa(len(filterValues))
	}
	query += ` ORDER BY a.created_time `
	filterValues = append(filterValues, f.Limit)
	query += ` DESC LIMIT $` + strconv.Itoa(len(filterValues))

	//
	// Execute our custom built SQL query to the database.
	//

	// log.Println("SERVE LOG | WorkOrderCommentRepo | queryRowsWithFilter | query", query)
	// log.Println("SERVE LOG | WorkOrderCommentRepo | queryRowsWithFilter | filterValues", filterValues)

	return s.db.QueryContext(ctx, query, filterValues...)
}

func (s *WorkOrderCommentRepo) ListByFilter(ctx context.Context, filter *models.WorkOrderCommentFilter) ([]*models.WorkOrderComment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	querySelect := `
    SELECT
        a.id,
		a.uuid,
		a.tenant_id,
		a.order_id,
		a.comment_id,
		a.created_time,
		b.text
    FROM
        work_order_comments
	AS
	    a
	JOIN
	    comments
	AS
	    b
	ON
	    a.comment_id = b.id
    `

	rows, err := s.queryRowsWithFilter(ctx, querySelect, filter)
	if err != nil {
		return nil, err
	}

	var arr []*models.WorkOrderComment
	defer rows.Close()
	for rows.Next() {
		m := new(models.WorkOrderComment)
		err := rows.Scan(
			&m.Id,
			&m.Uuid,
			&m.TenantId,
			&m.OrderId,
			&m.CommentId,
			&m.CreatedTime,
			&m.Text,
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

func (s *WorkOrderCommentRepo) CountByFilter(ctx context.Context, f *models.WorkOrderCommentFilter) (uint64, error) {
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
	    work_order_comments
	WHERE
		tenant_id = $` + strconv.Itoa(len(filterValues))

	//
	// The following code will add our filters
	//

	if !f.CreatedTime.IsZero() {
		filterValues = append(filterValues, f.CreatedTime)
		query += ` AND created_time >= $` + strconv.Itoa(len(filterValues))
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
