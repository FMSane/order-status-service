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

func NewOrderStatusRepository(db *mongo.Database) *OrderStatusRepository {
	return &OrderStatusRepository{
		Collection: db.Collection("order_statuses"),
	}
}

// Create inserts a new OrderStatus document
func (r *OrderStatusRepository) Create(ctx context.Context, status model.OrderStatus) error {
	if status.ID.IsZero() {
		status.ID = primitive.NewObjectID()
	}
	now := time.Now()
	if status.CreatedAt.IsZero() {
		status.CreatedAt = now
	}
	status.UpdatedAt = now
	_, err := r.Collection.InsertOne(ctx, status)
	return err
}

// FindByID retrieves an OrderStatus by its ObjectID
func (r *OrderStatusRepository) FindByID(ctx context.Context, id primitive.ObjectID) (model.OrderStatus, error) {
	var res model.OrderStatus
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	return res, err
}

// ExistsByOrderID checks if there's already a OrderStatus for an order_id
func (r *OrderStatusRepository) ExistsByOrderID(ctx context.Context, orderID string) (bool, error) {
	count, err := r.Collection.CountDocuments(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateStatusWithEntry atomically updates current status and pushes a history entry
func (r *OrderStatusRepository) UpdateStatusWithEntry(ctx context.Context, id primitive.ObjectID, statusID primitive.ObjectID, statusName string, entry model.StatusEntry) error {
	update := bson.M{
		"$set": bson.M{
			"status_id":  statusID,
			"status":     statusName,
			"updated_at": time.Now(),
		},
		"$push": bson.M{
			"history": entry,
		},
	}
	_, err := r.Collection.UpdateByID(ctx, id, update)
	return err
}

// GetBaseStatuses (returns distinct status names) - retained for compatibility
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

func (r *OrderStatusRepository) FindByStatusID(ctx context.Context, statusID primitive.ObjectID) ([]model.OrderStatus, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"status_id": statusID})
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
