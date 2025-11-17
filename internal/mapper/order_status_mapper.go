// order_status_mapper.go
package mapper

import (
	"order-status-service/internal/dto"
	"order-status-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Convierte un DTO de creación a entidad de base de datos
func ToOrderStatusEntity(req dto.CreateOrderStatusRequest) model.OrderStatus {
	return model.OrderStatus{
		ID:        primitive.NewObjectID(),
		OrderID:   req.OrderID,
		UserID:    req.UserID,
		Status:    req.Status,
		UpdatedAt: time.Now(),
	}
}

// Convierte una entidad de base de datos a DTO
func ToOrderStatusDTO(entity model.OrderStatus) dto.OrderStatusDTO {
	return dto.OrderStatusDTO{
		ID:        entity.ID.Hex(),
		OrderID:   entity.OrderID,
		UserID:    entity.UserID,
		Status:    entity.Status,
		UpdatedAt: entity.UpdatedAt,
	}
}

// Convierte múltiples entidades de base de datos a DTOs
func ToOrderStatusDTOs(entities []model.OrderStatus) []dto.OrderStatusDTO {
	dtos := make([]dto.OrderStatusDTO, len(entities))
	for i, e := range entities {
		dtos[i] = ToOrderStatusDTO(e)
	}
	return dtos
}

func ToShippingEntity(s dto.ShippingDTO) model.ShippingInfo {
	return model.ShippingInfo{
		AddressLine1: s.AddressLine1,
		AddressLine2: s.AddressLine2,
		City:         s.City,
		Province:     s.Province,
		Country:      s.Country,
		Zipcode:      s.Zipcode,
		Comments:     s.Comments,
	}
}
