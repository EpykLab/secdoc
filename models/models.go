package models

import (
	"go/token"
)

type SecurityControl struct {
	ControlID    string         `json:"control_id"`
	Description  string         `json:"description"`
	References   string         `json:"references"`
	Verification string         `json:"verification"`
	FilePath     string         `json:"file_path"`
	Position     token.Position `json:"position"`
}

type Requirement struct {
	RequirementID string         `json:"requirement_id"`
	Description   string         `json:"description"`
	Verification  string         `json:"verification"`
	Stakeholder   string         `json:"stakeholder"`
	FilePath      string         `json:"file_path"`
	Position      token.Position `json:"position"`
}

type Report struct {
	SecurityControls []SecurityControl `json:"security_controls"`
	Requirements     []Requirement     `json:"requirements"`
}
