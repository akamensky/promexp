package promexp

type seriesType string

const (
	gauge     seriesType = "gauge"
	counter   seriesType = "counter"
	summary   seriesType = "summary"
	histogram seriesType = "histogram"
)
