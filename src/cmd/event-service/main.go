package main

import (
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/cuongcb/micro-event/dbproxy"
	"github.com/cuongcb/micro-event/msgqueue"
	queueAmqp "github.com/cuongcb/micro-event/msgqueue/amqp"
	"github.com/cuongcb/micro-event/msgqueue/contracts"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type eventServiceHandler struct {
	dbHandler    *dbproxy.MongoDBLayer
	eventEmitter msgqueue.Emitter
}

type findEventRequest struct {
	ID   *string `form:"id"`
	Name *string `form:"name"`
}

func (eh *eventServiceHandler) findEventHandler(ctx *gin.Context) {
	var req findEventRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	var event dbproxy.Event
	switch {
	case req.ID != nil:
		id, err := hex.DecodeString(*req.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "invalid request")
			return
		}

		event, err = eh.dbHandler.FindEvent(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "invalid request")
			return
		}
	case req.Name != nil:
		e, err := eh.dbHandler.FindEventByName(*req.Name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "invalid request")
			return
		}

		event = e
	}

	ctx.JSON(http.StatusOK, event)
}

func (eh *eventServiceHandler) allEventsHandler(ctx *gin.Context) {
	events, err := eh.dbHandler.FindAllAvailableEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func (eh *eventServiceHandler) newEventHandler(ctx *gin.Context) {
	var event dbproxy.Event
	if err := ctx.Bind(&event); err != nil {
		log.Println("binding error:", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	event, err := eh.dbHandler.AddEvent(event)
	if err != nil {
		log.Printf("adding event %v error %v\n", event, err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	msg := &contracts.EventCreatedEvent{
		ID:         event.ID.String(),
		Name:       event.Name,
		LocationID: event.Location.ID.String(),
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
	}

	eh.eventEmitter.Emit(msg)

	ctx.JSON(http.StatusOK, event)
}

func main() {
	log.SetPrefix("[event-service] ")
	dbHandler, err := dbproxy.NewMongoDBLayer("")
	if err != nil {
		panic(err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	emitter, err := queueAmqp.NewEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	esh := eventServiceHandler{dbHandler: dbHandler, eventEmitter: emitter}
	router := gin.Default()

	v1 := router.Group("/v1/event")
	{
		v1.POST("/", esh.newEventHandler)
		v1.GET("/", esh.allEventsHandler)
		v1.GET("/search", esh.findEventHandler)
	}

	go router.Run(":8080")

	router.RunTLS(":8081", "./cert/cert.pem", "./cert/key.pem")
}
