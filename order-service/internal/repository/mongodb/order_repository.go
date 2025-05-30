package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/hsibAD/order-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

type mongoOrder struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty"`
	UserID          string              `bson:"user_id"`
	Items           []mongoOrderItem    `bson:"items"`
	TotalPrice      float64             `bson:"total_price"`
	Currency        string              `bson:"currency"`
	Status          string              `bson:"status"`
	DeliveryAddress *mongoDeliveryAddress `bson:"delivery_address"`
	DeliveryTime    time.Time           `bson:"delivery_time"`
	CreatedAt       time.Time           `bson:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at"`
}

type mongoOrderItem struct {
	ProductID   string  `bson:"product_id"`
	ProductName string  `bson:"product_name"`
	Quantity    int32   `bson:"quantity"`
	UnitPrice   float64 `bson:"unit_price"`
	TotalPrice  float64 `bson:"total_price"`
}

type mongoDeliveryAddress struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string            `bson:"user_id"`
	FullName      string            `bson:"full_name"`
	StreetAddress string            `bson:"street_address"`
	Apartment     string            `bson:"apartment"`
	City          string            `bson:"city"`
	State         string            `bson:"state"`
	PostalCode    string            `bson:"postal_code"`
	Country       string            `bson:"country"`
	Phone         string            `bson:"phone"`
	IsDefault     bool              `bson:"is_default"`
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		db:         db,
		collection: db.Collection("orders"),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	mOrder := toMongoOrder(order)
	result, err := r.collection.InsertOne(ctx, mOrder)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidOrderID
	}

	var mOrder mongoOrder
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrInvalidOrderID
		}
		return nil, err
	}

	return fromMongoOrder(&mOrder), nil
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*domain.Order, int, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var mOrders []mongoOrder
	if err = cursor.All(ctx, &mOrders); err != nil {
		return nil, 0, err
	}

	orders := make([]*domain.Order, len(mOrders))
	for i, mOrder := range mOrders {
		orders[i] = fromMongoOrder(&mOrder)
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	objectID, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return domain.ErrInvalidOrderID
	}

	mOrder := toMongoOrder(order)
	mOrder.ID = objectID

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, mOrder)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrInvalidOrderID
	}

	return nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return domain.ErrInvalidOrderID
	}

	update := bson.M{
		"$set": bson.M{
			"status":     string(status),
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrInvalidOrderID
	}

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidOrderID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrInvalidOrderID
	}

	return nil
}

func toMongoOrder(order *domain.Order) *mongoOrder {
	items := make([]mongoOrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = mongoOrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		}
	}

	var deliveryAddress *mongoDeliveryAddress
	if order.DeliveryAddress != nil {
		deliveryAddress = &mongoDeliveryAddress{
			UserID:        order.DeliveryAddress.UserID,
			FullName:      order.DeliveryAddress.FullName,
			StreetAddress: order.DeliveryAddress.StreetAddress,
			Apartment:     order.DeliveryAddress.Apartment,
			City:          order.DeliveryAddress.City,
			State:         order.DeliveryAddress.State,
			PostalCode:    order.DeliveryAddress.PostalCode,
			Country:       order.DeliveryAddress.Country,
			Phone:         order.DeliveryAddress.Phone,
			IsDefault:     order.DeliveryAddress.IsDefault,
		}
		if order.DeliveryAddress.ID != "" {
			if objectID, err := primitive.ObjectIDFromHex(order.DeliveryAddress.ID); err == nil {
				deliveryAddress.ID = objectID
			}
		}
	}

	mOrder := &mongoOrder{
		UserID:          order.UserID,
		Items:           items,
		TotalPrice:      order.TotalPrice,
		Currency:        order.Currency,
		Status:          string(order.Status),
		DeliveryAddress: deliveryAddress,
		DeliveryTime:    order.DeliveryTime,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}

	if order.ID != "" {
		if objectID, err := primitive.ObjectIDFromHex(order.ID); err == nil {
			mOrder.ID = objectID
		}
	}

	return mOrder
}

func fromMongoOrder(mOrder *mongoOrder) *domain.Order {
	items := make([]domain.OrderItem, len(mOrder.Items))
	for i, item := range mOrder.Items {
		items[i] = domain.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		}
	}

	var deliveryAddress *domain.DeliveryAddress
	if mOrder.DeliveryAddress != nil {
		deliveryAddress = &domain.DeliveryAddress{
			ID:            mOrder.DeliveryAddress.ID.Hex(),
			UserID:        mOrder.DeliveryAddress.UserID,
			FullName:      mOrder.DeliveryAddress.FullName,
			StreetAddress: mOrder.DeliveryAddress.StreetAddress,
			Apartment:     mOrder.DeliveryAddress.Apartment,
			City:          mOrder.DeliveryAddress.City,
			State:         mOrder.DeliveryAddress.State,
			PostalCode:    mOrder.DeliveryAddress.PostalCode,
			Country:       mOrder.DeliveryAddress.Country,
			Phone:         mOrder.DeliveryAddress.Phone,
			IsDefault:     mOrder.DeliveryAddress.IsDefault,
		}
	}

	return &domain.Order{
		ID:              mOrder.ID.Hex(),
		UserID:          mOrder.UserID,
		Items:           items,
		TotalPrice:      mOrder.TotalPrice,
		Currency:        mOrder.Currency,
		Status:          domain.OrderStatus(mOrder.Status),
		DeliveryAddress: deliveryAddress,
		DeliveryTime:    mOrder.DeliveryTime,
		CreatedAt:       mOrder.CreatedAt,
		UpdatedAt:       mOrder.UpdatedAt,
	}
} 