package prometheusclient

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

type Client interface {
	UpdateTemperatureGaugeValue(value float64)
	UpdateBearing(value float64)
	GetPromHttpHandler() http.Handler
}

type prometheusClient struct {
	handler          http.Handler
	temperatureGauge prometheus.Gauge
	bearingGauge prometheus.Gauge
}

func (client *prometheusClient) UpdateTemperatureGaugeValue(value float64) {
	client.temperatureGauge.Set(value)
}

func (client *prometheusClient) UpdateBearing(value float64) {
	client.bearingGauge.Set(value)
}

func (client *prometheusClient) GetPromHttpHandler() http.Handler {
	return client.handler
}

func GetPrometheusClient() Client {
	temperatureGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ambient_temperature",
		Help: "Ambient temperature reading",
	})

	bearingGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "current_bearing",
		Help: "Current bearing reading",
	})

	prometheus.Register(temperatureGauge)
	prometheus.Register(bearingGauge)
	return &prometheusClient{
		handler:          promhttp.Handler(),
		temperatureGauge: temperatureGauge,
		bearingGauge: bearingGauge,
	}
}
