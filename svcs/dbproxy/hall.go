package dbproxy

import (
	"gopkg.in/mgo.v2/bson"
)

// Hall ...
type Hall struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string        `json:"name"`
	Location string        `json:"location,omitempty"`
	Capacity int           `json:"capacity"`
}
