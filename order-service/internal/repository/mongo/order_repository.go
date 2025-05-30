package mongo

import (
	"context"
	"time"

	"github.com/hsibAD/order-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		collection: db.Collection("orders"),
	}
}

type OrderModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          string            `bson:"user_id"`
	Items           []OrderItemModel   `bson:"items"`
	TotalPrice      float64           `bson:"total_price"`
	Currency        string            `bson:"currency"`
	Status          string            `bson:"status"`
	DeliveryAddress *AddressModel     `bson:"delivery_address,omitempty"`
	DeliveryTime    time.Time         `bson:"delivery_time"`
	CreatedAt       time.Time         `bson:"created_at"`
	UpdatedAt       time.Time         `bson:"updated_at"`
}

type OrderItemModel struct {
	ProductID   string  `bson:"product_id"`
	ProductName string  `bson:"product_name"`
	Quantity    int32   `bson:"quantity"`
	UnitPrice   float64 `bson:"unit_price"`
	TotalPrice  float64 `bson:"total_price"`
}

type AddressModel struct {
	ID            string `bson:"id,omitempty"`
	UserID        string `bson:"user_id"`
	StreetAddress string `bson:"street_address"`
	City          string `bson:"city"`
	State         string `bson:"state"`
	Country       string `bson:"country"`
	PostalCode    string `bson:"postal_code"`
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	model := &OrderModel{
		UserID:       order.UserID,
		TotalPrice:   order.TotalPrice,
		Currency:     order.Currency,
		Status:       string(order.Status),
		DeliveryTime: order.DeliveryTime,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Convert items
	model.Items = make([]OrderItemModel, len(order.Items))
	for i, item := range order.Items {
		model.Items[i] = OrderItemModel{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		}
	}

	// Convert delivery address if present
	if order.DeliveryAddress != nil {
		model.DeliveryAddress = &AddressModel{
			ID:            order.DeliveryAddress.ID,
			UserID:        order.DeliveryAddress.UserID,
			StreetAddress: order.DeliveryAddress.StreetAddress,
			City:          order.DeliveryAddress.City,
			State:         order.DeliveryAddress.State,
			Country:       order.DeliveryAddress.Country,
			PostalCode:    order.DeliveryAddress.PostalCode,
		}
	}

	result, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	// Update order ID
	order.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *OrderRepository) Get(ctx context.Context, id string) (*domain.Order, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var model OrderModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	// Convert items
	items := make([]domain.OrderItem, len(model.Items))
	for i, item := range model.Items {
		items[i] = domain.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		}
	}

	// Convert delivery address if present
	var address *domain.DeliveryAddress
	if model.DeliveryAddress != nil {
		address = &domain.DeliveryAddress{
			ID:            model.DeliveryAddress.ID,
			UserID:        model.DeliveryAddress.UserID,
			StreetAddress: model.DeliveryAddress.StreetAddress,
			City:          model.DeliveryAddress.City,
			State:         model.DeliveryAddress.State,
			Country:       model.DeliveryAddress.Country,
			PostalCode:    model.DeliveryAddress.PostalCode,
		}
	}

	return &domain.Order{
		ID:              model.ID.Hex(),
		UserID:          model.UserID,
		Items:           items,
		TotalPrice:      model.TotalPrice,
		Currency:        model.Currency,
		Status:          domain.OrderStatus(model.Status),
		DeliveryAddress: address,
		DeliveryTime:    model.DeliveryTime,
	}, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	objectID, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     string(order.Status),
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
} 