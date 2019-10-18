package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type WebsocketMessage struct {
	Type    string          `json:"type"`
	Data    string          `json:"data"`
	ZoneID  bson.ObjectId   `json:"-"`
	UserIDs []bson.ObjectId `json:"-"`
}
