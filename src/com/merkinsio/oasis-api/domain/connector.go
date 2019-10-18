package domain

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

// DB exported
var DB *mgo.Database

//MongoConnector Holds the MongoDB database session
type MongoConnector struct {
	Session *mgo.Session
}

// NewConnector Establish a new MongoDB connection
func NewConnector(url string) *MongoConnector {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(fmt.Errorf("Error connecting to the database: %s", err))
	}
	DB = session.DB("")
	return &MongoConnector{Session: session}
}

// Close Closes the MongoDB connection
func (c *MongoConnector) Close() {
	c.Session.Close()
}
