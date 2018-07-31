package bleservices_test

import "github.com/go-ble/ble"

type TestBleClient struct {
	discoverProfileParams DiscoverProfileParams
	subscribeParams       SubscribeParams
	err error
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
	return nil, client.err
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
	return client.err
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
