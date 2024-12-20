package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"os"

	"log"
	"shop/internal/modules/public/jobs"
)

func main() {

	schedule := asynq.NewScheduler(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_DB"), os.Getenv("REDIS_PORT"))}, &asynq.SchedulerOpts{})

	registerSchedules(schedule)

	err := schedule.Start()
	if err != nil {
		log.Fatal("[x] error start scheduler :", err)
	}

}

// registerSchedules is responsible to register our tasks with specific cronspec
func registerSchedules(schedule *asynq.Scheduler) {

	//---------------------- example task
	exampleTask, exampleTaskErr := jobs.TaskExample("your data ")
	if exampleTaskErr == nil {
		_, err := schedule.Register("@every 3s", exampleTask)
		if err != nil {
			log.Println("[x] error `exampleTaskErr` run schedule:", err)
			return
		}
	}
	//------------------------
}
