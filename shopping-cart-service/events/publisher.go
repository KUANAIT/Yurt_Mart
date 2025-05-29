package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

var natsConn *nats.Conn

func InitNATS(url string) {
	var err error
	natsConn, err = nats.Connect(url, nats.Timeout(5*time.Second))
	if err != nil {
		log.Fatalf("[NATS] Failed to connect: %v", err)
	}
	log.Println("[NATS] Connected")
}

type CartItemAddedEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

func PublishCartItemAdded(event CartItemAddedEvent) {
	if natsConn == nil {
		log.Println("[NATS] Not initialized")
		return
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[NATS] Marshal error: %v", err)
		return
	}

	err = natsConn.Publish("cart.item.added", data)
	if err != nil {
		log.Printf("[NATS] Publish error: %v", err)
	} else {
		log.Printf("[NATS] Published: user_id=%s, product_id=%s", event.UserID, event.ProductID)
	}
}

func SubscribeToCartEvents() {
	if natsConn == nil {
		log.Println("[NATS] Not initialized")
		return
	}

	_, err := natsConn.Subscribe("cart.item.added", func(m *nats.Msg) {
		var event CartItemAddedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Println("Failed to unmarshal event:", err)
			return
		}

		// Используем user_id как email
		toEmail := event.UserID

		body := "You add product " + event.ProductID +
			" (x" + fmt.Sprint(event.Quantity) + ") to card."

		SendEmailSMTP(toEmail, "🛒 Cart updated", body)

	})
	if err != nil {
		log.Println("Failed to subscribe to cart.item.added:", err)
	}
}
