package emitter

import "github.com/tj/go-debug"

var dbg = debug.Debug("dbus:emitter")

//Callback is a function to be invoked when an event happens
type Callback func(ev Event)

//Event contains information about what happened
type Event interface {
	GetName() string
	GetData() interface{}
}

//BaseEvent contains unspecialized information about what happened
type BaseEvent struct {
	name string
	data interface{}
}

//GetName return the event name
func (e BaseEvent) GetName() string {
	return e.name
}

//GetData return the event data
func (e BaseEvent) GetData() interface{} {
	return e.data
}

var pipe chan Event
var events = make(map[string][]Callback, 0)

func loop() {
	dbg("loop: Started")
	for {
		if pipe == nil {
			dbg("loop: Closed")
			return
		}
		ev := <-pipe
		dbg("loop: Trigger event %s", ev.GetName())
		if events[ev.GetName()] != nil {
			size := len(events[ev.GetName()])
			if size == 0 {
				return
			}
			for i := 0; i < size; i++ {
				dbg("Call fn")
				events[ev.GetName()][i](ev)
			}
		}
	}
}

func getPipe() {
	if pipe == nil {
		pipe = make(chan Event)
		go loop()
	}
}

//On registers to an event
func On(event string, callback Callback) {

	if event == "" {
		panic("Cannot use an empty string as event name")
	}

	if events[event] == nil {
		getPipe()
		events[event] = make([]Callback, 0)
	}

	events[event] = append(events[event], callback)
	dbg("Added %s events, len is %d", event, len(events[event]))
}

// Emit an event
func Emit(name string, data interface{}) {
	dbg("Emit event %s\n", name)
	getPipe()
	ev := BaseEvent{name, data}
	pipe <- ev
}

//Off Removes all callbacks from an event
func Off(name string) {
	dbg("Off %s", name)
	if name == "*" {
		for name := range events {
			if name != "*" {
				Off(name)
			}
		}
	}

	if events[name] != nil {
		delete(events, name)
	}

	if len(events) == 0 {
		close(pipe)
		pipe = nil // will stop the go routine
	}

}