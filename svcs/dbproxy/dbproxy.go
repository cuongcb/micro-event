package dbproxy

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// DB ..
	DB = "micro_events"
	// USERS ...
	USERS = "users"
	// EVENTS ...
	EVENTS = "events"
)

const (
	mongoURI = "mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb"
)

// MongoDBLayer ...
type MongoDBLayer struct {
	session *mgo.Session
}

// NewMongoDBLayer ...
func NewMongoDBLayer(uri string) (*MongoDBLayer, error) {
	if uri == "" {
		uri = mongoURI
	}

	s, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	return &MongoDBLayer{session: s}, nil
}

func (m *MongoDBLayer) getFreshSession() *mgo.Session {
	return m.session.Copy()
}

// AddEvent ...
func (m *MongoDBLayer) AddEvent(e Event) (Event, error) {
	s := m.getFreshSession()
	defer s.Close()

	if e.ID.Valid() == false {
		e.ID = bson.NewObjectId()
	}

	if e.Location.ID.Valid() == false {
		e.Location.ID = bson.NewObjectId()
	}

	return e, s.DB(DB).C(EVENTS).Insert(e)
}

// FindEvent ...
func (m *MongoDBLayer) FindEvent(id []byte) (Event, error) {
	s := m.getFreshSession()
	defer s.Close()

	var e Event
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)

	return e, err
}

// FindEventByName ...
func (m *MongoDBLayer) FindEventByName(name string) (Event, error) {
	s := m.getFreshSession()
	defer s.Close()

	var e Event
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)

	return e, err
}

// FindAllAvailableEvents ...
func (m *MongoDBLayer) FindAllAvailableEvents() ([]Event, error) {
	s := m.getFreshSession()
	defer s.Close()

	var events []Event
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)

	return events, err
}
