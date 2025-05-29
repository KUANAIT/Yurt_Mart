package model

type CartItem struct {
	UserID    string `bson:"user_id"`
	ProductID string `bson:"product_id"`
	Quantity  int32  `bson:"quantity"`
}
