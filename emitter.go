package goble

import (
	"log"
)

const (
	ALL = "__allEvents__"
)

type Event struct {
	Name       string
	State      string
	DeviceUUID UUID
	Peripheral Peripheral
}

type EventHandlerFunc func(Event)

type Emitter struct {
	handlers map[string]EventHandlerFunc
	event    chan Event
}

func (e *Emitter) Init() {
	e.handlers = make(map[string]EventHandlerFunc)
	e.event = make(chan Event)

	// event handler
	go func() {
		for {
			ev := <-e.event

			if fn, ok := e.handlers[ev.Name]; ok {
				fn(ev)
			} else if fn, ok := e.handlers[ALL]; ok {
				fn(ev)
			} else {
				log.Println("unhandled Emit", ev)
			}
		}
	}()
}

func (e *Emitter) Emit(ev Event) {
	e.event <- ev
}

func (e *Emitter) On(event string, fn EventHandlerFunc) {
	if fn == nil {
		delete(e.handlers, event)
	} else {
		e.handlers[event] = fn
	}
}
