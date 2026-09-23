package cmd

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"log"
	"os"
	"shop/infrastructure/events"
	//adminJob "shop/interfaces/worker/adminjobs"
	"shop/bootstrap"
)

func RunWorker(ctx context.Context, dep *bootstrap.Dependencies, em *events.EventManager) {
	//load bootstrap
	//dep, err := bootstrap.Initialize()
	//if err != nil {
	//	log.Fatalf("[x] failed to initialize dependencies: %v", err)
	//}

	//config asynq server
	server := asynq.
		NewServer(
			asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOSTNAME"), os.Getenv("REDIS_PORT"))},
			asynq.Config{
				Concurrency: 10,
				//Queues:                   nil,
			},
		)

	//mux server
	mux := asynq.NewServeMux()
	//mux.HandleFunc(publicJob.TypeSendWelcomeSMS, publicJob.HandleTaskSendWelcomeSMS)
	//mux.HandleFunc(publicJob.TypeExample, publicJob.HandleExampleTask)

	//mux.Handle(adminJob.CancelPendingOrders, adminJob.NewCancelJob(dep, em))
	//mux.Handle(jobs.TypeSendEmail, jobs.NewSendEmailJob(dep))
	//>>>>> run serve r<<<<<
	log.Println("[info] worker started")
	err := server.Run(mux)
	if err != nil {
		log.Fatal("[x] job worker start failed:", err)
	}
}

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "run asynq worker",
	Run: func(cmd *cobra.Command, args []string) {
		RunWorker(context.Background(), nil, nil)
	},
}
