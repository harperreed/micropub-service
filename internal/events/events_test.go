package events

import (
	"sync"
	"sync/atomic"
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
	eventType := "file"
	eventData := "test.txt"
	received := make(chan interface{}, 1)

	emitter.On(eventType, func(data interface{}) {
		received <- data
	})

	go emitter.Emit(FileEvent{Type: eventType, Filename: eventData})

	select {
	case data := <-received:
		if fe, ok := data.(FileEvent); !ok || fe.Filename != eventData {
			t.Errorf("Expected FileEvent with Filename %v, got %v", eventData, data)
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

func TestEventEmitter_MultipleListeners(t *testing.T) {
	emitter := NewEventEmitter()
	eventType := "test"
	eventData := "test_data"
	listenerCount := 3
	var callCount int32

	for i := 0; i < listenerCount; i++ {
		emitter.On(eventType, func(data interface{}) {
			atomic.AddInt32(&callCount, 1)
		})
	}

	emitter.Emit(FileEvent{Type: eventType, Filename: eventData})

	// Wait for all listeners to be called
	time.Sleep(100 * time.Millisecond)

	if int(atomic.LoadInt32(&callCount)) != listenerCount {
		t.Errorf("Expected %d listener calls, got %d", listenerCount, callCount)
	}
}

func TestEventEmitter_RemoveListener(t *testing.T) {
	emitter := NewEventEmitter()
	eventType := "test"
	eventData := "test_data"
	var callCount int32

	listener := func(data interface{}) {
		atomic.AddInt32(&callCount, 1)
	}

	emitter.On(eventType, listener)
	emitter.Off(eventType, listener)

	emitter.Emit(FileEvent{Type: eventType, Filename: eventData})

	// Wait for potential listener calls
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&callCount) != 0 {
		t.Errorf("Expected 0 listener calls after removal, got %d", callCount)
	}

	// Test removing a non-existent listener
	nonExistentListener := func(data interface{}) {}
	emitter.Off(eventType, nonExistentListener) // This should not panic
}

func TestEventEmitter_ErrorHandling(t *testing.T) {
	emitter := NewEventEmitter()
	eventType := "test"
	eventData := "test_data"
	var normalListenerCalled bool
	var panicListenerCalled bool

	emitter.On(eventType, func(data interface{}) {
		panic("This listener panics")
	})

	emitter.On(eventType, func(data interface{}) {
		normalListenerCalled = true
	})

	// This deferred function will recover from the panic and set panicListenerCalled
	defer func() {
		if r := recover(); r != nil {
			panicListenerCalled = true
		}
	}()

	emitter.Emit(FileEvent{Type: eventType, Filename: eventData})

	// Wait for all listeners to be called
	time.Sleep(100 * time.Millisecond)

	if !normalListenerCalled {
		t.Error("Normal listener was not called")
	}

	if !panicListenerCalled {
		t.Error("Panic listener did not panic as expected")
	}
}
