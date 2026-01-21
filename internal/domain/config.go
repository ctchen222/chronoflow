package domain

// TimelineConfig holds user preferences for the timeline display
type TimelineConfig struct {
	DayStart    string `json:"day_start"`    // "HH:MM" format, default "08:00"
	DayEnd      string `json:"day_end"`      // "HH:MM" format, default "18:00"
	SlotMinutes int    `json:"slot_minutes"` // display granularity, default 30
	MoveMinutes int    `json:"move_minutes"` // movement granularity for Shift+J/K, default same as SlotMinutes
}

// DefaultTimelineConfig returns the default timeline configuration
func DefaultTimelineConfig() TimelineConfig {
	return TimelineConfig{
		DayStart:    "08:00",
		DayEnd:      "18:00",
		SlotMinutes: 30,
		MoveMinutes: 30, // Default to same as SlotMinutes
	}
}

// Config holds all user configuration
type Config struct {
	Timeline TimelineConfig `json:"timeline"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Timeline: DefaultTimelineConfig(),
	}
}
