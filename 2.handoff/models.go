package main

import (
	"errors"
	"strings"
	"time"
)

type Incident struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Service   string          `json:"service"`
	Severity  string          `json:"severity"` // SEV1, SEV2, SEV3
	Status    string          `json:"status"`   // triggered, acknowledged, investigating, mitigated, resolved
	OpenedBy  string          `json:"opened_by"`
	OnCall    string          `json:"on_call"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Entries   []TimelineEntry `json:"entries"`
}

type TimelineEntry struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"time"`
	Author string    `json:"author"`
	Type   string    `json:"type"` // observation, action, discovery, open_question, state_change
	Text   string    `json:"text"`
}

// type IncidentFilter struct {
// }

// type IncidentUpdate struct {
// }

type CreateIncidentRequest struct {
	Title    string `json:"title"`
	Service  string `json:"service"`
	Severity string `json:"severity"` // SEV1, SEV2, SEV3
	OpenedBy string `json:"opened_by"`
}

type GetIncidentRequest struct {
	IncidentID string `json:"incident_id"`
}

func (c *CreateIncidentRequest) Validate() error {
	if strings.Trim(c.Title, " ") == "" {
		return errors.New("Request doesn't contain title")
	}
	if strings.Trim(c.Service, " ") == "" {
		return errors.New("Request doesn't contain service")
	}
	if strings.Trim(c.Severity, " ") == "" {
		return errors.New("Request doesn't contain severity")
	}
	if strings.Trim(c.OpenedBy, " ") == "" {
		return errors.New("Request doesn't contain opened_by")
	}
	return nil
}
