//example

package events

import (
	"context"
	"fmt"
)

// UserCreatedEvent event name
const UserCreatedEvent = "user.created"

// SendWelcomeNotification is a listener
func SendWelcomeNotification(ctx context.Context, data any) {
	notificationData := data.(map[eventName][]any)
	select {
	case <-ctx.Done():
		fmt.Println("SendWelcomeNotification : Execution canceled or timed out")
		return
	default:
		fmt.Printf("SendWelcomeNotification executed: %v\n", notificationData)
	}
}

// UserCreatedListener is a listener
func UserCreatedListener(ctx context.Context, data any) {
	userData := data.(map[string]interface{})
	select {
	case <-ctx.Done():
		fmt.Println("UserCreatedListener: Execution canceled or timed out")
		return
	default:
		fmt.Printf("UserCreatedListener executed: %v\n", userData)
	}
}
