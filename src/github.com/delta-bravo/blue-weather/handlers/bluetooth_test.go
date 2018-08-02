package handlers

import (
	"net/http"
	"testing"
)

type testPrometheusClient struct {
	handlerReturned  bool
	temperatureGauge float64
}

type testHandler struct {
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (client *testPrometheusClient) UpdateTemperatureGaugeValue(value float64) {
	client.temperatureGauge = value
}

func (client *testPrometheusClient) GetPromHttpHandler() http.Handler {
	client.handlerReturned = true
	return &testHandler{}
}

func TestCreateBluetoothHandlers_ShouldCreateHandlersWithExpectedPrometheusAndUpdateTemperature(t *testing.T) {
	prometheusClient := &testPrometheusClient{}
	bluetoothHandlers := CreateBluetoothHandlers(prometheusClient)

	tempHandler := bluetoothHandlers.CreateTemperatureHandler()

	tempHandler([]byte{42})

	if prometheusClient.temperatureGauge != 42 {
		t.Error("Expected to set Prometheus client temperature to correct value")
	}
	if bluetoothHandlers.GetCurrentTemperature() != 42 {
		t.Error("Expected to set handlers current temperature to correct value")
	}
}
