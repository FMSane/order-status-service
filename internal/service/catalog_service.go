package service

import (
	"log"
	"time"

	"order-status-service/internal/model"
	"order-status-service/internal/repository"
)

type CatalogService struct {
	Repo *repository.CatalogRepository
}

func NewCatalogService(repo *repository.CatalogRepository) *CatalogService {
	return &CatalogService{Repo: repo}
}

// Se ejecuta automáticamente al iniciar el microservicio
func (s *CatalogService) SeedDefaultStatuses() error {
	count, err := s.Repo.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("✅ Estados base ya existen, no se vuelven a crear")
		return nil
	}

	defaults := []interface{}{
		model.StatusCatalog{Name: "Pendiente", CreatedAt: time.Now()},
		model.StatusCatalog{Name: "En preparación", CreatedAt: time.Now()},
		model.StatusCatalog{Name: "Enviado", CreatedAt: time.Now()},
		model.StatusCatalog{Name: "Entregado", CreatedAt: time.Now()},
		model.StatusCatalog{Name: "Cancelado", CreatedAt: time.Now()},
	}

	if err := s.Repo.InsertMany(defaults); err != nil {
		return err
	}

	log.Println("✅ Estados base creados automáticamente")
	return nil
}

func (s *CatalogService) GetAll() ([]model.StatusCatalog, error) {
	return s.Repo.GetAll()
}
