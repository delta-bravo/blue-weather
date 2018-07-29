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

type TestBleClient struct {
	discoverProfileParams               DiscoverProfileParams
	subscribeParams                     SubscribeParams
}

func (TestBleClient) Addr() ble.Addr {
	panic("implement me")
}

func (TestBleClient) Name() string {
	panic("implement me")
}

func (TestBleClient) Profile() *ble.Profile {
	panic("implement me")
}

func (client *TestBleClient) DiscoverProfile(force bool) (*ble.Profile, error) {
	client.discoverProfileParams = DiscoverProfileParams{
		called:    true,
		withForce: force,
	}
	return nil, nil
}

func (TestBleClient) DiscoverServices(filter []ble.UUID) ([]*ble.Service, error) {
	panic("implement me")
}

func (TestBleClient) DiscoverIncludedServices(filter []ble.UUID, s *ble.Service) ([]*ble.Service, error) {
	panic("implement me")
}

func (TestBleClient) DiscoverCharacteristics(filter []ble.UUID, s *ble.Service) ([]*ble.Characteristic, error) {
	panic("implement me")
}

func (TestBleClient) DiscoverDescriptors(filter []ble.UUID, c *ble.Characteristic) ([]*ble.Descriptor, error) {
	panic("implement me")
}

func (TestBleClient) ReadCharacteristic(c *ble.Characteristic) ([]byte, error) {
	panic("implement me")
}

func (TestBleClient) ReadLongCharacteristic(c *ble.Characteristic) ([]byte, error) {
	panic("implement me")
}

func (TestBleClient) WriteCharacteristic(c *ble.Characteristic, value []byte, noRsp bool) error {
	panic("implement me")
}

func (TestBleClient) ReadDescriptor(d *ble.Descriptor) ([]byte, error) {
	panic("implement me")
}

func (TestBleClient) WriteDescriptor(d *ble.Descriptor, v []byte) error {
	panic("implement me")
}

func (TestBleClient) ReadRSSI() int {
	panic("implement me")
}

func (TestBleClient) ExchangeMTU(rxMTU int) (txMTU int, err error) {
	panic("implement me")
}

func (client *TestBleClient) Subscribe(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
	subscribeParams := SubscribeParams{
		characteristic: *c,
		ind:            ind,
		handler:        h,
	}
	client.subscribeParams = subscribeParams
	return nil
}

func (TestBleClient) Unsubscribe(c *ble.Characteristic, ind bool) error {
	panic("implement me")
}

func (TestBleClient) ClearSubscriptions() error {
	panic("implement me")
}

func (TestBleClient) CancelConnection() error {
	panic("implement me")
}

func (TestBleClient) Disconnected() <-chan struct{} {
	panic("implement me")
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
