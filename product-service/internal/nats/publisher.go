package nats

import (
	"encoding/json"
	"log"
	"product-service/internal/domain"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(url string) *Publisher {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("❌ Failed to connect to NATS: %v", err)
	}
	log.Println("✅ Connected to NATS at", url)
	return &Publisher{conn: nc}
}

func (p *Publisher) PublishProductCreated(product *domain.Product) {
	data, err := json.Marshal(product)
	if err != nil {
		log.Println("❌ Failed to marshal product:", err)
		return
	}
	_ = p.conn.Publish("product.created", data)
	log.Println("📤 Published product.created for ID:", product.ID)
}

func (p *Publisher) PublishProductDeleted(id string) {
	payload := map[string]string{"id": id}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("❌ Failed to marshal deletion payload:", err)
		return
	}
	_ = p.conn.Publish("product.deleted", data)
	log.Println("📤 Published product.deleted for ID:", id)
}
