package domain

import "gopkg.in/mgo.v2/bson"

type Province struct {
	Entity `bson:",inline"`
	Name   string `json:"name" bson:"name"`
	Code   int    `json:"code" bson:"code"`
}

func GetProvinces() ([]Province, error) {
	provinces := []Province{}
	if err := DB.C("provinces").Find(bson.M{}).All(&provinces); err != nil {
		return provinces, err
	}

	return provinces, nil
}
