package promexp

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type series struct {
	name    string
	t       seriesType
	help    string
	metrics map[string]*metric
}

func newSeries(name string, t seriesType, help string) *series {
	return &series{
		name:    name,
		t:       t,
		help:    help,
		metrics: make(map[string]*metric),
	}
}

func (s *series) setMetric(ls labelsT, value float64) {
	s.metrics[ls.String()] = newMetric(ls, value)
}

func (s *series) expire(ttl time.Duration) {
	for _, m := range s.metrics {
		if m.isExpired(ttl) {
			delete(s.metrics, m.labels.String())
		}
	}
}

func (s *series) typeString() string {
	return fmt.Sprintf("# TYPE %s %s", s.name, s.t)
}

func (s *series) helpString() string {
	return fmt.Sprintf("# HELP %s %s", s.name, s.help)
}

func (s *series) String() string {
	if len(s.metrics) == 0 {
		return ""
	}

	lines := make([]string, 0)

	for _, m := range s.metrics {
		lines = append(lines, fmt.Sprintf("%s%s", s.name, m.String()))
	}
	sort.Strings(lines)

	lines = append([]string{s.helpString(), s.typeString()}, lines...)

	return strings.Join(lines, "\n")
}
