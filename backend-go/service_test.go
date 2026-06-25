package main

import (
	"testing"
	"time"
)

func TestBuildHandoffBrief(t *testing.T) {
	now := time.Now()

	t.Run("counts actions and open questions", func(t *testing.T) {
		inc := Incident{
			Severity:  SEV1,
			Status:    TRIGGERED,
			Service:   "api",
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: ACTION, Text: "restarted"},
				{Author: "anh", Type: OPEN_QUESTION, Text: "why?"},
				{Author: "anh", Type: OBSERVATION, Text: "cpu high"},
				{Author: "anh", Type: ACTION, Text: "scaled up"},
			},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.TakenActions != 2 {
			t.Errorf("TakenActions expected 2, got %d", brief.TakenActions)
		}
		if brief.OpenQuestion != 1 {
			t.Errorf("OpenQuestion expected 1, got %d", brief.OpenQuestion)
		}
		if brief.TotalEntry != 4 {
			t.Errorf("TotalEntry expected 4, got %d", brief.TotalEntry)
		}
	})

	t.Run("handoff count tracks author changes", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: OBSERVATION, Text: "a"},
				{Author: "anh", Type: OBSERVATION, Text: "b"},
				{Author: "bernd", Type: OBSERVATION, Text: "c"},
				{Author: "anh", Type: OBSERVATION, Text: "d"},
			},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.HandoffCount != 2 {
			t.Errorf("HandoffCount expected 2, got %d", brief.HandoffCount)
		}
	})

	t.Run("single author zero handoffs", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: OBSERVATION, Text: "a"},
				{Author: "anh", Type: OBSERVATION, Text: "b"},
			},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.HandoffCount != 0 {
			t.Errorf("HandoffCount expected 0, got %d", brief.HandoffCount)
		}
	})

	t.Run("empty entries", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.TotalEntry != 0 {
			t.Errorf("TotalEntry expected 0, got %d", brief.TotalEntry)
		}
		if brief.HandoffCount != 0 {
			t.Errorf("HandoffCount expected 0, got %d", brief.HandoffCount)
		}
	})

	t.Run("maps incident fields correctly", func(t *testing.T) {
		inc := Incident{
			Severity:  SEV2,
			Status:    INVESTIGATING,
			Service:   "payments",
			CreatedAt: now.Add(-30 * time.Minute),
			Entries:   []TimelineEntry{},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.Severity != SEV2 {
			t.Errorf("Severity expected %s, got %s", SEV2, brief.Severity)
		}
		if brief.Status != INVESTIGATING {
			t.Errorf("Status expected %s, got %s", INVESTIGATING, brief.Status)
		}
		if brief.Service != "payments" {
			t.Errorf("Service expected payments, got %s", brief.Service)
		}
		if brief.ElapsedMinute < 29 || brief.ElapsedMinute > 31 {
			t.Errorf("ElapsedMinute expected ~30, got %d", brief.ElapsedMinute)
		}
		if !brief.CreatedAt.Equal(inc.CreatedAt) {
			t.Errorf("CreatedAt mismatch")
		}
	})

	t.Run("nil flagStore skips detailed brief", func(t *testing.T) {
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, nil, "")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
		if brief.OpenQuestionList != nil {
			t.Error("expected nil OpenQuestionList")
		}
	})

	t.Run("detailed brief when flag enabled", func(t *testing.T) {
		fs := CreateFlagStore()
		fs.Create(FeatureFlag{
			Name:     "detailed_handoff_brief",
			Enabled:  true,
			Rollout:  100,
			Variants: []string{"detailed"},
		})
		inc := Incident{
			CreatedAt: now,
			Entries: []TimelineEntry{
				{Author: "anh", Type: ACTION, Text: "restarted"},
				{Author: "anh", Type: OPEN_QUESTION, Text: "why?"},
			},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList == nil {
			t.Fatal("expected non-nil TakenActionsList")
		}
		if len(*brief.TakenActionsList) != 1 {
			t.Errorf("expected 1 action, got %d", len(*brief.TakenActionsList))
		}
		if brief.OpenQuestionList == nil {
			t.Fatal("expected non-nil OpenQuestionList")
		}
		if len(*brief.OpenQuestionList) != 1 {
			t.Errorf("expected 1 question, got %d", len(*brief.OpenQuestionList))
		}
	})

	t.Run("no detailed brief when flag disabled", func(t *testing.T) {
		fs := CreateFlagStore()
		fs.Create(FeatureFlag{
			Name:     "detailed_handoff_brief",
			Enabled:  false,
			Rollout:  100,
			Variants: []string{"detailed"},
		})
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
	})

	t.Run("no detailed brief when flag not found", func(t *testing.T) {
		fs := CreateFlagStore()
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
	})

	t.Run("no detailed brief when variant is not detailed", func(t *testing.T) {
		fs := CreateFlagStore()
		fs.Create(FeatureFlag{
			Name:     "detailed_handoff_brief",
			Enabled:  true,
			Rollout:  100,
			Variants: []string{"control"},
		})
		inc := Incident{
			CreatedAt: now,
			Entries:   []TimelineEntry{{Author: "anh", Type: ACTION, Text: "a"}},
		}
		brief := buildHandoffBrief(inc, &fs, "user1")

		if brief.TakenActionsList != nil {
			t.Error("expected nil TakenActionsList")
		}
	})
}
