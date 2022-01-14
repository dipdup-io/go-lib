package prometheus

import (
	"log"
	"net/http"
	"sync"

	"github.com/dipdup-net/go-lib/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Service -
type Service struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
	gauge      map[string]*prometheus.GaugeVec
	server     *http.Server
	wg         sync.WaitGroup
}

// NewService -
func NewService(cfg *config.Prometheus) *Service {
	var s Service
	s.counters = make(map[string]*prometheus.CounterVec)
	s.histograms = make(map[string]*prometheus.HistogramVec)
	s.gauge = make(map[string]*prometheus.GaugeVec)

	if cfg != nil && cfg.URL != "" {
		s.server = &http.Server{Addr: cfg.URL}
		http.Handle("/metrics", promhttp.Handler())
	}

	return &s
}

// Start -
func (service *Service) Start() {
	if service.server == nil {
		return
	}

	service.wg.Add(1)
	go func() {
		defer service.wg.Done()

		if err := service.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
}

// Close -
func (service *Service) Close() error {
	if service.server != nil {
		if err := service.server.Close(); err != nil {
			return err
		}
	}

	service.wg.Wait()

	return nil
}

// RegisterCounter -
func (service *Service) RegisterCounter(name, help string, labels ...string) {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	}, labels)
	service.counters[name] = vec
	prometheus.MustRegister(vec)
}

// Counter -
func (service *Service) Counter(name string) *prometheus.CounterVec {
	counter, ok := service.counters[name]
	if ok {
		return counter
	}
	return nil
}

// IncrementCounter -
func (service *Service) IncrementCounter(name string, labels map[string]string) {
	counter, ok := service.counters[name]
	if ok {
		counter.With(labels).Inc()
	}
}

// RegisterGoBuildMetrics -
func (service *Service) RegisterGoBuildMetrics() {
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
}

// RegisterHistogram -
func (service *Service) RegisterHistogram(name, help string, labels ...string) {
	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name,
		Help: help,
	}, labels)
	service.histograms[name] = vec
	prometheus.MustRegister(vec)
}

// Histogram -
func (service *Service) Histogram(name string) *prometheus.HistogramVec {
	histogram, ok := service.histograms[name]
	if ok {
		return histogram
	}
	return nil
}

// AddHistogramValue -
func (service *Service) AddHistogramValue(name string, labels map[string]string, observe float64) {
	histogram, ok := service.histograms[name]
	if ok {
		histogram.With(labels).Observe(observe)
	}
}

// RegisterGauge -
func (service *Service) RegisterGauge(name, help string, labels ...string) {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, labels)
	service.gauge[name] = vec
	prometheus.MustRegister(vec)
}

// Gauge -
func (service *Service) Gauge(name string) *prometheus.GaugeVec {
	gauge, ok := service.gauge[name]
	if ok {
		return gauge
	}
	return nil
}

// SetGaugeValue -
func (service *Service) SetGaugeValue(name string, labels map[string]string, observe float64) {
	gauge, ok := service.gauge[name]
	if ok {
		gauge.With(labels).Set(observe)
	}
}

// IncGaugeValue -
func (service *Service) IncGaugeValue(name string, labels map[string]string) {
	gauge, ok := service.gauge[name]
	if ok {
		gauge.With(labels).Inc()
	}
}

// DecGaugeValue -
func (service *Service) DecGaugeValue(name string, labels map[string]string) {
	gauge, ok := service.gauge[name]
	if ok {
		gauge.With(labels).Dec()
	}
}
