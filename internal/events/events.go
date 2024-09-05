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
	eventType := getEventType(event)
	if listeners, ok := e.listeners[eventType]; ok {
		for _, listener := range listeners {
			go listener(event)
		}
	}
}

func getEventType(event interface{}) string {
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
