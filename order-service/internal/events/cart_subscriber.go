package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	CartInfoSubject = "cart.info"
)

type CartInfo struct {
	CartID     string  `json:"cart_id"`
	TotalPrice float64 `json:"total_price"`
	Currency   string  `json:"currency"`
}

type CartSubscriber struct {
	nc *nats.Conn
	js nats.JetStreamContext
	cartInfoCache map[string]*CartInfo
}

func NewCartSubscriber(url string) (*CartSubscriber, error) {
	// Try to connect with retries
	var nc *nats.Conn
	var err error
	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to NATS, attempt %d: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Create the stream if it doesn't exist
	stream := &nats.StreamConfig{
		Name:     "CARTS",
		Subjects: []string{"cart.*"},
	}

	if _, err := js.AddStream(stream); err != nil {
		if err != nats.ErrStreamNameAlreadyInUse {
			return nil, err
		}
	}

	sub := &CartSubscriber{
		nc: nc,
		js: js,
		cartInfoCache: make(map[string]*CartInfo),
	}

	// Subscribe to cart info updates
	if _, err := js.Subscribe(CartInfoSubject, sub.handleCartInfo); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *CartSubscriber) handleCartInfo(msg *nats.Msg) {
	var cartInfo CartInfo
	if err := json.Unmarshal(msg.Data, &cartInfo); err != nil {
		log.Printf("Error unmarshaling cart info: %v", err)
		return
	}

	s.cartInfoCache[cartInfo.CartID] = &cartInfo
}

func (s *CartSubscriber) GetCartInfo(ctx context.Context, cartID string) (*CartInfo, error) {
	// First check the cache
	if info, ok := s.cartInfoCache[cartID]; ok {
		return info, nil
	}

	// If not in cache, request it directly
	msg, err := s.nc.Request(CartInfoSubject+".get", []byte(cartID), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var cartInfo CartInfo
	if err := json.Unmarshal(msg.Data, &cartInfo); err != nil {
		return nil, err
	}

	// Cache the result
	s.cartInfoCache[cartID] = &cartInfo
	return &cartInfo, nil
}

func (s *CartSubscriber) Close() error {
	s.nc.Close()
	return nil
} 