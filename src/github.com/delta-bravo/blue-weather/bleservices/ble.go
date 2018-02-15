package bleservices

import (
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"log"
	"time"
	"golang.org/x/net/context"
	"strings"
)

const temperatureServiceUuid = "E95D6100251D470AA062FA1922DFA9A8"
const temperatureCharacteristicUuid = "E95D9250251D470AA062FA1922DFA9A8"

var ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), time.Second*120))

func filter(a ble.Advertisement) bool {
	return strings.Contains(a.LocalName(), "micro:bit")
}

func StartReadingMicrobitTemperature(temperatureHandler ble.NotificationHandler) {
	device, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("Unable to init device: %s\n", err)
		return
	}
	ble.SetDefaultDevice(device)

	log.Println("Connecting...")

	client, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Fatalf("Unable to connect: %s\n", err)
		return
	}

	profile, err := client.DiscoverProfile(true)

	if err != nil {
		log.Fatalf("Unable to discover device profile: %s\n", err)
		return
	}

	discoverAndSubscribeTemperatureService(profile, client, temperatureHandler)
}

func discoverAndSubscribeTemperatureService(profile *ble.Profile, client ble.Client, temperatureHandler ble.NotificationHandler) {
	for _, service := range profile.Services {
		serviceUuid, _ := ble.Parse(temperatureServiceUuid)
		if service.UUID.Equal(serviceUuid) {
			log.Println("Temperature service discovered service, yay")
			subscribeForTemperatureUpdates(service, client, temperatureHandler)
		}
	}
}

func subscribeForTemperatureUpdates(service *ble.Service, client ble.Client, temperatureHandler ble.NotificationHandler) {
	for _, characteristic := range service.Characteristics {
		characteristicUuid, _ := ble.Parse(temperatureCharacteristicUuid)
		if characteristic.UUID.Equal(characteristicUuid) {
			log.Println("Temperature characteristic discovered!")
			client.Subscribe(characteristic, true, temperatureHandler)
		}
	}
}
