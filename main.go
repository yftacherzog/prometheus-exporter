package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a struct for you collector that contains pointers
// to prometheus descriptors for each metric you wish to expose.
// Note you can also include fields of other types if they provide utility
// but we just won't be exposing them as metrics.
type fooCollector struct {
	fakeUp    *prometheus.Desc
	fooMetric *prometheus.Desc
	barMetric *prometheus.Desc
}

// You must create a constructor for you collector that
// initializes every descriptor and returns a pointer to the collector
func newFooCollector() *fooCollector {
	return &fooCollector{
		fakeUp: prometheus.NewDesc("fake_up", "Shows whether the fake app is up",
			nil, nil),
		fooMetric: prometheus.NewDesc("foo_metric",
			"Shows whether a foo has occurred in our cluster",
			nil, nil,
		),
		barMetric: prometheus.NewDesc("bar_metric",
			"Shows whether a bar has occurred in our cluster",
			nil, nil,
		),
	}
}

func getFakeUp(mean float64) int {
	randomValue := rand.Float64()
	if randomValue < mean {
		return 1
	} else {
		return 0
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *fooCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.fakeUp
	ch <- collector.fooMetric
	ch <- collector.barMetric
}

// Collect implements required collect function for all promehteus collectors
func (collector *fooCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	if 1 == 1 {
		metricValue += rand.Float64()
	}

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	m0 := prometheus.MustNewConstMetric(
		collector.fakeUp, prometheus.GaugeValue, float64(getFakeUp(0.75)),
	)
	m1 := prometheus.MustNewConstMetric(collector.fooMetric, prometheus.GaugeValue, metricValue)
	m2 := prometheus.MustNewConstMetric(collector.barMetric, prometheus.GaugeValue, metricValue)
	// m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
	// m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	ch <- m0
	ch <- m1
	ch <- m2
}

func main() {
	reg := prometheus.NewPedanticRegistry()

	foo := newFooCollector()
	reg.MustRegister(foo)

	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
			// Pass custom registry
			Registry: reg,
		},
	))
	log.Fatal(http.ListenAndServe(":9101", nil))
}
