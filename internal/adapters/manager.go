package adapters

import (
	"fmt"
	"multitrack-bot/internal/adapters/russianpost"
	"multitrack-bot/internal/config"
	"multitrack-bot/internal/domain"
)

type AdapterManager struct {
	adapters map[string]domain.CourierAdapter
}

func NewAdapterManager(cfg *config.Config) *AdapterManager {
	manager := &AdapterManager{
		adapters: make(map[string]domain.CourierAdapter),
	}

	if cfg.RussianPostLogin == "" || cfg.RussianPostPass == "" {
		fmt.Println("⚠️ Warning: Missing Russian Post API credentials in environment variables")
	}

	// register adapter
	manager.RegisterAdapter(russianpost.NewRussianPostAdapter(
		cfg.RussianPostLogin, cfg.RussianPostPass,
	))

	return manager
}

func (m *AdapterManager) RegisterAdapter(adapter domain.CourierAdapter) {
	m.adapters[adapter.Name()] = adapter
}

func (m *AdapterManager) GetAdapter(name string) (domain.CourierAdapter, error) {
	adapter, exists := m.adapters[name]
	if !exists {
		return nil, fmt.Errorf("adapter not found: %s", name)
	}
	return adapter, nil
}
