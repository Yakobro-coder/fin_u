package services

import (
	"FinUslugi/internal/models"
	"FinUslugi/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type MaterialService struct {
	repo   *repository.MaterialRepository
	logger *zap.SugaredLogger
}

func NewMaterialService(db *pgxpool.Pool, logger *zap.SugaredLogger) *MaterialService {
	repo := repository.NewMaterialRepository(db)
	return &MaterialService{repo: repo, logger: logger}
}

func (s *MaterialService) CreateMaterial(materialType, status, title, content string) (*models.Material, error) {
	material := &models.Material{
		UUID:          uuid.New(),
		TypeMaterials: materialType,
		Status:        status,
		Title:         title,
		Content:       content,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := s.repo.Create(context.Background(), material)
	return material, err
}

func (s *MaterialService) GetMaterialByID(id uuid.UUID) (*models.Material, error) {
	return s.repo.GetByID(context.Background(), id)
}

func (s *MaterialService) UpdateMaterial(id uuid.UUID, status, title, content string) (*models.Material, error) {
	material, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	material.Status = status
	material.Title = title
	material.Content = content
	material.UpdatedAt = time.Now()

	err = s.repo.Update(context.Background(), material)
	return material, err
}

func (s *MaterialService) GetAllMaterials(typeMaterials *string, createdAtFrom *time.Time, createdAtTo *time.Time, limit, offset int) ([]*models.Material, error) {
	return s.repo.GetAll(context.Background(), typeMaterials, createdAtFrom, createdAtTo, limit, offset-1)
}
