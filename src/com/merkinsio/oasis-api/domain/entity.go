package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Entity Holds the shared properties between the domain entities
type Entity struct {
	ID        bson.ObjectId `bson:"_id" json:"id" valid:"required"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt" valid:"required"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt" valid:"required"`
	Active    bool          `bson:"active" json:"active" valid:"required"`
}

//InitializeNewData initializes the entity data to a newly created one
func (e *Entity) InitializeNewData() {
	e.ID = bson.NewObjectId()
	e.CreatedAt = time.Now().UTC()
	e.UpdatedAt = time.Now().UTC()
	e.Active = true
}

//UpdateData updates the entity data
func (e *Entity) UpdateData() {
	e.UpdatedAt = time.Now().UTC()
}
