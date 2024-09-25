package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	log "github.com/sirupsen/logrus"
)

const (
	METRIC_PORT int = 5000
)

const (
	metric_namespace = "finance"
	metric_subsystem = "symbol"
)

type NewRelicCounter interface {
	Inc()
}

type NewRelicGauge interface {
	Set(value float64)
}

type NewRelicMetric struct {
	name string
}

var (
	weekday        time.Weekday
	metricCounters map[string]NewRelicCounter
	metricGauges   map[string]NewRelicGauge
)

func init() {
	weekday = time.Now().Weekday()
	metricCounters = make(map[string]NewRelicCounter)
	metricGauges = make(map[string]NewRelicGauge)
	log.Infof("Metrics initialized")
}

func GetOrCreateCounter(name string) NewRelicCounter {
	if metricCounters[name] == nil {
		metricCounters[name] = createCounter(name)
	}
	return metricCounters[name]
}

func GetOrCreateGauge(name string) NewRelicGauge {
	if metricGauges[name] == nil {
		metricGauges[name] = createGauge(name)
	}
	return metricGauges[name]
}

func createCounter(name string) NewRelicCounter {
	return &NewRelicMetric{fmt.Sprintf("%s_%s_%s", metric_namespace, metric_subsystem, name)}
}

func createGauge(name string) NewRelicGauge {
	return &NewRelicMetric{fmt.Sprintf("%s_%s_%s", metric_namespace, metric_subsystem, name)}
}

func (m *NewRelicMetric) Inc() {
	harvester.RecordMetric(telemetry.Count{
		Name:      m.name,
		Value:     1,
		Timestamp: time.Now(),
		Attributes: map[string]interface{}{
			"weekday": weekday,
		},
	})
}

func (m *NewRelicMetric) Set(value float64) {
	harvester.RecordMetric(telemetry.Gauge{
		Name:      m.name,
		Value:     value,
		Timestamp: time.Now(),
		Attributes: map[string]interface{}{
			"weekday": weekday,
		},
	})
}

func HarvestNow() {
	harvester.HarvestNow(context.Background())
}
