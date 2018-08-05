package main

import (
	"github.com/delta-bravo/blue-weather/bleservices"
	"github.com/delta-bravo/blue-weather/httpserver"
	"github.com/delta-bravo/blue-weather/prometheusclient"
	"github.com/delta-bravo/blue-weather/handlers"
)

func main() {
	bluetoothClient := bleservices.GetClient()
	prometheusClient := prometheusclient.GetPrometheusClient()
	bluetoothHandlers:= handlers.CreateBluetoothHandlers(prometheusClient)

	go bleservices.StartBluetoothServices(bluetoothClient, bluetoothHandlers.CreateTemperatureHandler())
	httpServer := httpserver.CreateServer(bluetoothHandlers, prometheusClient)
	httpServer.StartServer(":8015")
}
