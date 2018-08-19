package handlers

import (
	"net/http"
	"testing"
)

type testPrometheusClient struct {
	handlerReturned  bool
	temperatureGauge float64
	bearing          float64
}

type testHandler struct {
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (client *testPrometheusClient) UpdateTemperatureGaugeValue(value float64) {
	client.temperatureGauge = value
}

func (client *testPrometheusClient) UpdateBearing(value float64) {
	client.bearing = value
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

func TestCreateBluetoothHandlers_ShouldCreateHandlersWithExpectedPrometheusAndUpdateBearing(t *testing.T) {
	prometheusClient := &testPrometheusClient{}
	bluetoothHandlers := CreateBluetoothHandlers(prometheusClient)

	bearingHandler := bluetoothHandlers.CreateBearingHandler()

	bearingHandler([]byte{180})

	if prometheusClient.bearing != 180 {
		t.Error("Expected to set Prometheus client temperature to correct value")
	}
	if bluetoothHandlers.GetCurrentBearing() != 180 {
		t.Error("Expected to set handlers current temperature to correct value")
	}
}
