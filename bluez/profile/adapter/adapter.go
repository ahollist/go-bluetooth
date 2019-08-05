package adapter

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

func (a *Adapter1) GetAdapterID() (string, error) {
	return ParseAdapterID(a.Path())
}

var defaultAdapterName = "hci0"

func SetDefaultAdapterName(a string) {
	defaultAdapterName = a
}

func GetDefaultAdapterName() string {
	return defaultAdapterName
}

// ParseAdapterID read the adapterID from an object path in the form /org/bluez/hci[0-9]*[/...]
func ParseAdapterID(path dbus.ObjectPath) (string, error) {

	spath := string(path)

	if !strings.HasPrefix(spath, bluez.OrgBluezPath) {
		return "", fmt.Errorf("Failed to parse adapterID from %s", path)
	}

	parts := strings.Split(spath[len(bluez.OrgBluezPath)+1:], "/")
	adapterID := parts[0]

	if adapterID[:3] != "hci" {
		return "", fmt.Errorf("adapterID missing hci* prefix from %s", path)
	}

	return adapterID, nil
}

// AdapterExists checks if an adapter is available
func AdapterExists(adapterID string) (bool, error) {

	om, err := bluez.GetObjectManager()
	if err != nil {
		return false, err
	}

	objects, err := om.GetManagedObjects()
	if err != nil {
		return false, err
	}

	path := dbus.ObjectPath(fmt.Sprintf("%s/%s", bluez.OrgBluezPath, adapterID))
	_, exists := objects[path]

	return exists, nil
}

func GetDefaultAdapter() (*Adapter1, error) {
	return GetAdapter(GetDefaultAdapterName())
}

// GetAdapter return an adapter object instance
func GetAdapter(adapterID string) (*Adapter1, error) {

	if exists, err := AdapterExists(adapterID); !exists {
		if err != nil {
			return nil, fmt.Errorf("AdapterExists: %s", err)
		}
		return nil, fmt.Errorf("Adapter %s not found", adapterID)
	}

	return NewAdapter1FromAdapterID(adapterID)
}

func GetAdapterFromDevicePath(path dbus.ObjectPath) (*Adapter1, error) {

	d, err := device.NewDevice1(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load device %s", path)
	}

	a, err := NewAdapter1(d.Properties.Adapter)
	if err != nil {
		return nil, err
	}

	return a, nil
}
