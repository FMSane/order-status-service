package service

import (
	"context"
	"errors"
	"fmt"
	"order-status-service/internal/dto"
	"order-status-service/internal/mapper"
	"order-status-service/internal/model"
	"order-status-service/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderStatusService es el servicio principal para manejar estados de órdenes
type OrderStatusService struct {
	repo *repository.OrderStatusRepository
	// nota: requiere que CatalogRepository tenga los métodos FindByID / ExistsByID
	catalogRepo *repository.CatalogRepository
}

func NewOrderStatusService(repo *repository.OrderStatusRepository) *OrderStatusService {
	return &OrderStatusService{
		repo:        repo,
		catalogRepo: repository.NewCatalogRepository(repo.Collection.Database()),
	}
}

// CreateStatus crea un nuevo documento OrderStatus (usado para inicialización)
// Acepta StatusID (preferido) o Status (nombre) en la request.
func (s *OrderStatusService) CreateStatus(req dto.CreateOrderStatusRequest) (dto.OrderStatusDTO, error) {
	ctx := context.Background()

	// Resolver id y nombre del estado: priorizar StatusID si se envía
	var statusID primitive.ObjectID
	var statusName string
	if req.StatusID != "" {
		id, err := primitive.ObjectIDFromHex(req.StatusID)
		if err != nil {
			return dto.OrderStatusDTO{}, fmt.Errorf("invalid status_id: %v", err)
		}
		// verificar existencia en catálogo
		exists, err := s.catalogRepo.ExistsByID(ctx, id)
		if err != nil {
			return dto.OrderStatusDTO{}, err
		}
		if !exists {
			return dto.OrderStatusDTO{}, fmt.Errorf("status id %s not found in catalog", req.StatusID)
		}
		statusID = id
		// obtener nombre
		cat, err := s.catalogRepo.FindByID(ctx, id)
		if err != nil {
			return dto.OrderStatusDTO{}, err
		}
		statusName = cat.Name
	} else if req.Status != "" {
		// alternativa: aceptar nombre de estado (compatibilidad hacia atrás)
		exists, err := s.catalogRepo.ExistsByName(req.Status)
		if err != nil {
			return dto.OrderStatusDTO{}, err
		}
		if !exists {
			return dto.OrderStatusDTO{}, fmt.Errorf("status '%s' does not exist in catalog", req.Status)
		}
		// buscar id por nombre
		cat, err := s.catalogRepo.FindByName(ctx, req.Status)
		if err != nil {
			return dto.OrderStatusDTO{}, err
		}
		statusID = cat.ID
		statusName = cat.Name
	} else {
		// valor por defecto: Pendiente
		cat, err := s.catalogRepo.FindByName(ctx, "Pendiente")
		if err != nil {
			return dto.OrderStatusDTO{}, errors.New("default status 'Pendiente' not found")
		}
		statusID = cat.ID
		statusName = cat.Name
	}

	// Prevenir múltiples order_status para la misma orden (idempotencia en inicialización)
	if req.OrderID != "" {
		exists, err := s.repo.ExistsByOrderID(context.Background(), req.OrderID)
		if err != nil {
			return dto.OrderStatusDTO{}, err
		}
		if exists {
			// buscar documento existente y retornarlo
			// intento de búsqueda por order id (no existe método directo FindByOrderID, se usa FindAll)
			all, err := s.repo.FindAll(context.Background())
			if err == nil {
				for _, o := range all {
					if o.OrderID == req.OrderID {
						return mapper.ToOrderStatusDTO(o), nil
					}
				}
			}
			// alternativa: conflicto
			return dto.OrderStatusDTO{}, fmt.Errorf("order_status for order_id %s already exists", req.OrderID)
		}
	}

	// Construir entidad
	entity := model.OrderStatus{
		ID:        primitive.NewObjectID(),
		OrderID:   req.OrderID,
		UserID:    req.UserID,
		StatusID:  statusID,
		Status:    statusName,
		Shipping:  mapper.ToShippingEntity(req.Shipping),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		History: []model.StatusEntry{
			{
				ID:     primitive.NewObjectID(),
				Status: statusName,
				UserID: req.UserID,
				Role:   "system",
				Reason: "initial",
				At:     time.Now(),
			},
		},
	}

	if err := s.repo.Create(context.Background(), entity); err != nil {
		return dto.OrderStatusDTO{}, err
	}
	return mapper.ToOrderStatusDTO(entity), nil
}

// ChangeStatus cambia el estado actual aplicando reglas de negocio
func (s *OrderStatusService) ChangeStatus(orderStatusID string, newStatusID string, actorID string, actorRole string, reason string) (dto.OrderStatusDTO, error) {
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(orderStatusID)
	if err != nil {
		return dto.OrderStatusDTO{}, fmt.Errorf("invalid order status id")
	}
	newID, err := primitive.ObjectIDFromHex(newStatusID)
	if err != nil {
		return dto.OrderStatusDTO{}, fmt.Errorf("invalid new status id")
	}

	// resolver el nombre del nuevo estado desde el catálogo
	cat, err := s.catalogRepo.FindByID(ctx, newID)
	if err != nil {
		return dto.OrderStatusDTO{}, err
	}
	newName := cat.Name

	// buscar documento existente
	doc, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	// Idempotencia: si es el mismo status_id -> retornar sin cambios
	if doc.StatusID == newID {
		return mapper.ToOrderStatusDTO(doc), nil
	}

	// Estados terminales donde no se permiten cambios
	terminal := map[string]bool{
		"Cancelado": true,
		"Entregado": true,
		"Rechazado": true,
	}
	if terminal[doc.Status] {
		return dto.OrderStatusDTO{}, fmt.Errorf("cannot change status from terminal state '%s'", doc.Status)
	}

	// REGLAS DE NEGOCIO
	// Cancelado -> solo cliente y solo si el estado actual NO es Enviado/Entregado/Rechazado
	if newName == "Cancelado" {
		if actorRole != "client" {
			return dto.OrderStatusDTO{}, fmt.Errorf("only client can cancel the order")
		}
		forbidden := map[string]bool{"Enviado": true, "Entregado": true, "Rechazado": true}
		if forbidden[doc.Status] {
			return dto.OrderStatusDTO{}, fmt.Errorf("cannot cancel when current status is '%s'", doc.Status)
		}
	}

	// Rechazado -> solo admin o seller y solo si el estado actual NO es Enviado o Cancelado
	if newName == "Rechazado" {
		if actorRole != "admin" && actorRole != "seller" {
			return dto.OrderStatusDTO{}, fmt.Errorf("only admin or seller can reject the order")
		}
		forbidden := map[string]bool{"Enviado": true, "Cancelado": true}
		if forbidden[doc.Status] {
			return dto.OrderStatusDTO{}, fmt.Errorf("cannot reject when current status is '%s'", doc.Status)
		}
	}

	// Todas las validaciones pasaron — construir entrada de historial y actualizar
	entry := model.StatusEntry{
		ID:     primitive.NewObjectID(),
		Status: newName,
		UserID: actorID,
		Role:   actorRole,
		Reason: reason,
		At:     time.Now(),
	}

	if err := s.repo.UpdateStatusWithEntry(ctx, objID, newID, newName, entry); err != nil {
		return dto.OrderStatusDTO{}, err
	}

	// retornar documento actualizado (volver a buscar)
	updated, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return mapper.ToOrderStatusDTO(updated), nil
}

// Otros getters auxiliares reutilizando el repo
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

func (s *OrderStatusService) GetByStatusID(statusID string) ([]dto.OrderStatusDTO, error) {
	objID, err := primitive.ObjectIDFromHex(statusID)
	if err != nil {
		return nil, fmt.Errorf("invalid status id")
	}

	statuses, err := s.repo.FindByStatusID(context.Background(), objID)
	if err != nil {
		return nil, err
	}

	return mapper.ToOrderStatusDTOs(statuses), nil
}

func (s *OrderStatusService) IsCancelStatus(id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid status id")
	}

	st, err := s.catalogRepo.FindByID(context.Background(), objID)
	if err != nil {
		return false, err
	}

	return st.Name == "Cancelado", nil
}
