package domain

type Crop struct {
	Entity `bson:",inline"`

	Name string `json:"name" bson:"name"`
}
