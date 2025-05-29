package tests

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"user-service/internal/core/domain"
	"user-service/internal/infrastructure/database"
)

func TestMongoUserCreation(t *testing.T) {
	client, err := database.NewMongoClient("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.CloseMongoClient(client)

	collection := client.Database("user_service").Collection("users")

	testUser := &domain.User{
		ID:       "test123",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, testUser)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	var foundUser domain.User
	err = collection.FindOne(ctx, bson.M{"_id": testUser.ID}).Decode(&foundUser)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if foundUser.Email != testUser.Email {
		t.Errorf("Expected email %s, got %s", testUser.Email, foundUser.Email)
	}
	if foundUser.Name != testUser.Name {
		t.Errorf("Expected name %s, got %s", testUser.Name, foundUser.Name)
	}
	if foundUser.Password != testUser.Password {
		t.Errorf("Expected password %s, got %s", testUser.Password, foundUser.Password)
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": testUser.ID})
	if err != nil {
		t.Errorf("Failed to clean up test user: %v", err)
	}
}
