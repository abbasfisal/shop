package events

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

var (
	once                 sync.Once
	eventManagerInstance *EventManager
)

type ListenerFunc func(ctx context.Context, data any)

type eventName string

type EventManager struct {
	dep       *EventManagerDep
	listeners map[eventName][]ListenerFunc
	mu        sync.RWMutex
}

type EventManagerDep struct {
	AsynqClient *asynq.Client
	DB          *gorm.DB
	RedisClient *redis.Client
	MongoClient *mongo.Client
}

func NewEventManager(dep *EventManagerDep) *EventManager {
	once.Do(func() {
		eventManagerInstance = &EventManager{
			dep:       dep,
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
	defer em.mu.Unlock()

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

func GetDep() *EventManagerDep {
	return eventManagerInstance.dep
}
