package domain

import "time"

type SocialAccounts struct {
	Google   *SocialAccount `json:"google" bson:"google"`
	Facebook *SocialAccount `json:"facebook" bson:"facebook"`
}

type SocialAccount struct {
	ID        string    `json:"id" bson:"id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
