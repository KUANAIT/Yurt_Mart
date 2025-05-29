package repository

import (
	"context"
	"log"
	"time"

	"meaningfullname/Yurt_Mart/common/database"
	"meaningfullname/Yurt_Mart/orderhistory/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Order, error)
	GetRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.Order, error)
	Delete(ctx context.Context, id string) error
}

type orderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository() OrderRepository {
	if database.GetDB() == nil {
		log.Println("database.GetDB() returned nil in NewOrderRepository!")
		// Depending on desired behavior, you might return nil or an error here
		// For now, we'll proceed and expect the service layer to handle the nil repo
		return &orderRepository{}
	}
	return &orderRepository{
		collection: database.GetDB().Collection("orders"),
	}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	order.BeforeCreate() // Generate ID if empty
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	var orders []*model.Order

	filter := bson.M{"user_id": userID}

	// Options to sort by creation date descending
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) GetRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.Order, error) {
	var orders []*model.Order

	filter := bson.M{"user_id": userID}

	// Options to sort by creation date descending and limit the results
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
