package service_example

import (
	"time"

	"github.com/muka/go-bluetooth/api/service"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	"github.com/muka/go-bluetooth/hw/linux/btmgmt"
	log "github.com/sirupsen/logrus"
)

func serve(adapterID string) error {
	btman := btmgmt.NewBtMgmt(adapterID)
	err := btman.SetSsp(true)
	if err != nil {
		return err
	}

	simpleAgent := agent.NewSimpleAgent()
	simpleAgent.SetPassKey(1111)
	simpleAgent.SetPassCode("1111")
	options := service.AppOptions{
		AdapterID:   adapterID,
		AgentCaps:   agent.CapDisplayOnly,
		UUIDSuffix:  "-0000-1000-8000-00805F9B34FB",
		UUID:        "1234",
		CustomAgent: simpleAgent,
	}

	a, err := service.NewApp(options)
	if err != nil {
		return err
	}
	defer a.Close()

	a.SetName("go_bluetooth")

	log.Infof("HW address %s", a.Adapter().Properties.Address)

	if !a.Adapter().Properties.Powered {
		err = a.Adapter().SetPowered(true)
		if err != nil {
			log.Fatalf("Failed to power the adapter: %s", err)
		}
	}

	service1, err := a.NewService("2233")
	if err != nil {
		return err
	}

	err = a.AddService(service1)
	if err != nil {
		return err
	}

	char1, err := service1.NewChar("3344")
	if err != nil {
		return err
	}

	char1.Properties.Flags = []string{
		gatt.FlagCharacteristicEncryptAuthenticatedRead,
		gatt.FlagCharacteristicEncryptAuthenticatedWrite,
	}

	char1.OnRead(service.CharReadCallback(func(c *service.Char, options map[string]interface{}) ([]byte, error) {
		log.Warnf("GOT READ REQUEST")
		return []byte{42}, nil
	}))

	char1.OnWrite(service.CharWriteCallback(func(c *service.Char, value []byte) ([]byte, error) {
		log.Warnf("GOT WRITE REQUEST")
		return value, nil
	}))

	err = service1.AddChar(char1)
	if err != nil {
		return err
	}

	descr1, err := char1.NewDescr("4455")
	if err != nil {
		return err
	}

	descr1.Properties.Flags = []string{
		gatt.FlagCharacteristicEncryptAuthenticatedRead,
		gatt.FlagCharacteristicEncryptAuthenticatedWrite,
	}

	descr1.OnRead(service.DescrReadCallback(func(c *service.Descr, options map[string]interface{}) ([]byte, error) {
		log.Warnf("GOT READ REQUEST")
		return []byte{42}, nil
	}))
	descr1.OnWrite(service.DescrWriteCallback(func(d *service.Descr, value []byte) ([]byte, error) {
		log.Warnf("GOT WRITE REQUEST")
		return value, nil
	}))

	err = char1.AddDescr(descr1)
	if err != nil {
		return err
	}

	err = a.Run()
	if err != nil {
		return err
	}

	log.Infof("Exposed service %s", service1.Properties.UUID)

	timeout := uint32(6 * 3600) // 6h
	log.Infof("Advertising for %ds...", timeout)
	cancel, err := a.Advertise(timeout)
	if err != nil {
		return err
	}

	defer cancel()

	wait := make(chan bool)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		wait <- true
	}()

	<-wait

	return nil
}
