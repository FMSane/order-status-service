// catalog_repository.go
package repository

import (
	"context"
	"time"

	"order-status-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CatalogRepository struct {
	Collection *mongo.Collection
}

func NewCatalogRepository(db *mongo.Database) *CatalogRepository {
	return &CatalogRepository{
		Collection: db.Collection("statuses_catalog"),
	}
}

func (r *CatalogRepository) GetAll() ([]model.StatusCatalog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.StatusCatalog
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *CatalogRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.Collection.CountDocuments(ctx, bson.M{})
}

func (r *CatalogRepository) InsertMany(defaults []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.InsertMany(ctx, defaults)
	return err
}

func (r *CatalogRepository) ExistsByName(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := r.Collection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *CatalogRepository) InsertOne(ctx context.Context, status model.StatusCatalog) error {
	_, err := r.Collection.InsertOne(ctx, status)
	return err
}

func (r *CatalogRepository) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.Collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *CatalogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (model.StatusCatalog, error) {
	var res model.StatusCatalog
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	return res, err
}

func (r *CatalogRepository) FindByName(ctx context.Context, name string) (model.StatusCatalog, error) {
	var res model.StatusCatalog
	err := r.Collection.FindOne(ctx, bson.M{"name": name}).Decode(&res)
	return res, err
}

func (r *CatalogRepository) GetByID(id string) (*model.StatusCatalog, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result model.StatusCatalog
	err = r.Collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
