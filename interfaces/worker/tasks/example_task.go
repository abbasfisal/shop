package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TypeExample = "example:task"

func TaskExample(data any) (*asynq.Task, error) {
	payload, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(TypeExample, " | json marshal err:", err)
		return nil, err
	}

	// e.g. return asynq.NewTask(TypeExample, payload,asynq.Queue("low")) =>pass some opts
	return asynq.NewTask(TypeExample, payload), nil
}

func HandleExampleTask(ctx context.Context, t *asynq.Task) error {
	var data any
	err := json.Unmarshal(t.Payload(), &data)
	if err != nil {
		log.Println("unmarshal json got err:", err)
		return err
	}

	// implement business logic

	fmt.Println(TypeExample, " implement business logic ;)")

	//-------------------------

	return nil
}
