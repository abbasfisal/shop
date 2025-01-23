package events

import (
	"context"
	"fmt"
	"shop/internal/modules/public/jobs"
)

const SendEmailEvent = "send.email"

type SendEmailPayload struct {
	ID    int
	Name  string
	Email string
	Text  string
}

func SendEmailListener(ctx context.Context, data any) {
	payload := data.(SendEmailPayload)
	select {
	case <-ctx.Done():
		fmt.Println("SendEmailListener: Execution canceled or timed out")
		return
	default:
		dep := GetDep()
		result, err := dep.RedisClient.Ping(ctx).Result()

		//
		//dep.AsynqClient.Enqueue(jobs.TaskSendEmail(jobs.SendEmailPayload{Name: "ali"}), asynq.ProcessIn(time.Second*10))
		dep.AsynqClient.Enqueue(jobs.TaskSendEmail(jobs.SendEmailPayload{Name: "ali"}))

		if err != nil {
			fmt.Println("redis ping failed")
		} else {
			fmt.Println("redis ping success:" + result)
		}

		fmt.Printf("SendEmailListener executed: %v\n", payload)

	}
}
