package main

import (
	"encoding/hex"
	"net/http"

	"github.com/cuongcb/micro-event/svcs/dbproxy"
	"github.com/gin-gonic/gin"
)

type eventServiceHandler struct {
	dbHandler *dbproxy.MongoDBLayer
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	event, err := eh.dbHandler.AddEvent(event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func main() {
	dbHandler, err := dbproxy.NewMongoDBLayer("")
	if err != nil {
		panic(err)
	}

	esh := eventServiceHandler{dbHandler: dbHandler}
	router := gin.Default()

	v1 := router.Group("/v1/event")
	{
		v1.POST("/", esh.newEventHandler)
		v1.GET("/", esh.allEventsHandler)
		v1.GET("/search", esh.findEventHandler)
	}

	router.RunTLS(":8081", "./cert/cert.pem", "./cert/key.pem")
}
