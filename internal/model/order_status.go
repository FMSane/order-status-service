package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID   string             `bson:"order_id" json:"order_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Status    string             `bson:"status" json:"status"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
