package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/hsibAD/payment-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

type mongoPayment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OrderID       string             `bson:"order_id"`
	UserID        string             `bson:"user_id"`
	Amount        float64            `bson:"amount"`
	Currency      string             `bson:"currency"`
	Status        string             `bson:"status"`
	PaymentMethod string             `bson:"payment_method"`
	TransactionID string             `bson:"transaction_id,omitempty"`
	ErrorMessage  string             `bson:"error_message,omitempty"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func NewPaymentRepository(db *mongo.Database) *PaymentRepository {
	return &PaymentRepository{
		db:         db,
		collection: db.Collection("payments"),
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	mPayment := toMongoPayment(payment)
	result, err := r.collection.InsertOne(ctx, mPayment)
	if err != nil {
		return err
	}

	payment.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidPaymentID
	}

	var mPayment mongoPayment
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mPayment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrInvalidPaymentID
		}
		return nil, err
	}

	return fromMongoPayment(&mPayment), nil
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) ([]*domain.Payment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var mPayments []mongoPayment
	if err = cursor.All(ctx, &mPayments); err != nil {
		return nil, err
	}

	payments := make([]*domain.Payment, len(mPayments))
	for i, mPayment := range mPayments {
		payments[i] = fromMongoPayment(&mPayment)
	}

	return payments, nil
}

func (r *PaymentRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*domain.Payment, int, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var mPayments []mongoPayment
	if err = cursor.All(ctx, &mPayments); err != nil {
		return nil, 0, err
	}

	payments := make([]*domain.Payment, len(mPayments))
	for i, mPayment := range mPayments {
		payments[i] = fromMongoPayment(&mPayment)
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, 0, err
	}

	return payments, int(total), nil
}

func (r *PaymentRepository) GetPendingPayments(ctx context.Context, page, limit int) ([]*domain.Payment, int, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{"status": domain.PaymentStatusPending}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var mPayments []mongoPayment
	if err = cursor.All(ctx, &mPayments); err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	count, err := r.collection.CountDocuments(ctx, bson.M{"status": domain.PaymentStatusPending})
	if err != nil {
		return nil, 0, err
	}

	payments := make([]*domain.Payment, len(mPayments))
	for i, mPayment := range mPayments {
		payments[i] = fromMongoPayment(&mPayment)
	}

	return payments, int(count), nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	objectID, err := primitive.ObjectIDFromHex(payment.ID)
	if err != nil {
		return domain.ErrInvalidPaymentID
	}

	mPayment := toMongoPayment(payment)
	mPayment.ID = objectID

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, mPayment)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrInvalidPaymentID
	}

	return nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, paymentID string, status domain.PaymentStatus) error {
	objectID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return domain.ErrInvalidPaymentID
	}

	update := bson.M{
		"$set": bson.M{
			"status":     string(status),
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// Добавляем функции конвертации между domain.Payment и mongoPayment
func toMongoPayment(payment *domain.Payment) mongoPayment {
	var id primitive.ObjectID
	if payment.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(payment.ID)
		if err == nil {
			id = objectID
		}
	}
	return mongoPayment{
		ID:            id,
		OrderID:       payment.OrderID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		PaymentMethod: string(payment.PaymentMethod),
		TransactionID: payment.TransactionID,
		ErrorMessage:  payment.ErrorMessage,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
}

func fromMongoPayment(m *mongoPayment) *domain.Payment {
	return &domain.Payment{
		ID:            m.ID.Hex(),
		OrderID:       m.OrderID,
		UserID:        m.UserID,
		Amount:        m.Amount,
		Currency:      m.Currency,
		Status:        m.Status,
		PaymentMethod: m.PaymentMethod,
		TransactionID: m.TransactionID,
		ErrorMessage:  m.ErrorMessage,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
