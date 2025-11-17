package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShippingInfo struct {
	AddressLine1 string `bson:"address_line1" json:"address_line1"`
	AddressLine2 string `bson:"address_line2,omitempty" json:"address_line2,omitempty"`
	City         string `bson:"city" json:"city"`
	Province     string `bson:"province,omitempty" json:"province,omitempty"`
	Country      string `bson:"country" json:"country"`
	Zipcode      string `bson:"zipcode,omitempty" json:"zipcode,omitempty"`
	Comments     string `bson:"comments,omitempty" json:"comments,omitempty"`
}

type StatusEntry struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status string             `bson:"status" json:"status"`
	UserID string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Role   string             `bson:"role,omitempty" json:"role,omitempty"`
	Reason string             `bson:"reason,omitempty" json:"reason,omitempty"`
	At     time.Time          `bson:"at" json:"at"`
}

type OrderStatus struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID string             `bson:"order_id" json:"order_id"`
	UserID  string             `bson:"user_id" json:"user_id"`
	// Nuevo: referenciamos el estado por id y mantenemos un nombre legible
	StatusID  primitive.ObjectID `bson:"status_id" json:"status_id"`
	Status    string             `bson:"status" json:"status"`
	Shipping  ShippingInfo       `bson:"shipping" json:"shipping"`
	History   []StatusEntry      `bson:"history" json:"history"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
