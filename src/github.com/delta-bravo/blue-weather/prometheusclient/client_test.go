package prometheusclient

import (
	"testing"
	"github.com/prometheus/client_golang/prometheus"
)

func Test_UpdateTemperatureGaugeValueShouldSetExpectedGaugeValue(t *testing.T) {

	client := GetPrometheusClient()
	registry := prometheus.DefaultGatherer

	client.UpdateTemperatureGaugeValue(42)

	families, e := registry.Gather()
	if e != nil {
		panic(e)
	}
	expected := false
	for _, family := range families {
		if *family.Name == "ambient_temperature" {
			for _, metric := range family.Metric {
				if *metric.Gauge.Value == 42 {
					expected = true
				}
			}

		}
	}
	if !expected {
		t.Error("Temperature gauge expected value not set or metric not found")
	}
}
