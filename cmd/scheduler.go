package cmd

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"os"
	"shop/bootstrap"
	"shop/infrastructure/events"
	adminJob "shop/interfaces/worker/adminjobs"

	"log"
)

func RunScheduler(ctx context.Context, dep *bootstrap.Dependencies, em *events.EventManager) {

	schedule := asynq.NewScheduler(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOSTNAME"), os.Getenv("REDIS_PORT"))}, &asynq.SchedulerOpts{})

	registerSchedules(schedule, dep)

	err := schedule.Run()
	if err != nil {
		log.Fatal("[x] error start scheduler :", err)
	}

}

// registerSchedules is responsible to register our tasks with specific cronspec
func registerSchedules(schedule *asynq.Scheduler, dep *bootstrap.Dependencies) {

	//---------------------- example task blueprint
	//	exampleTask, exampleTaskErr := jobs.TaskExample("your data ")
	//	if exampleTaskErr == nil {
	//		_, err := schedule.Register("@every 3s", exampleTask)
	//		if err != nil {
	//			log.Println("[x] error `exampleTaskErr` run schedule:", err)
	//			return
	//		}
	//}
	//------------------------

	_, err := schedule.Register("@every 1m", adminJob.TaskCancelPendingOrders())
	if err != nil {
		log.Println("CancelPendingOrders err:", err)
	}

}

func init() {
	rootCmd.AddCommand(schedulerCmd)
}

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "run asynq scheduler",
	Run: func(cmd *cobra.Command, args []string) {
		RunScheduler(context.Background(), nil, nil)
	},
}
