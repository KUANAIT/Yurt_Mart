package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	UserID     string    `bson:"user_id" json:"user_id"`
	Username   string    `bson:"username" json:"username"`
	ProductIDs []string  `bson:"product_ids" json:"product_ids"`
	Total      float64   `bson:"total" json:"total"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

func (o *Order) BeforeCreate() {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
}
