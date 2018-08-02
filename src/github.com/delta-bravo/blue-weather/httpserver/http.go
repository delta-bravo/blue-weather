package httpserver

import "github.com/delta-bravo/blue-weather/handlers"

type HttpServer interface {
	StartServer(addr string, bluetoothHandlers handlers.BluetoothHandlers)
	addHandlers()
}







