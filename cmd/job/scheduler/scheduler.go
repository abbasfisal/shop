package scheduler

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"os"
	"shop/internal/events"
	adminJob "shop/internal/modules/admin/jobs"
	"shop/internal/pkg/bootstrap"

	"log"
)

func RunScheduler(ctx context.Context, dep *bootstrap.Dependencies, em *events.EventManager) {

	schedule := asynq.NewScheduler(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOSTNAME"), os.Getenv("REDIS_PORT"))}, &asynq.SchedulerOpts{})

	if err := schedule.Ping(); err != nil {
		log.Fatal("[x] scheduler ping failed:", err)
	}

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
