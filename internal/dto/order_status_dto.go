package dto

import "time"

// Request para crear un nuevo estado de orden
type CreateOrderStatusRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
	Status  string `json:"status"`
}

// Request para actualizar el estado de una orden
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Response devuelto al listar estados
type OrderStatusDTO struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	OrderID   string    `json:"order_id" bson:"order_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Status    string    `json:"status" bson:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}
