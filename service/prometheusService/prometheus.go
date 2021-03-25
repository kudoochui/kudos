package prometheusService

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// RpcClient
	RpcClient = New().
		WithHistogram(
			"kudos_rpcclient_histogram_milliseconds_cost",
			"rpc time consumption histogram",
			[]string{"module", "method"},
			prometheus.ExponentialBuckets(0.001, 10, 8)).
		WithCounter(
			"kudos_rpcclient_counter",
			"rpc times counter",
			[]string{"module", "method"}).
		WithGauge(
			"kudos_rpcclient_state_milliseconds_cost",
			"current rpc time consumption",
			[]string{"module", "method"})
)

func Handler() http.Handler {
	return promhttp.Handler()
}

type Prometheus struct {
	histogram   *prometheus.HistogramVec
	summary *prometheus.SummaryVec
	counter *prometheus.CounterVec
	gauge   *prometheus.GaugeVec
}

// New creates a Prom instance.
func New() *Prometheus {
	return &Prometheus{}
}

// WithHistogram sets HistogramVec
func (p *Prometheus) WithHistogram(name, help string, labels []string, buckets []float64) *Prometheus {
	if p == nil || p.histogram != nil {
		return p
	}

	p.histogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: help,
			Buckets: buckets,
		}, labels)

	prometheus.MustRegister(p.histogram)
	return p
}

// WithSummary with summary timer
func (p *Prometheus) WithSummary(name, help string, labels []string, objectives map[float64]float64) *Prometheus {
	if p == nil || p.summary != nil {
		return p
	}

	p.summary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: name,
			Help: help,
			Objectives: objectives,
		}, labels)

	prometheus.MustRegister(p.summary)

	return p
}

// WithCounter sets counter.
func (p *Prometheus) WithCounter(name, help string, labels []string) *Prometheus {
	if p == nil || p.counter != nil {
		return p
	}

	p.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, labels)

	prometheus.MustRegister(p.counter)

	return p
}

// WithState sets state.
func (p *Prometheus) WithGauge(name, help string, labels []string) *Prometheus {
	if p == nil || p.gauge != nil {
		return p
	}

	p.gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels)

	prometheus.MustRegister(p.gauge)

	return p
}

// Timing log timing information (in milliseconds) without sampling
func (p *Prometheus) Timing(name string, time float64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.histogram != nil {
		p.histogram.WithLabelValues(label...).Observe(time)
	}
}

func (p *Prometheus) Sampling(name string, time float64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.summary != nil {
		p.summary.WithLabelValues(label...).Observe(time)
	}
}

func (p *Prometheus) CounterIncr(name string, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Inc()
	}
}

func (p *Prometheus) CounterAdd(name string, v float64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Add(v)
	}
}

// Incr increments one stat counter without sampling
func (p *Prometheus) Incr(name string, extra ...string) {
	label := append([]string{name}, extra...)
	if p.gauge != nil {
		p.gauge.WithLabelValues(label...).Inc()
	}
}

// Decr decrements one stat counter without sampling
func (p *Prometheus) Decr(name string, extra ...string) {
	if p.gauge != nil {
		label := append([]string{name}, extra...)
		p.gauge.WithLabelValues(label...).Dec()
	}
}

// State set state
func (p *Prometheus) Set(name string, v float64, extra ...string) {
	if p.gauge != nil {
		label := append([]string{name}, extra...)
		p.gauge.WithLabelValues(label...).Set(v)
	}
}

// Add add count v must > 0
func (p *Prometheus) Add(name string, v float64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.gauge != nil {
		p.gauge.WithLabelValues(label...).Add(v)
	}
}