package domain

type Product struct {
	ID          string  `bson:"_id,omitempty"`
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Category    string  `bson:"category"`
	Price       float64 `bson:"price"`
	Quantity    int32   `bson:"quantity"`
	UserID      string  `bson:"user_id"`
}
