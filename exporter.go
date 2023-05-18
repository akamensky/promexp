package prommetric

import (
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	// DefaultExpiration is a time that it takes metric without any updates to expire and be removed from exporter
	DefaultExpiration = time.Duration(5) * time.Minute
	// DefaultInterval is time interval between expiration evaluations
	DefaultInterval = time.Duration(1) * time.Minute
)

// Exporter is an abstraction for exporter
// it does not actually implement any of HTTP functions,
// but it does group the metrics, handles their expiration
// and renders them into a text. Should be instantiated using
// NewExporter() method.
type Exporter struct {
	mutex  *sync.Mutex
	ttl    time.Duration
	series []*series
}

// NewExporter creates an instance of Exporter
func NewExporter() *Exporter {
	e := &Exporter{
		mutex:  &sync.Mutex{},
		ttl:    DefaultExpiration,
		series: make([]*series, 0),
	}

	go func() {
		for true {
			time.Sleep(DefaultInterval)
			e.expire()
		}
	}()

	return e
}

func (e *Exporter) expire() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, s := range e.series {
		s.expire(e.ttl)
	}
}

func (e *Exporter) setMetric(name string, labels map[string]string, value float64, t seriesType, help string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var metricSeries *series
	for _, s := range e.series {
		if name == s.name {
			metricSeries = s
		}
	}

	if metricSeries == nil {
		metricSeries = newSeries(name, t, help)
		e.series = append(e.series, metricSeries)
		sort.Slice(e.series, func(i, j int) bool {
			return e.series[i].name < e.series[j].name
		})
	}

	metricSeries.setMetric(newLabels(labels), value)
}

// SetGauge creates new or updates existing gauge metric,
// the timestamp is considered that when the metric
// is being set using this method
func (e *Exporter) SetGauge(name string, value float64, help string, labels map[string]string) {
	e.setMetric(name, labels, value, gauge, help)
}

// SetCounter creates new or updates existing counter metric,
// the timestamp is considered that when the metric
// is being set using this method
func (e *Exporter) SetCounter(name string, value float64, help string, labels map[string]string) {
	e.setMetric(name, labels, value, counter, help)
}

// SetHistogram creates new or updates existing histogram metric,
// the timestamp is considered that when the metric
// is being set using this method
func (e *Exporter) SetHistogram(name string, value float64, help string, labels map[string]string) {
	e.setMetric(name, labels, value, histogram, help)
}

// SetSummary creates new or updates existing summary metric,
// the timestamp is considered that when the metric
// is being set using this method
func (e *Exporter) SetSummary(name string, value float64, help string, labels map[string]string) {
	e.setMetric(name, labels, value, summary, help)
}

// Render renders all active metrics into a text output
// that can be then plugged into HTTP handler.
func (e *Exporter) Render() string {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	values := make([]string, 0)
	for _, s := range e.series {
		value := s.String()
		if len(value) > 0 {
			values = append(values, s.String())
		}
	}

	return strings.Join(values, "\n")
}
