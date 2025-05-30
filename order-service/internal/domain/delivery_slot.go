package domain

import "time"

type DeliverySlot struct {
	StartTime time.Time
	EndTime   time.Time
	Available bool
} 