package main

import (
	"github.com/delta-bravo/blue-weather/bleservices"
	"github.com/delta-bravo/blue-weather/httpserver"
	"github.com/delta-bravo/blue-weather/prometheusclient"
	"github.com/delta-bravo/blue-weather/handlers"
	"log"
)

func main() {
	log.Println("Staring App")
	bluetoothClient := bleservices.GetClient()
	prometheusClient := prometheusclient.GetPrometheusClient()
	bluetoothHandlers:= handlers.CreateBluetoothHandlers(prometheusClient)

	go func() {
		err := bleservices.StartBluetoothServices(bluetoothClient, bluetoothHandlers.CreateTemperatureHandler())
		if err!=nil {
			panic(err)
		}
	}()
	httpServer := httpserver.CreateServer(bluetoothHandlers, prometheusClient)
	httpServer.StartServer(":8015")
}
