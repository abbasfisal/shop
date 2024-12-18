package events

import (
	"context"
	"fmt"
	"sync"
)

type ListenerFunc func(ctx context.Context, data any)

type eventName string

type EventManager struct {
	listeners map[eventName][]ListenerFunc
	mu        sync.RWMutex
}

var (
	once                 sync.Once
	eventManagerInstance *EventManager
)

func NewEventManager() *EventManager {
	once.Do(func() {
		eventManagerInstance = &EventManager{
			listeners: make(map[eventName][]ListenerFunc),
		}
	})
	return eventManagerInstance
}

func (em *EventManager) Register(eventName eventName, listeners ...ListenerFunc) {
	em.listeners[eventName] = append(em.listeners[eventName], listeners...)
}

func (em *EventManager) Emit(ctx context.Context, eventName eventName, data any, async bool) {
	em.mu.Lock()
	defer em.mu.RUnlock()

	if listeners, ok := em.listeners[eventName]; ok {
		for _, listener := range listeners {
			if async {
				go listener(ctx, data)
			} else {
				listener(ctx, data)
			}
		}
	} else {
		fmt.Printf("No listeners found for event: %s\n", eventName)
	}
}
