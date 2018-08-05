package httpserver

import (
	"github.com/delta-bravo/blue-weather/handlers"
	"github.com/delta-bravo/blue-weather/prometheusclient"
	"net/http"
	"html/template"
	"fmt"
)

type HttpServer interface {
	StartServer(addr string)
	addHandlers()
}

func CreateServer(handlers handlers.BluetoothHandlers, prometheusClient prometheusclient.Client) HttpServer {
	return &httpServer{
		prometheusClient:  prometheusClient,
		bluetoothHandlers: handlers,
	}
}

type httpServer struct {
	bluetoothHandlers handlers.BluetoothHandlers
	prometheusClient  prometheusclient.Client
}

func (server *httpServer) StartServer(addr string) {
	server.addHandlers()
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func (server *httpServer) addHandlers() {
	http.Handle("/metrics", server.prometheusClient.GetPromHttpHandler())
	htmlTemplate, err := template.ParseFiles("index.html")

	if err != nil {
		panic("Unparsable HTML")
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		htmlTemplate.Execute(writer, map[string]interface{}{"temperature": server.bluetoothHandlers.GetCurrentTemperature()})
	})
	http.HandleFunc("/temperature", server.refreshHandler)
	http.HandleFunc("/style.css", createFileHandler("style.css"))
	http.HandleFunc("/refresh.js", createFileHandler("refresh.js"))
}

func (server *httpServer) refreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, server.bluetoothHandlers.GetCurrentTemperature())
}

func createFileHandler(fileName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fileName)
	}
}
