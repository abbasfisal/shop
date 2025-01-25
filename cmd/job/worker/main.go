package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os"
	adminJob "shop/internal/modules/admin/jobs"
	"shop/internal/pkg/bootstrap"
)

func main() {
	//load bootstrap
	dep, err := bootstrap.Initialize()
	if err != nil {
		log.Fatalf("[x] failed to initialize dependencies: %v", err)
	}

	//config asynq server
	server := asynq.
		NewServer(
			asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOSTNAME"), os.Getenv("REDIS_PORT"))},
			asynq.Config{
				Concurrency: 10,
				//Queues:                   nil,
			},
		)

	errs := server.Ping()
	if errs != nil {
		log.Fatal("------[ping err:]", errs)
		return
	}

	//mux server
	mux := asynq.NewServeMux()
	//mux.HandleFunc(publicJob.TypeSendWelcomeSMS, publicJob.HandleTaskSendWelcomeSMS)
	//mux.HandleFunc(publicJob.TypeExample, publicJob.HandleExampleTask)

	mux.Handle(adminJob.CancelPendingOrders, adminJob.NewCancelJob(dep))
	//mux.Handle(jobs.TypeSendEmail, jobs.NewSendEmailJob(dep))
	//>>>>> run serve r<<<<<
	err = server.Run(mux)
	if err != nil {
		log.Fatal("[x] job worker start failed:", err)
	}
}
