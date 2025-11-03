package domain

import (
	"context"
)

type CourierAdapter interface {
	Name() string
	Validate(trackingNumber string) bool
	Track(ctx context.Context, trackingNumber string) (*RawTrackingResult, error)
}

type Normalizer interface {
	CanNormalize(courierName string) bool
	Normalize(raw *RawTrackingResult) *TrackingResult
}
