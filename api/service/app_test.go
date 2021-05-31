package service

import (
	"testing"
	"unsafe"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createTestApp(t *testing.T) *App {

	log.SetLevel(log.TraceLevel)

	a, err := NewApp(AppOptions{
		AdapterID: api.GetDefaultAdapterID(),
	})
	if err != nil {
		t.Fatal(err)
	}

	s1, err := a.NewService("2233")
	if err != nil {
		t.Fatal(err)
	}

	c1, err := s1.NewChar("3344")
	if err != nil {
		t.Fatal(err)
	}

	c1.
		OnRead(CharReadCallback(func(c *Char, options map[string]interface{}) ([]byte, error) {
			return nil, nil
		})).
		OnWrite(CharWriteCallback(func(c *Char, value []byte) ([]byte, error) {
			return nil, nil
		}))

	d1, err := c1.NewDescr("4455")
	if err != nil {
		t.Fatal(err)
	}

	err = c1.AddDescr(d1)
	if err != nil {
		t.Fatal(err)
	}

	err = s1.AddChar(c1)
	if err != nil {
		t.Fatal(err)
	}

	err = a.AddService(s1)
	if err != nil {
		t.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		t.Fatal(err)
	}

	return a
}

func TestApp(t *testing.T) {
	a := createTestApp(t)
	defer a.Close()
}

func TestAppAgentIsNil(t *testing.T) {
	a := new(App)
	assert.Equal(t, nil, a.agent)
}

func TestAppPassCodePersistsWithCustomAgent(t *testing.T) {
	passCode := "043210"

	ag := agent.NewSimpleAgent()
	ag.SetPassCode(passCode)
	opt := AppOptions{CustomAgent: ag}

	a, err := NewApp(opt)
	if err != nil {
		t.Fatal(err)
	}

	appAgent := a.Agent()
	appAgentAsSimple := (*agent.SimpleAgent)(unsafe.Pointer(&appAgent)) // This might be bad

	assert.Equal(t, passCode, appAgentAsSimple.PassCode())
}
