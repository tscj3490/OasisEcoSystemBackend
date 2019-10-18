package domain

import "gopkg.in/mgo.v2/bson"

type Town struct {
	Entity       `bson:",inline"`
	Name         string `json:"name" bson:"name"`
	Code         int    `json:"code" bson:"code"`
	Province     int    `json:"province" bson:"province"`
	ProvinceName string `json:"provinceName" bson:"provinceName"`
}

func GetTownsByProvince(province int) ([]Town, error) {
	towns := []Town{}
	// TODO: ROLES diffs
	if err := DB.C("towns").Find(bson.M{"province": province}).All(&towns); err != nil {
		return towns, err
	}

	return towns, nil
}
