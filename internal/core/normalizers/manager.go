package normalizers

import (
	"fmt"
	"multitrack-bot/internal/domain"
)

type NormilazerManager struct {
	normalizers []domain.Normalizer
}

func NewNormalizerManager() *NormilazerManager {
	manager := &NormilazerManager{
		normalizers: []domain.Normalizer{},
	}

	// register all normalizers
	manager.RegisterNormalizer(NewRussianPostNormalizer())

	return manager
}

func (m *NormilazerManager) RegisterNormalizer(normalizer domain.Normalizer) {
	m.normalizers = append(m.normalizers, normalizer)
}

func (m *NormilazerManager) Normalize(raw *domain.RawTrackingResult) (*domain.TrackingResult, error) {
	for _, normalizer := range m.normalizers {
		if normalizer.CanNormalize(raw.Courier) {
			return normalizer.Normalize(raw), nil
		}
	}

	return nil, fmt.Errorf("no normalizer found for this courier: %s", raw.Courier)
}
