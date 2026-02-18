package cafelisting

import "strings"

const (
	VisitStatusToVisit = "to_visit"
	VisitStatusVisited = "visited"
)

func normalizeVisitStatus(input string) (string, error) {
	status := strings.TrimSpace(strings.ToLower(input))
	if status == "" {
		return VisitStatusToVisit, nil
	}
	if status != VisitStatusToVisit && status != VisitStatusVisited {
		return "", ErrInvalidVisitStatus
	}
	return status, nil
}
