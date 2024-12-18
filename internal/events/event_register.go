package events

// RegisterEvents register all your event/listener
func RegisterEvents(em *EventManager) {

	// e.g. user.created event
	em.Register(UserCreatedEvent, SendWelcomeNotification, UserCreatedListener)

	// add another event/listener

}
