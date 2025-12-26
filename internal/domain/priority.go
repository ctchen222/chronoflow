package domain

// Priority represents the priority level of a todo
type Priority int

const (
	PriorityNone   Priority = 0
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

// String returns the string representation of the priority
func (p Priority) String() string {
	switch p {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "None"
	}
}

// Icon returns the icon representation of the priority
func (p Priority) Icon() string {
	switch p {
	case PriorityHigh:
		return "!!!"
	case PriorityMedium:
		return "!!"
	case PriorityLow:
		return "!"
	default:
		return ""
	}
}

// IsValid checks if the priority value is valid
func (p Priority) IsValid() bool {
	return p >= PriorityNone && p <= PriorityHigh
}
