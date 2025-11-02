package domain

import (
	"context"
)

type CourierAdapter interface {
	Name() string
	Validate(trackingNumber string) bool
	Track(ctx context.Context, trackingNumber string) (*RawTrackingResult, error)
}
