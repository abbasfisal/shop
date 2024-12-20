package events

import "context"

// e.g., this main.go file is just an example of how to use events and listeners.

func main() {

	em := NewEventManager()

	//register events
	RegisterEvents(em)

	//in your routes e.g. /users/welcome
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
