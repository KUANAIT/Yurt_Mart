package events

import (
	"context"
	"encoding/json"
	"fmt"
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
	nc            *nats.Conn
	js            nats.JetStreamContext
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
		log.Printf("[ERROR] Could not connect to NATS for CartSubscriber: %v. Using noop subscriber.", err)
		return &CartSubscriber{
			nc:            nil,
			js:            nil,
			cartInfoCache: make(map[string]*CartInfo),
		}, nil
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Printf("[ERROR] Could not get JetStream context: %v. Using noop subscriber.", err)
		return &CartSubscriber{
			nc:            nc,
			js:            nil,
			cartInfoCache: make(map[string]*CartInfo),
		}, nil
	}

	// Create the stream if it doesn't exist
	stream := &nats.StreamConfig{
		Name:     "CARTS",
		Subjects: []string{"cart.*"},
	}

	if _, err := js.AddStream(stream); err != nil {
		if err != nats.ErrStreamNameAlreadyInUse {
			log.Printf("[ERROR] Could not add NATS stream: %v. Using noop subscriber.", err)
			return &CartSubscriber{
				nc:            nc,
				js:            js,
				cartInfoCache: make(map[string]*CartInfo),
			}, nil
		}
	}

	sub := &CartSubscriber{
		nc:            nc,
		js:            js,
		cartInfoCache: make(map[string]*CartInfo),
	}

	// Subscribe to cart info updates
	if js != nil {
		if _, err := js.Subscribe(CartInfoSubject, sub.handleCartInfo); err != nil {
			log.Printf("[ERROR] Could not subscribe to cart info: %v. Using noop subscriber.", err)
			sub.nc = nil
			sub.js = nil
		}
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

	// Если нет подключения к NATS, возвращаем ошибку-заглушку
	if s.nc == nil {
		return nil, ErrNATSUnavailable()
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

func ErrNATSUnavailable() error {
	return fmt.Errorf("NATS unavailable: CartSubscriber is in noop mode")
}

func (s *CartSubscriber) Close() error {
	if s.nc != nil {
		s.nc.Close()
	}
	return nil
}
