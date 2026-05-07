package main

import "time"

func buildHandoffBrief(inc Incident) HandoffBrief {
	actions := []TimelineEntry{}
	openQuestions := []TimelineEntry{}

	for _, entry := range inc.Entries {
		switch entry.Type {
		case ACTION:
			actions = append(actions, entry)
		case OPEN_QUESTION:
			openQuestions = append(openQuestions, entry)
		}
	}

	return HandoffBrief{
		Severity:      inc.Severity,
		Status:        inc.Status,
		Service:       inc.Service,
		ElapsedMinute: int(time.Since(inc.CreatedAt).Minutes()),
		TotalEntry:    len(inc.Entries),
		TakenActions:  actions,
		OpenQuestion:  openQuestions,
		CreatedAt:     inc.CreatedAt,
	}
}
