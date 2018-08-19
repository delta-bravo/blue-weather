package prometheusclient

import (
	"testing"
	"github.com/prometheus/client_golang/prometheus"
)

func Test_UpdateTemperatureGaugeValueShouldSetExpectedGaugeValue(t *testing.T) {

	client := GetPrometheusClient()
	registry := prometheus.DefaultGatherer

	client.UpdateTemperatureGaugeValue(42)
	client.UpdateBearing(1138)

	families, e := registry.Gather()
	if e != nil {
		panic(e)
	}
	expectedTemperature := false
	expectedBearing := false
	for _, family := range families {
		if *family.Name == "ambient_temperature" {
			for _, metric := range family.Metric {
				if *metric.Gauge.Value == 42 {
					expectedTemperature = true
				}
			}

		} else if *family.Name == "current_bearing" {
			for _, metric := range family.Metric {
				if *metric.Gauge.Value == 1138 {
					expectedBearing = true
				}
			}

		}
	}
	if !expectedTemperature {
		t.Error("Temperature gauge expected value not set or metric not found")
	}
	if !expectedBearing {
		t.Error("Bearing gauge expected value not set or metric not found")
	}
}
