package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"
	"strings"
)

type pgsqlTagRepository struct {
	db *sql.DB
}

// NewPgsqlTagRepository will create new an todoRepository object representation of TagRepository interface
func NewPgsqlTagRepository(db *sql.DB) *pgsqlTagRepository {
	return &pgsqlTagRepository{
		db: db,
	}
}

func (r *pgsqlTagRepository) Create(ctx context.Context, tag *domain.Tag) (err error) {
	query := `INSERT INTO tags (id, name, description, is_active, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err = r.db.ExecContext(ctx, query, tag.ID, tag.Name, tag.Description, tag.IsActive, tag.CreatedBy, tag.CreatedAt); err != nil {
		return err
	}

	return
}

func (r *pgsqlTagRepository) Update(ctx context.Context, tag *domain.Tag) (err error) {
	// Build dynamic SET clauses from Tag struct
	sets := []string{}
	args := []interface{}{}
	idx := 1

	if tag.Name != "" {
		sets = append(sets, fmt.Sprintf("name = $%d", idx))
		args = append(args, tag.Name)
		idx++
	}

	if tag.Description != "" {
		sets = append(sets, fmt.Sprintf("description = $%d", idx))
		args = append(args, tag.Description)
		idx++
	}

	if tag.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, tag.IsActive)
		idx++
	}

	// kalau ada sesuatu untuk di‐update, commit ke SQL
	if len(sets) > 0 {
		// Update stamp
		sets = append(sets, fmt.Sprintf("updated_at = $%d", idx))
		args = append(args, tag.UpdatedAt)
		idx++

		sets = append(sets, fmt.Sprintf("updated_by = $%d", idx))
		args = append(args, tag.UpdatedBy)
		idx++

		// tambahkan WHERE id = $idx
		args = append(args, tag.ID)
		query := fmt.Sprintf(
			"UPDATE tags SET %s WHERE id = $%d",
			strings.Join(sets, ", "),
			idx,
		)

		if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
			return
		}
	}

	return
}

func (r *pgsqlTagRepository) Delete(ctx context.Context, tag *domain.Tag) (rowsAffected int64, err error) {
	query := "DELETE FROM tags WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, tag.ID)
	if err != nil {
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	return
}

func (r *pgsqlTagRepository) GetList(ctx context.Context, request *request.GetListTagReq) (res []response.GetListTagRes, meta response.MetaRes, err error) {
	// 1. Build WHERE clauses
	wheres := []string{}
	args := []interface{}{}
	idx := 1

	if request.Name != "" {
		wheres = append(wheres, fmt.Sprintf("name ILIKE $%d", idx))
		args = append(args, "%"+request.Name+"%")
		idx++
	}

	if request.IsActive != nil {
		wheres = append(wheres, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, request.IsActive)
		idx++
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// --- 2. Hitung totalCount dulu (tanpa LIMIT/OFFSET) ---
	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM tags %s",
		whereSQL,
	)
	if err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&meta.TotalData); err != nil {
		return nil, meta, err
	}

	// 2. Calculate LIMIT & OFFSET
	perPage := request.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	page := request.Page
	if page <= 0 {
		page = 1
	}

	// total pages = ceil(total / perPage)
	meta.Page = page
	meta.PerPage = perPage
	meta.TotalPages = (meta.TotalData + perPage - 1) / perPage

	offset := (page - 1) * perPage

	// add LIMIT & OFFSET to args
	args = append(args, perPage, offset)
	limitPos, offsetPos := idx, idx+1

	// 3. Final query
	query := fmt.Sprintf(`
        SELECT
            id, name, description,
            is_active
        FROM tags
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereSQL, limitPos, offsetPos)

	// 4. Execute
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()

	// 5. Scan results
	for rows.Next() {
		var item response.GetListTagRes

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.IsActive,
		); err != nil {
			return nil, meta, err
		}

		res = append(res, item)
	}
	if errRow := rows.Err(); errRow != nil {
		return nil, meta, errRow
	}

	return
}

func (r *pgsqlTagRepository) GetDetail(ctx context.Context, request *request.GetDetailTagReq) (res domain.Tag, err error) {

	const query = `
					SELECT
					  id,
					  name,
					  description,
					  is_active,
					  created_at,
					  created_by,
					  updated_at,
					  updated_by,
					  deleted_at
					FROM tags
					WHERE id = $1
					LIMIT 1
					`

	// 1. QueryRowContext untuk ambil satu baris
	row := r.db.QueryRowContext(ctx, query, request.ID)

	// 2. Scan kolom ke field di domain.Tag
	// since created_at is NOT NULL int8:
	var createdAt int64
	// updated_at/deleted_at can be NULL, so use NullInt64:
	var updatedAt, deletedAt sql.NullInt64
	var updatedBy sql.NullString

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.Description,
		&res.IsActive,
		&createdAt,
		&res.CreatedBy,
		&updatedAt,
		&updatedBy,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, utils.NewNotFoundError("partner not found")
		}
		return res, err
	}

	// assign into your domain fields
	res.CreatedAt = createdAt
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Int64
	}
	if deletedAt.Valid {
		res.DeletedAt = deletedAt.Int64
	}
	if updatedBy.Valid {
		res.UpdatedBy = updatedBy.String
	}

	return
}
