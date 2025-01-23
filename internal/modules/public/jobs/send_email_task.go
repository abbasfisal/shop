package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"shop/internal/entities"
	"shop/internal/pkg/bootstrap"
)

const TypeSendEmail = "send:email"

type SendEmailPayload struct {
	Name string
}

func TaskSendEmail(payload SendEmailPayload) *asynq.Task {
	fmt.Println("-------------- task send email called ")
	payloadJson, _ := json.Marshal(payload)
	return asynq.NewTask(TypeSendEmail, payloadJson)
}

type SendEmailJob struct {
	dep *bootstrap.Dependencies
}

func NewSendEmailJob(dep *bootstrap.Dependencies) *SendEmailJob {
	fmt.Println("----------- New Send Email Job called")
	return &SendEmailJob{dep: dep}
}

func (sj *SendEmailJob) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Println("------------ Process Task for SendEmailJob called")
	var payload SendEmailPayload

	d := sj.dep.DB
	var u entities.User
	ferr := d.WithContext(ctx).Find(&u).Error
	if ferr != nil {
		fmt.Println("find err:", ferr)
	} else {
		fmt.Println("users :", u)
	}

	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		fmt.Println("log SendEmailJob : err : ", err)
		return err
	}

	fmt.Println("SendEmailJob : payload : ", payload)

	return nil
}
