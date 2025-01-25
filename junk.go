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
	dependencies, err := bootstrap.Initialize()
	if err != nil {
		log.Fatal("[x] error initializing project :", err)
	}
	defer mysql.Close()
	defer mongodb.Disconnect()
	defer dependencies.AsynqClient.Close()

	//
	r := gin.Default()
	setupRoutes(r, dependencies)

	//
	fmt.Printf("http://localhost:8484/")
	log.Fatalln(r.Run("localhost:8484"))
}

func setupRoutes(r *gin.Engine, dependencies *bootstrap.Dependencies) {

	em := events.NewEventManager(&events.EventManagerDep{
		AsynqClient: dependencies.AsynqClient,
		DB:          dependencies.DB,
		RedisClient: dependencies.RedisClient,
		MongoClient: dependencies.MongoClient,
	})
	events.RegisterEvents(em)

	//
	hndlr := NewClientHandler(dependencies, em)
	//client
	r.GET("/client/a", hndlr.ClientSendEmail)
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

	payload := events.SendEmailPayload{
		ID:    10,
		Name:  "ali",
		Email: "ali@gmail.com",
		Text:  "welcome dear ali , please give us some money",
	}

	ch.em.Emit(c.Request.Context(), events.SendEmailEvent, payload, true)

	c.JSON(200, gin.H{
		"payload": payload,
		"message": "email sent",
		"keys":    ch.dep.RedisClient.Ping(c.Request.Context()),
	})
	return
}
