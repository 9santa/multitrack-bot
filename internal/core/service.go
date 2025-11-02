package core

import (
	"context"
	"fmt"
	"log"
	"multitrack-bot/internal/adapters"
	"multitrack-bot/internal/adapters/russianpost"
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

func (n *NormalizationEngine) Normalize(raw *domain.RawTrackingResult) *domain.TrackingResult {
	result := &domain.TrackingResult{
		Courier: raw.Courier,
	}

	switch data := raw.RawData.(type) {
	case []russianpost.HistoryRecord:
		if len(data) > 0 {
			// last operation is the last record
			last := data[len(data)-1]
			result.Number = last.Barcode

			var statusMap = map[string]string{
				"Вручение":  "Доставлено",
				"Обработка": "В пути",
				"Прием":     "Принято",
				"Присвоение идентификатора": "Создана",
			}

			humanStatus := statusMap[last.OperType]
			if humanStatus != "" {
				result.Status = humanStatus
			}

			result.Description = last.OperType

			for _, r := range data {
				t, _ := time.Parse(time.RFC3339, r.OperDate)
				result.Checkpoints = append(result.Checkpoints, domain.Checkpoint{
					Date:        t,
					Location:    r.Address,
					Status:      r.OperAttr,
					Description: r.OperType,
				})
			}
		}
	}

	result.LastUpdated = time.Now()
	return result

}
