package repository

import (
	"context"
	"log"
	"shopping-cart-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartRepository struct {
	collection *mongo.Collection
}

func NewCartRepository(db *mongo.Database) *CartRepository {
	return &CartRepository{
		collection: db.Collection("cart"),
	}
}

func (r *CartRepository) AddToCart(ctx context.Context, item model.CartItem) error {
	log.Printf("[DB] ➕ AddToCart: user_id=%s, product_id=%s, qty=%d", item.UserID, item.ProductID, item.Quantity)

	filter := bson.M{"user_id": item.UserID, "product_id": item.ProductID}
	update := bson.M{"$inc": bson.M{"quantity": item.Quantity}}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("[DB] AddToCart error: %v", err)
		return err
	}

	log.Printf("[DB] Updated: matched=%d, modified=%d", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (r *CartRepository) GetCart(ctx context.Context, userID string) ([]model.CartItem, error) {
	log.Printf("[DB] 📥 GetCart: user_id=%s", userID)

	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("[DB] GetCart error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []model.CartItem
	if err = cursor.All(ctx, &items); err != nil {
		log.Printf("[DB] Cursor error: %v", err)
		return nil, err
	}

	log.Printf("[DB] Found %d items", len(items))
	return items, nil
}

func (r *CartRepository) RemoveFromCart(ctx context.Context, userID string, productID string) error {
	log.Printf("[DB] RemoveFromCart: user_id=%s, product_id=%s", userID, productID)

	filter := bson.M{"user_id": userID, "product_id": productID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("[DB] RemoveFromCart error: %v", err)
	}
	return err
}
