package domain

type User struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"-" bson:"password"`
}

type UserEvent struct {
	Type    string `json:"type"`
	UserID  string `json:"user_id"`
	Payload []byte `json:"payload,omitempty"`
}
