package models

import (
	"fmt"
	"time"
)

type Theme struct {
	Name      string      `json:"name"`
	Status    ThemeStatus `json:"status"`
	Type      ThemeType   `json:"category"`
	CreatedAt time.Time   `json:"created_at"`
}

type ThemeType string

const (
	ThemeTypeCategory ThemeType = "category"
	ThemeTypeTheme    ThemeType = "person"
)

func ParseThemeType(s string) (ThemeType, error) {
	switch s {
	case "category":
		return ThemeTypeCategory, nil
	case "person":
		return ThemeTypeTheme, nil
	default:
		return "", fmt.Errorf("invalid theme type %s", s)
	}
}

type ThemeStatus string

const (
	ThemeStatusActive   ThemeStatus = "active"
	ThemeStatusInactive ThemeStatus = "inactive"
)

func ParseThemeStatus(s string) (ThemeStatus, error) {
	switch s {
	case "active":
		return ThemeStatusActive, nil
	case "inactive":
		return ThemeStatusInactive, nil
	default:
		return "", fmt.Errorf("invalid theme status %s", s)
	}
}
