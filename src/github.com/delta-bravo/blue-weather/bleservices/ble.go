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

type BluetoothClient interface {
	Initialize() error
	DiscoverProfile() (*ble.Profile, error)
	Subscribe(characteristic *ble.Characteristic, ind bool, h ble.NotificationHandler) error
}

type BluetoothClientImpl struct {
	BleClient ble.Client
}

func (client *BluetoothClientImpl) Initialize() error {
	device, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("Unable to init device: %s\n", err)
		return err
	}
	ble.SetDefaultDevice(device)
	log.Println("Connecting...")
	bluetoothClient, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Fatalf("Unable to connect: %s\n", err)
		return err
	}
	log.Println("Connected")
	client.BleClient = bluetoothClient
	return nil
}

func (client BluetoothClientImpl) DiscoverProfile() (*ble.Profile, error) {
	return client.BleClient.DiscoverProfile(true)
}

func (client BluetoothClientImpl) Subscribe(characteristic *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
	log.Println("Subscribing for characteristic!")
	return client.BleClient.Subscribe(characteristic, ind, h)
}

func GetClient() BluetoothClient {
	bluetoothClient := &BluetoothClientImpl{}
	err := bluetoothClient.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize client: %s\n", err)
		return nil
	}
	return bluetoothClient
}

var ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), time.Second*120))

func filter(a ble.Advertisement) bool {
	return strings.Contains(a.LocalName(), "micro:bit")
}

func StartBluetoothServices(bluetoothClient BluetoothClient, notificationHandler ble.NotificationHandler) error {
	profile, err := bluetoothClient.DiscoverProfile()
	temperatureServiceUuid, err := ble.Parse(temperatureServiceUuid)
	temperatureCharacteristicUuid, err := ble.Parse(temperatureCharacteristicUuid)
	if err != nil {
		return err
	}
	for _, service := range profile.Services {
		if service.UUID.Equal(temperatureServiceUuid) {
			log.Println("Temperature service discovered service, yay")
			return subscribeForCharacteristicWithService(service.Characteristics, temperatureCharacteristicUuid,
				bluetoothClient, notificationHandler)
		}
	}
	return nil
}

func subscribeForCharacteristicWithService(characteristics []*ble.Characteristic, uuid ble.UUID,
	client BluetoothClient, handler ble.NotificationHandler) error {
	for _, characteristic := range characteristics {
		if characteristic.UUID.Equal(uuid) {
			return client.Subscribe(characteristic, false, handler)
		}
	}
	return nil
}
