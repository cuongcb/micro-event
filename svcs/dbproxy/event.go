package dbproxy

import (
	"gopkg.in/mgo.v2/bson"
)

// Event ...
type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}
