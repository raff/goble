package goble

import (
	"log"
)

type Event struct {
	Name          string
	State         string
	DeviceUUID    UUID
	Advertisement Advertisement
	Rssi          int
	Services      []string
}

type EventHandlerFunc func(Event)

type Emitter struct {
	handlers map[string]EventHandlerFunc
}

func (e *Emitter) Emit(ev Event) {
	if e.handlers == nil {
		e.handlers = make(map[string]EventHandlerFunc)
	}

	if fn, ok := e.handlers[ev.Name]; ok {
		fn(ev)
	} else {
		log.Println("unhandled emit", ev)
	}
}

func (e *Emitter) On(event string, fn EventHandlerFunc) {
	if e.handlers == nil {
		e.handlers = make(map[string]EventHandlerFunc)
	}

	if fn == nil {
		delete(e.handlers, event)
	} else {
		e.handlers[event] = fn
	}
}
