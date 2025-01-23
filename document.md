

## How to Use Asynq with Events

### Step-by-Step Guide

#### 1. Create an Event File
Create a new event file such as `send_email_event.go` in the `internal/events` directory. You can also add multiple listeners to the event.

#### 2. Trigger a Job with Asynq
To trigger a job with Asynq, you need to use `dep := GetDep()` to get access to the dependencies and then enqueue a job using `dep.AsynqClient.Enqueue()`.

#### 3. Register Your Event
Register your event in the `event_register.go` file.

#### 4. Create a Job in the `jobs` Directory
Under the `internal/modules/admin` or `internal/modules/public` directory, create a job file like `send_email_task.go`. Here’s an example:

```go
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
    fmt.Println("-------------- Task send email called")
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
        fmt.Println("Find error:", ferr)
    } else {
        fmt.Println("Users:", u)
    }

    err := json.Unmarshal(t.Payload(), &payload)
    if err != nil {
        fmt.Println("Log SendEmailJob: Error:", err)
        return err
    }

    fmt.Println("SendEmailJob: Payload:", payload)

    return nil
}
```

#### 5. Register Your Jobs in the `cmd/job/worker/main.go`
In the `cmd/job/worker/main.go`, register your job handlers:

```go
package main

import (
    "fmt"
    "github.com/hibiken/asynq"
    "log"
    "os"
    adminJob "shop/internal/modules/admin/jobs"
    "shop/internal/modules/public/jobs"
    "shop/internal/pkg/bootstrap"
)

func main() {
    // Initialize dependencies
    dep, err := bootstrap.Initialize()
    if err != nil {
        log.Fatalf("[x] Failed to initialize dependencies: %v", err)
    }

    // Configure Asynq server
    server := asynq.NewServer(
        asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOSTNAME"), os.Getenv("REDIS_PORT"))},
        asynq.Config{
            Concurrency: 10,
        },
    )

    errs := server.Ping()
    if errs != nil {
        log.Fatal("------[Ping error]:", errs)
        return
    }

    // Mux server setup
    mux := asynq.NewServeMux()

    // Register job handlers
    mux.Handle(adminJob.CancelPendingOrders, adminJob.NewCancelJob(dep))
    mux.Handle(jobs.TypeSendEmail, jobs.NewSendEmailJob(dep))

    // Run server
    err = server.Run(mux)
    if err != nil {
        log.Fatal("[x] Job worker start failed:", err)
    }
}
```

#### 6. Trigger Events Using `EventManager`
To trigger an event using the `EventManager`, follow the steps below in your API:

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
    "shop/internal/database/mongodb"
    "shop/internal/database/mysql"
    "shop/internal/events"
    "shop/internal/pkg/bootstrap"
)

func main() {
    // Initialize dependencies
    dependencies, err := bootstrap.Initialize()
    if err != nil {
        log.Fatal("[x] Error initializing project:", err)
    }
    defer mysql.Close()
    defer mongodb.Disconnect()
    defer dependencies.AsynqClient.Close()

    // Configure and start the Gin server
    r := gin.Default()
    setupRoutes(r, dependencies)

    fmt.Printf("http://localhost:8484/")
    log.Fatalln(r.Run("localhost:8484"))
}

func setupRoutes(r *gin.Engine, dependencies *bootstrap.Dependencies) {
    // Create and configure EventManager
    em := events.NewEventManager(&events.EventManagerDep{
        AsynqClient: dependencies.AsynqClient,
        DB:          dependencies.DB,
        RedisClient: dependencies.RedisClient,
        MongoClient: dependencies.MongoClient,
    })

    // Register all events
    events.RegisterEvents(em)

    // Create handler
    hndlr := NewClientHandler(dependencies, em)

    // Define route
    r.GET("/client/send-email", hndlr.ClientSendEmail)
}

type ClientHandler struct {
    dep *bootstrap.Dependencies
    em  *events.EventManager
}

func NewClientHandler(dep *bootstrap.Dependencies, em *events.EventManager) *ClientHandler {
    return &ClientHandler{
        dep: dep,
        em:  em,
    }
}

func (ch *ClientHandler) ClientSendEmail(c *gin.Context) {
    // Prepare payload for your event
    payload := events.SendEmailPayload{
        ID:    10,
        Name:  "Ali",
        Email: "ali@gmail.com",
        Text:  "Welcome dear Ali, please give us some money.",
    }

    // Trigger your event
    ch.em.Emit(c.Request.Context(), events.SendEmailEvent, payload, true)

    c.JSON(200, gin.H{
        "payload": payload,
        "message": "Email sent",
        "keys":    ch.dep.RedisClient.Ping(c.Request.Context()),
    })
}
```

---

Let me know if you need further adjustments or explanations!