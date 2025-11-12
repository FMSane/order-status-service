package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Modelo de los estados base del catálogo
type StatusCatalog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
