package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		"scrape_collector_duration_seconds",
		"Duration of a collector scrape.",
		nil,
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		"scrape_collector_success",
		"Whether a collector succeeded.",
		nil,
		nil,
	)
)

// Collector implements the prometheus.Collector interface.
type Collector struct {
}

// NewCollector creates a new collector
func NewCollector() *Collector {
	return &Collector{}
}

// Describe implements the prometheus.Collector interface.
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)

	var success float64
	if err == nil {
		success = 1
	}

	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds())
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success)
}
