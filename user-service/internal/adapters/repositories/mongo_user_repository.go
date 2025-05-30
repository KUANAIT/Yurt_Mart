package repositories

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-service/internal/core/domain"
)

type MongoUserRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoUserRepository(client *mongo.Client, dbName string) *MongoUserRepository {
	return &MongoUserRepository{
		client:     client,
		database:   dbName,
		collection: "users",
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *domain.User) error {
	collection := r.client.Database(r.database).Collection(r.collection)
	fmt.Printf("Inserting user into database '%s', collection '%s'\n", r.database, r.collection)
	fmt.Printf("User data to insert: %+v\n", user)

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	fmt.Printf("Successfully inserted user with ID: %v\n", result.InsertedID)

	var insertedUser domain.User
	err = collection.FindOne(ctx, bson.M{"_id": user.ID}).Decode(&insertedUser)
	if err != nil {
		fmt.Printf("Warning: Could not verify insertion: %v\n", err)
	} else {
		fmt.Printf("Verified insertion - Found user in database: %+v\n", insertedUser)
	}

	return nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"_id": id}
	fmt.Printf("Finding user by ID '%s' in database '%s'\n", id, r.database)
	fmt.Printf("Using filter: %+v\n", filter)

	var user domain.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Printf("No user found with ID '%s'\n", id)
			return nil, errors.New("user not found")
		}
		fmt.Printf("Error finding user: %v\n", err)
		return nil, err
	}

	fmt.Printf("Found user: %+v\n", user)
	return &user, nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"email": email}
	fmt.Printf("Finding user by email '%s' in database '%s'\n", email, r.database)

	var user domain.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	fmt.Printf("Found user: %+v\n", user)
	return &user, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *domain.User) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	fmt.Printf("Starting update operation for user ID: %s\n", user.ID)
	fmt.Printf("Update request data: %+v\n", user)

	existingUser, err := r.FindByID(ctx, user.ID)
	if err != nil {
		fmt.Printf("Error finding user to update: %v\n", err)
		return fmt.Errorf("failed to find user: %v", err)
	}

	fmt.Printf("Found existing user to update: %+v\n", existingUser)

	if existingUser.Email != user.Email {
		fmt.Printf("Email is being changed from '%s' to '%s'\n", existingUser.Email, user.Email)
		var existingUserWithEmail domain.User
		err := collection.FindOne(ctx, bson.M{
			"email": user.Email,
			"_id":   bson.M{"$ne": user.ID},
		}).Decode(&existingUserWithEmail)

		if err == nil {
			fmt.Printf("Email '%s' is already in use by another user\n", user.Email)
			return errors.New("email already in use")
		} else if err != mongo.ErrNoDocuments {
			fmt.Printf("Error checking email uniqueness: %v\n", err)
			return fmt.Errorf("failed to check email uniqueness: %v", err)
		}
		fmt.Printf("Email '%s' is available for use\n", user.Email)
	}

	update := bson.M{
		"$set": bson.M{
			"email": user.Email,
			"name":  user.Name,
		},
	}

	filter := bson.M{"_id": existingUser.ID}

	fmt.Printf("Attempting to update user with filter: %+v and update: %+v\n", filter, update)

	opts := options.Update().SetUpsert(false)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return fmt.Errorf("failed to update user: %v", err)
	}

	fmt.Printf("Update result - Matched: %d, Modified: %d, Upserted: %d\n",
		result.MatchedCount, result.ModifiedCount, result.UpsertedCount)

	if result.MatchedCount == 0 {
		fmt.Printf("No user found to update with ID: %s\n", existingUser.ID)
		return errors.New("user not found")
	}

	if result.ModifiedCount == 0 {
		fmt.Printf("Warning: Update operation matched but did not modify any documents\n")
	}

	fmt.Printf("Successfully updated user. Modified count: %d\n", result.ModifiedCount)

	updatedUser, err := r.FindByID(ctx, existingUser.ID)
	if err != nil {
		fmt.Printf("Warning: Could not verify update: %v\n", err)
	} else {
		fmt.Printf("Verified update - Updated user in database: %+v\n", updatedUser)
		if updatedUser.Email != user.Email || updatedUser.Name != user.Name {
			fmt.Printf("Warning: Update verification failed - Data mismatch!\n")
			fmt.Printf("Expected: Email=%s, Name=%s\n", user.Email, user.Name)
			fmt.Printf("Actual: Email=%s, Name=%s\n", updatedUser.Email, updatedUser.Name)
		}
	}

	return nil
}

func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	collection := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"_id": id}
	fmt.Printf("Deleting user with ID '%s' from database '%s'\n", id, r.database)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %d documents\n", result.DeletedCount)
	return nil
}

func (r *MongoUserRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})
	return err
}

func (r *MongoUserRepository) List(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	// Calculate skip value for pagination
	skip := (page - 1) * pageSize

	// Get total count
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Find users with pagination
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
