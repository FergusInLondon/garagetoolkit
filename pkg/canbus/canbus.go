package canbus

import (
	"github.com/brutella/can"
)

type CANBusListener can.Handler

var (
	bus *can.Bus
)

func Connect(iface string) (err error) {
	bus, err = can.NewBusForInterfaceWithName(iface)
	return err
}

func RegisterListener(listener can.Handler) {
	bus.Subscribe(listener)
}

func Run() error {
	return bus.ConnectAndPublish()
}

func Stop() error {
	return bus.Disconnect()
}
