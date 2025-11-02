package core

import (
	"context"
	"fmt"
	"log"
	"multitrack-bot/internal/adapters"
	"multitrack-bot/internal/domain"
	"time"
)

type TrackingService struct {
	adapterManager *adapters.AdapterManager
	normalizer     *NormalizationEngine
}

func NewTrackingService(adapterManager *adapters.AdapterManager) *TrackingService {
	return &TrackingService{
		adapterManager: adapterManager,
		normalizer:     NewNormalizationEngine(),
	}
}

func (s *TrackingService) Track(ctx context.Context, trackingNumber, courier string) (*domain.TrackingResult, error) {

	log.Printf("[Track] courier=%s number=%s", courier, trackingNumber)

	// get adapter for this courier
	adapter, err := s.adapterManager.GetAdapter(courier)
	if err != nil {
		log.Printf("[ERROR] get adapter: %v", err)
		return nil, fmt.Errorf("failed to get adapter: %w", err)
	}

	// tracking with adapter
	rawResult, err := adapter.Track(ctx, trackingNumber)
	if err != nil {
		log.Printf("[ERROR] adapter.Track: %v", err)
		return nil, fmt.Errorf("tracking failed: %w", err)
	}

	// normalize data
	normalizedResult := s.normalizer.Normalize(rawResult)
	normalizedResult.LastUpdated = time.Now()

	return normalizedResult, nil
}

// NormalizationEngine transforms raw data to universal format
type NormalizationEngine struct{}

func NewNormalizationEngine() *NormalizationEngine {
	return &NormalizationEngine{}
}

func (n *NormalizationEngine) Normalize(rawResult *domain.RawTrackingResult) *domain.TrackingResult {
	result := &domain.TrackingResult{
		Number:      "tracking_number", // todo
		Courier:     rawResult.Courier,
		Status:      "Processing",
		Description: "Package is processing",
		Checkpoints: []domain.Checkpoint{
			{
				Date:        time.Now(),
				Location:    "Saint-Petersburg",
				Status:      "Accepted",
				Description: "Package was Delivered and Received",
			},
		},
	}

	return result
}
