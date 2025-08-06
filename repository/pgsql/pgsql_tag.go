package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/domain"
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
