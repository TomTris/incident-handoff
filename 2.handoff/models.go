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

func (c *TimelineEntry) Validate() error {
	if strings.TrimSpace(c.Author) == "" {
		return ErrNoAuthor
	}
	if validEntryTypes[strings.TrimSpace(c.Type)] == false {
		return ErrBadEntryType
	}
	if strings.TrimSpace(c.Text) == "" {
		return ErrNoText
	}
	return nil
}

type IncidentFilter struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func (f *IncidentFilter) Validate() error {
	if f.Status != "" && !IncidentStatus[strings.TrimSpace(f.Status)] {
		return errors.New("Invalid Incident status")
	}
	return nil
}

type IncidentUpdate struct {
	Status   *string `json:"status"`
	Severity *string `json:"severity"`
	OnCall   *string `json:"on_call"`
}

func (f *IncidentUpdate) Validate() error {
	if f.Status != nil && IncidentStatus[strings.TrimSpace(*f.Status)] == false {
		return ErrBadStatus
	}
	if f.Severity != nil && IncidentSeverity[strings.TrimSpace(*f.Severity)] == false {
		return ErrInvalidSeverity
	}
	if f.OnCall != nil && strings.TrimSpace(*f.OnCall) == "" {
		return ErrOnCall
	}
	return nil
}

type CreateIncidentRequest struct {
	Title    string  `json:"title"`
	Service  string  `json:"service"`
	Severity string  `json:"severity"` // SEV1, SEV2, SEV3
	OpenedBy string  `json:"opened_by"`
	OnCall   *string `json:"on_call"`
}

func (c *CreateIncidentRequest) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return ErrNoTitle
	}
	if strings.TrimSpace(c.Service) == "" {
		return ErrNoService
	}
	if IncidentSeverity[strings.TrimSpace(c.Severity)] == false {
		return ErrInvalidSeverity
	}
	if strings.TrimSpace(c.OpenedBy) == "" {
		return ErrOpenedBy
	}
	if c.OnCall != nil && strings.TrimSpace(*c.OnCall) == "" {
		return ErrOnCall
	}
	return nil
}

type HandoffBrief struct {
	Severity      string          `json:"severity"`
	Status        string          `json:"status"`
	Service       string          `json:"service"`
	TotalEntry    int             `json:"total_entry"`
	ElapsedMinute int             `json:"elapsed_minute"`
	TakenActions  []TimelineEntry `json:"taken_actions"`
	OpenQuestion  []TimelineEntry `json:"open_question"`
	CreatedAt     time.Time       `json:"created_at"`
}
