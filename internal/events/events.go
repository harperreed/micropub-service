package events

import (
	"sync"
)

type EventEmitter struct {
	listeners map[string][]func(interface{})
	mu        sync.RWMutex
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]func(interface{})),
	}
}

func (e *EventEmitter) On(eventType string, listener func(interface{})) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[eventType] = append(e.listeners[eventType], listener)
}

func (e *EventEmitter) Emit(event interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	eventType := GetEventType(event)
	if listeners, ok := e.listeners[eventType]; ok {
		for _, listener := range listeners {
			go listener(event)
		}
	}
}

func (e *EventEmitter) Off(eventType string, listener func(interface{})) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if listeners, ok := e.listeners[eventType]; ok {
		for i, l := range listeners {
			if &l == &listener {
				e.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

func GetEventType(event interface{}) string {
	switch event.(type) {
	case FileEvent:
		return "file"
	default:
		return "unknown"
	}
}

type FileEvent struct {
	Type     string
	Filename string
}
