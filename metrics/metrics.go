package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RunDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "spiegel",
		Subsystem: "main",
		Name:      "run_duration_seconds",
		Help:      "Duration of each run in seconds",
		Buckets:   prometheus.LinearBuckets(10, 10, 30),
	}, []string{"episodes_downloaded", "succeeded"})

	TotalScrapes = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "spiegel",
		Subsystem: "scrapers",
		Name:      "scrapes_total",
		Help:      "Total number of scrapes performed by each scraper",
	}, []string{"scraper", "succeeded"})

	ScraperRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "spiegel",
		Subsystem: "scrapers",
		Name:      "request_duration_seconds",
		Help:      "Duration of each request in seconds",
		Buckets:   prometheus.ExponentialBuckets(0.005, 2, 12),
	}, []string{"scraper", "code", "method"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "spiegel",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Duration of each request in seconds",
		Buckets:   prometheus.ExponentialBuckets(0.005, 2, 12),
	}, []string{"path", "code", "method"})

	TotalEpisodesByShow = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "spiegel",
		Subsystem: "main",
		Name:      "episodes_downloaded",
		Help:      "Total number of episodes downloaded per show",
	}, []string{"show"})
)
