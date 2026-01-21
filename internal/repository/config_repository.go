package repository

import (
	"encoding/json"
	"os"
	"path/filepath"

	"ctchen222/chronoflow/internal/domain"
)

// ConfigRepository defines the interface for configuration data access
type ConfigRepository interface {
	// Load loads config from persistent storage
	Load() (domain.Config, error)

	// Save saves config to persistent storage
	Save(config domain.Config) error
}

// JSONConfigRepository implements ConfigRepository using JSON file storage
type JSONConfigRepository struct {
	filePath string
}

// NewJSONConfigRepository creates a new JSON-based config repository
func NewJSONConfigRepository(filePath string) *JSONConfigRepository {
	return &JSONConfigRepository{
		filePath: filePath,
	}
}

func (r *JSONConfigRepository) Load() (domain.Config, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config
			config := domain.DefaultConfig()
			if saveErr := r.Save(config); saveErr != nil {
				return config, saveErr
			}
			return config, nil
		}
		return domain.Config{}, err
	}
	if len(data) == 0 {
		return domain.DefaultConfig(), nil
	}

	var config domain.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return domain.Config{}, err
	}

	// Ensure defaults for any missing fields
	if config.Timeline.DayStart == "" {
		config.Timeline.DayStart = "08:00"
	}
	if config.Timeline.DayEnd == "" {
		config.Timeline.DayEnd = "18:00"
	}
	if config.Timeline.SlotMinutes == 0 {
		config.Timeline.SlotMinutes = 30
	}

	return config, nil
}

func (r *JSONConfigRepository) Save(config domain.Config) error {
	// Ensure directory exists
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}
