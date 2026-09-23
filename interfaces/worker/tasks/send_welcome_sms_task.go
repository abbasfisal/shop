package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

//this is just and example

const TypeSendWelcomeSMS = "user:welcome:sms"

func TaskSendWelcomeSMS(data any) (*asynq.Task, error) {
	payload, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(TypeSendWelcomeSMS, " | json marshal err:", err)
		return nil, err
	}

	// e.g. return asynq.NewTask(TypeSendWelcomeSMS, payload,asynq.Queue("low")) =>pass some opts
	return asynq.NewTask(TypeSendWelcomeSMS, payload), nil
}

func HandleTaskSendWelcomeSMS(ctx context.Context, t *asynq.Task) error {
	var data any
	err := json.Unmarshal(t.Payload(), &data)
	if err != nil {
		log.Println("unmarshal json got err:", err)
		return err
	}

	// implement business logic

	fmt.Println(TypeSendWelcomeSMS, " implement business logic ;)")

	//-------------------------

	return nil
}
