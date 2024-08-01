package repository

import (
	"FinUslugi/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MaterialRepository struct {
	DB *pgxpool.Pool
}

func NewMaterialRepository(db *pgxpool.Pool) *MaterialRepository {
	return &MaterialRepository{DB: db}
}

func (r *MaterialRepository) Create(ctx context.Context, material *models.Material) error {
	query := `INSERT INTO materials (uuid, type_materials, status, title, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.DB.Exec(ctx, query, material.UUID, material.TypeMaterials, material.Status, material.Title, material.Content, material.CreatedAt, material.UpdatedAt)
	return err
}

func (r *MaterialRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Material, error) {
	query := `SELECT uuid, type_materials, status, title, content, created_at, updated_at FROM materials WHERE uuid = $1`
	row := r.DB.QueryRow(ctx, query, id)

	var material models.Material
	err := row.Scan(&material.UUID, &material.TypeMaterials, &material.Status, &material.Title, &material.Content, &material.CreatedAt, &material.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &material, err
}

func (r *MaterialRepository) Update(ctx context.Context, material *models.Material) error {
	query := `UPDATE materials SET status = $1, title = $2, content = $3, updated_at = $4 WHERE uuid = $5`
	_, err := r.DB.Exec(ctx, query, material.Status, material.Title, material.Content, material.UpdatedAt, material.UUID)
	return err
}

func (r *MaterialRepository) GetAll(ctx context.Context, typeMaterials *string, createdAtFrom *time.Time, createdAtTo *time.Time, limit, offset int) ([]*models.Material, error) {
	query := `SELECT uuid, type_materials, status, title, content, created_at, updated_at FROM materials WHERE status = 'активный'`
	args := []interface{}{}
	argCount := 1

	if typeMaterials != nil {
		query += fmt.Sprintf(" AND type_materials = $%d", argCount)
		args = append(args, *typeMaterials)
		argCount++
	}

	if createdAtFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *createdAtFrom)
		argCount++
	}

	if createdAtTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *createdAtTo)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []*models.Material
	for rows.Next() {
		var material models.Material
		if err := rows.Scan(&material.UUID, &material.TypeMaterials, &material.Status, &material.Title, &material.Content, &material.CreatedAt, &material.UpdatedAt); err != nil {
			return nil, err
		}
		materials = append(materials, &material)
	}
	return materials, nil
}
