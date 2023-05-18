package promexp

import (
	"fmt"
	"reflect"
	"time"
)

type metric struct {
	labels    labelsT
	value     float64
	timestamp time.Time
}

func newMetric(ls labelsT, value float64) *metric {
	return &metric{
		labels:    ls,
		value:     value,
		timestamp: time.Now(),
	}
}

func (m *metric) String() string {
	return fmt.Sprintf("%s %g %d", m.labels.String(), m.value, m.timestamp.UnixMilli())
}

func (m *metric) isSame(v *metric) bool {
	if v == nil || m == nil {
		return false
	}

	if reflect.DeepEqual(m.labels, v.labels) {
		return true
	}

	return false
}

func (m *metric) isExpired(ttl time.Duration) bool {
	if m.timestamp.Add(ttl).Before(time.Now()) {
		return true
	}

	return false
}
