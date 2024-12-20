package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os"
	publicJob "shop/internal/modules/public/jobs"
)

func main() {
	server := asynq.
		NewServer(
			asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_DB"))},
			asynq.Config{
				Concurrency: 10,
				//Queues:                   nil,
			},
		)

	//mux server
	mux := asynq.NewServeMux()
	mux.HandleFunc(publicJob.TypeSendWelcomeSMS, publicJob.HandleTaskSendWelcomeSMS)
	mux.HandleFunc(publicJob.TypeExample, publicJob.HandleExampleTask)

	//>>>>> run serve r<<<<<
	err := server.Run(mux)
	if err != nil {
		log.Fatal("[x] job worker start failed:", err)
	}
}
