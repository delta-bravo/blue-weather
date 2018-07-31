package bleservices_test

import (
	"testing"
	"github.com/go-ble/ble"
	"github.com/delta-bravo/blue-weather/bleservices"
)

const _temperatureServiceUuid = "E95D6100251D470AA062FA1922DFA9A8"
const _temperatureCharacteristicUuid = "E95D9250251D470AA062FA1922DFA9A8"

type TestBluetoothClient struct {
	profile           ble.Profile
	err               error
	profileDiscovered bool
	testHandlerPassed bool
	subscribeParams   SubscribeParams
	testHandler       func(data []byte)
}

type SubscribeParams struct {
	characteristic ble.Characteristic
	ind            bool
	handler        ble.NotificationHandler
}

type DiscoverProfileParams struct {
	called    bool
	withForce bool
}

func (client *TestBluetoothClient) createTestHandler() {
	client.testHandler = func(data []byte) {
		client.testHandlerPassed = true
	}
}

func (client *TestBluetoothClient) setupGoodProfile() {
	serviceUuid, e := ble.Parse(_temperatureServiceUuid)
	characteristicUuid, e := ble.Parse(_temperatureCharacteristicUuid)
	if e != nil {
		panic("Failed to parse expected UUIDs")
	}
	service := ble.NewService(serviceUuid)
	service.Characteristics = []*ble.Characteristic{service.NewCharacteristic(characteristicUuid)}
	services := []*ble.Service{service}

	client.profile = ble.Profile{
		Services: services,
	}
	client.err = nil
}

func (client *TestBluetoothClient) setupError(err error) {
	client.err = err
}

func (client *TestBluetoothClient) DiscoverProfile() (*ble.Profile, error) {
	client.profileDiscovered = true
	return &client.profile, nil
}

func (client *TestBluetoothClient) Subscribe(characteristic *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
	client.subscribeParams = SubscribeParams{
		characteristic: *characteristic,
		ind:            ind,
		handler:        h,
	}
	client.testHandler([]byte{})
	return client.err
}

func (client *TestBluetoothClient) Initialize() error {
	return nil
}

func TestStartBluetoothServices(t *testing.T) {
	//given
	testBluetoothClient := &TestBluetoothClient{}
	testBluetoothClient.setupGoodProfile()
	testBluetoothClient.createTestHandler()

	//when
	bleservices.StartBluetoothServices(testBluetoothClient, testBluetoothClient.testHandler)

	//then
	if !testBluetoothClient.profileDiscovered {
		t.Error("Expected to call discover profile, never did")
	}

	subscribeParams := testBluetoothClient.subscribeParams

	tempCharacteristicUUID, e := ble.Parse(_temperatureCharacteristicUuid)
	if e != nil {
		panic("Failed to parse expected UUID")
	}

	if !subscribeParams.characteristic.UUID.Equal(tempCharacteristicUUID) {
		t.Errorf("Expected uuid %s but got %s", tempCharacteristicUUID, subscribeParams.characteristic.UUID.String())
	}

	if subscribeParams.ind {
		t.Error("Expected notification subscription (ind = false) but got (ind = true)")
	}

	if !testBluetoothClient.testHandlerPassed {
		t.Error("Didn't pass expected ble handler")
	}
}

func TestDiscoverProfile(t *testing.T) {
	//given
	mockBleClient := TestBleClient{
	}
	bleClientUnderTest := bleservices.BluetoothClientImpl{
		BleClient: &mockBleClient,
	}

	//when
	bleClientUnderTest.DiscoverProfile()

	//then
	if !mockBleClient.discoverProfileParams.called {
		t.Error("Expected BLE Client's DiscoverProfile to have been called")
	}

	if !mockBleClient.discoverProfileParams.withForce {
		t.Error("Expected BLE Client's DiscoverProfile to have been called with force = true")
	}
}

func TestSubscribe(t *testing.T) {
	//given
	mockBleClient := TestBleClient{
	}

	bleClientWrapperUnderTest := bleservices.BluetoothClientImpl{
		BleClient: &mockBleClient,
	}
	serviceUuid, e := ble.Parse(_temperatureServiceUuid)
	characteristicUuid, e := ble.Parse(_temperatureCharacteristicUuid)
	if e != nil {
		panic("Failed to parse expected UUIDs")
	}

	service := ble.NewService(serviceUuid)
	characteristic := service.NewCharacteristic(characteristicUuid)

	//when
	bleClientWrapperUnderTest.Subscribe(characteristic, false, nil)

	//then
	if !mockBleClient.subscribeParams.characteristic.UUID.Equal(characteristicUuid) {
		t.Error("Expected BLE Client's characteristic UUID to be equal")
	}

	if mockBleClient.subscribeParams.ind {
		t.Error("Expected BLE Client's ind to be false")
	}
}
