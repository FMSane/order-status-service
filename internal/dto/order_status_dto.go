package dto

import "time"

type ShippingDTO struct {
	AddressLine1 string `json:"address_line1" binding:"required"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city" binding:"required"`
	Province     string `json:"province,omitempty"`
	Country      string `json:"country" binding:"required"`
	Zipcode      string `json:"zipcode,omitempty"`
	Comments     string `json:"comments,omitempty"`
}

// CreateOrderStatusRequest para init
type CreateOrderStatusRequest struct {
	OrderID  string      `json:"order_id" binding:"required"`
	UserID   string      `json:"user_id" binding:"required"`
	Status   string      `json:"status"`              // opcional por compat
	StatusID string      `json:"status_id,omitempty"` // prefiero este
	Shipping ShippingDTO `json:"shipping" binding:"required"`
}

// Update request: ahora pedimos status_id
type UpdateOrderStatusRequest struct {
	StatusID string `json:"status_id" binding:"required"`
	Reason   string `json:"reason,omitempty"`
}

// DTO de respuesta
type OrderStatusDTO struct {
	ID        string      `json:"id" bson:"_id,omitempty"`
	OrderID   string      `json:"order_id" bson:"order_id"`
	UserID    string      `json:"user_id" bson:"user_id"`
	StatusID  string      `json:"status_id" bson:"status_id"`
	Status    string      `json:"status" bson:"status"`
	Shipping  ShippingDTO `json:"shipping"`
	History   []any       `json:"history,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
