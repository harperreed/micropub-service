package events

import (
	"sync"
	"testing"
	"time"
)

func TestNewEventEmitter(t *testing.T) {
	emitter := NewEventEmitter()
	if emitter == nil {
		t.Fatal("NewEventEmitter returned nil")
	}
	if emitter.listeners == nil {
		t.Error("listeners map is nil")
	}
}

func TestEventEmitter_On(t *testing.T) {
	emitter := NewEventEmitter()
	eventType := "test"
	listener := func(interface{}) {}

	emitter.On(eventType, listener)

	if len(emitter.listeners[eventType]) != 1 {
		t.Errorf("Expected 1 listener, got %d", len(emitter.listeners[eventType]))
	}
}

func TestEventEmitter_Emit(t *testing.T) {
	emitter := NewEventEmitter()
	eventType := "test"
	eventData := "test data"
	received := make(chan interface{}, 1)

	emitter.On(eventType, func(data interface{}) {
		received <- data
	})

	emitter.Emit(FileEvent{Type: eventType, Filename: eventData})

	select {
	case data := <-received:
		if fe, ok := data.(FileEvent); !ok || fe.Filename != eventData {
			t.Errorf("Expected %v, got %v", eventData, data)
		}
	case <-time.After(time.Second):
		t.Error("Timed out waiting for event")
	}
}

func TestEventEmitter_EmitUnknownType(t *testing.T) {
	emitter := NewEventEmitter()
	unknownEvent := struct{ message string }{"unknown"}

	// This should not panic
	emitter.Emit(unknownEvent)
}

func TestEventEmitter_ConcurrentAccess(t *testing.T) {
	emitter := NewEventEmitter()
	eventType := "test"
	iterations := 1000
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			emitter.On(eventType, func(interface{}) {})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			emitter.Emit(FileEvent{Type: eventType, Filename: "test"})
		}
	}()

	wg.Wait()

	if len(emitter.listeners[eventType]) != iterations {
		t.Errorf("Expected %d listeners, got %d", iterations, len(emitter.listeners[eventType]))
	}
}

func TestGetEventType(t *testing.T) {
	tests := []struct {
		name     string
		event    interface{}
		expected string
	}{
		{"FileEvent", FileEvent{Type: "test", Filename: "test.txt"}, "file"},
		{"UnknownEvent", struct{}{}, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEventType(tt.event); got != tt.expected {
				t.Errorf("GetEventType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
