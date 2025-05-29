package repository

import (
	"context"
	"log"
	"product-service/database"
	"product-service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(product *domain.Product) (string, error)
	GetProductsByCategory(category string) ([]*domain.Product, error)
	GetByID(id string) (*domain.Product, error)

	Delete(id string) error
	List() ([]*domain.Product, error)
}

type mongoRepo struct {
	collection *mongo.Collection
}

func NewMongoProductRepository() ProductRepository {
	client := database.ConnectMongo("mongodb://localhost:27017")
	collection := client.Database("onlinesupermarket").Collection("products")
	return &mongoRepo{collection: collection}
}

func (m *mongoRepo) Create(product *domain.Product) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := m.collection.InsertOne(ctx, product)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", err
	}

	return oid.Hex(), nil
}

func (m *mongoRepo) GetByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product domain.Product
	err = m.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (m *mongoRepo) GetProductsByCategory(category string) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"category": category}
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err == nil {
			products = append(products, &p)
		}
	}

	return products, nil
}

func (m *mongoRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (m *mongoRepo) List() ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			log.Println("decode error:", err)
			continue
		}
		products = append(products, &p)
	}

	return products, nil
}
