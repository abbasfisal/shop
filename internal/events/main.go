package events

import (
	"context"
	"log"
	"shop/internal/pkg/bootstrap"
)

// e.g., this main.go file is just an example of how to use events and listeners.

func main() {
	dep, err := bootstrap.Initialize()
	if err != nil {
		log.Fatalln("fail to initialize bootstrap in event main.go file :", err)
	}

	evntDep := EventManagerDep{
		AsynqClient: dep.AsynqClient,
		DB:          dep.DB,
		RedisClient: dep.RedisClient,
		MongoClient: dep.MongoClient,
	}
	em := NewEventManager(&evntDep)

	//register events
	RegisterEvents(em)

	// in your routes e.g. /users/welcome
	userWelcomePayload := map[string]any{
		"user_id": 1,
		"name":    "ali",
		"address": map[string]any{
			"city":        "tehran",
			"street_name": "pasdaran",
		},
	}
	em.Emit(context.Background(), UserCreatedEvent, userWelcomePayload, true)
}
