package models

import (
	"fmt"
	"time"
)

type Theme struct {
	Name      string      `json:"name"`
	Status    ThemeStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Color     string      `json:"color"`
	Icon      string      `json:"icon"`
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
