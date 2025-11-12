package repository

import (
	"context"
	"order-status-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderStatusRepository struct {
	Collection *mongo.Collection
}

// ✅ Esta es la única función que debe existir
func NewOrderStatusRepository(db *mongo.Database) *OrderStatusRepository {
	return &OrderStatusRepository{
		Collection: db.Collection("order_statuses"),
	}
}

func (r *OrderStatusRepository) Create(ctx context.Context, status model.OrderStatus) error {
	if status.ID.IsZero() {
		status.ID = primitive.NewObjectID()
	}
	if status.UpdatedAt.IsZero() {
		status.UpdatedAt = time.Now()
	}
	_, err := r.Collection.InsertOne(ctx, status)
	return err
}

func (r *OrderStatusRepository) Update(ctx context.Context, id string, newStatus string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}
	_, err = r.Collection.UpdateByID(ctx, objID, update)
	return err
}

func (r *OrderStatusRepository) GetBaseStatuses(ctx context.Context) ([]string, error) {
	cursor, err := r.Collection.Distinct(ctx, "status", bson.D{})
	if err != nil {
		return nil, err
	}
	names := make([]string, len(cursor))
	for i, v := range cursor {
		names[i] = v.(string)
	}
	return names, nil
}

func (r *OrderStatusRepository) FindAll(ctx context.Context) ([]model.OrderStatus, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []model.OrderStatus
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *OrderStatusRepository) FindByUser(ctx context.Context, userID string) ([]model.OrderStatus, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.OrderStatus
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *OrderStatusRepository) ExistsByName(name string) (bool, error) {
	count, err := r.Collection.CountDocuments(context.TODO(), bson.M{"status": name})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrderStatusRepository) Insert(ctx context.Context, status model.OrderStatus) error {
	_, err := r.Collection.InsertOne(ctx, status)
	return err
}

func (r *OrderStatusRepository) FindByStatus(ctx context.Context, status string) ([]model.OrderStatus, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []model.OrderStatus
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
