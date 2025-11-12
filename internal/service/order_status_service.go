package service

import (
	"context"
	"fmt"
	"order-status-service/internal/dto"
	"order-status-service/internal/mapper"
	"order-status-service/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatusService struct {
	repo *repository.OrderStatusRepository
}

func NewOrderStatusService(repo *repository.OrderStatusRepository) *OrderStatusService {
	return &OrderStatusService{repo: repo}
}

func (s *OrderStatusService) SeedDefaultStatuses() error {
	defaults := []string{"Pendiente", "En preparación", "Enviado", "Entregado", "Cancelado"}

	for _, name := range defaults {
		exists, err := s.repo.ExistsByName(name)
		if err != nil {
			return err
		}
		if !exists {
			_, err := s.CreateStatus(dto.CreateOrderStatusRequest{Status: name})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *OrderStatusService) CreateStatus(req dto.CreateOrderStatusRequest) (dto.OrderStatusDTO, error) {

	// Validar existencia del estado en el catálogo
	catalogRepo := repository.NewCatalogRepository(s.repo.Collection.Database())
	exists, err := catalogRepo.ExistsByName(req.Status)
	if err != nil {
		return dto.OrderStatusDTO{}, err
	}
	if !exists {
		return dto.OrderStatusDTO{}, fmt.Errorf("status '%s' does not exist in catalog", req.Status)
	}

	entity := mapper.ToOrderStatusEntity(req)

	// Asegurar que siempre haya ID y fecha
	if entity.ID.IsZero() {
		entity.ID = primitive.NewObjectID()
	}
	entity.UpdatedAt = time.Now()

	err = s.repo.Create(context.Background(), entity)
	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return mapper.ToOrderStatusDTO(entity), nil
}

func (s *OrderStatusService) UpdateStatus(id string, req dto.UpdateOrderStatusRequest) (dto.OrderStatusDTO, error) {
	err := s.repo.Update(context.Background(), id, req.Status)
	if err != nil {
		return dto.OrderStatusDTO{}, err
	}
	// no tenemos el documento actualizado, así que devolvemos uno mínimo
	return dto.OrderStatusDTO{ID: id, Status: req.Status}, nil
}

func (s *OrderStatusService) GetAllStatuses() ([]string, error) {
	return s.repo.GetBaseStatuses(context.Background())
}

func (s *OrderStatusService) GetAllOrderStatuses() ([]dto.OrderStatusDTO, error) {
	statuses, err := s.repo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	return mapper.ToOrderStatusDTOs(statuses), nil
}

func (s *OrderStatusService) GetStatusesByUser(userID string) ([]dto.OrderStatusDTO, error) {
	statuses, err := s.repo.FindByUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return mapper.ToOrderStatusDTOs(statuses), nil
}

func (s *OrderStatusService) GetByStatus(status string) ([]dto.OrderStatusDTO, error) {
	statuses, err := s.repo.FindByStatus(context.Background(), status)
	if err != nil {
		return nil, err
	}
	return mapper.ToOrderStatusDTOs(statuses), nil
}
