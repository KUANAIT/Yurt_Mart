package domain

const (
	UserCreatedEvent = "user.created"
	UserUpdatedEvent = "user.updated"
	UserDeletedEvent = "user.deleted"
)

type EventHandler interface {
	HandleUserCreated(user *User) error
	HandleUserUpdated(user *User) error
	HandleUserDeleted(userID string) error
}
