package core

import (
	"context"
	"fmt"
	"log"
	"multitrack-bot/internal/adapters"
	"multitrack-bot/internal/core/normalizers"
	"multitrack-bot/internal/domain"
)

type TrackingService struct {
	adapterManager    *adapters.AdapterManager
	normalizerManager *normalizers.NormilazerManager
}

func NewTrackingService(adapterManager *adapters.AdapterManager) *TrackingService {
	return &TrackingService{
		adapterManager:    adapterManager,
		normalizerManager: normalizers.NewNormalizerManager(),
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
	normalizedResult, err := s.normalizerManager.Normalize(rawResult)
	if err != nil {
		log.Printf("[ERROR] normalize: %v", err)
		return nil, fmt.Errorf("normalization failed: %w", err)
	}

	return normalizedResult, nil
}
