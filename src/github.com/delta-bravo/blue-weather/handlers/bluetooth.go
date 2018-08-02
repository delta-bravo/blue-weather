package handlers

import "github.com/delta-bravo/blue-weather/prometheus"

type BluetoothHandlers interface {
	CreateTemperatureHandler() func(data []byte)
	GetCurrentTemperature() float64
}

type GaugeHandlers struct {
	temperatureGauge float64
	prometheusClient prometheus.Client
}

func (handlers *GaugeHandlers) CreateTemperatureHandler() func(data []byte) {
	return func(data []byte) {
		handlers.temperatureGauge = float64(data[0])
		handlers.prometheusClient.UpdateTemperatureGaugeValue(handlers.temperatureGauge)
	}
}

func (handlers *GaugeHandlers) 	GetCurrentTemperature() float64 {
	return handlers.temperatureGauge
}

func CreateBluetoothHandlers(prometheusClient prometheus.Client) BluetoothHandlers {
	return &GaugeHandlers{
		prometheusClient:prometheusClient,
	}
}
