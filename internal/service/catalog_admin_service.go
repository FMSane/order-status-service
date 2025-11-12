package service

import (
	"context"
	"errors"
	"order-status-service/internal/model"
	"order-status-service/internal/repository"
	"time"
)

type CatalogAdminService struct {
	repo *repository.CatalogRepository
}

func NewCatalogAdminService(repo *repository.CatalogRepository) *CatalogAdminService {
	return &CatalogAdminService{repo: repo}
}

// Crea un nuevo estado en el catálogo base (solo admins)
func (s *CatalogAdminService) CreateStatus(name string) error {
	exists, err := s.repo.ExistsByName(name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("status already exists in catalog")
	}

	status := model.StatusCatalog{
		Name:      name,
		CreatedAt: time.Now(),
	}
	return s.repo.InsertOne(context.Background(), status)
}

// Devuelve todos los estados del catálogo base
func (s *CatalogAdminService) GetAll() ([]model.StatusCatalog, error) {
	return s.repo.GetAll()
}
