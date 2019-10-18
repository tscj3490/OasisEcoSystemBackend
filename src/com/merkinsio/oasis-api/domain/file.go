package domain

import (
	"time"
)

//File the file data we store in the db
type File struct {
	ID         string    `json:"id" bson:"id"`
	FileName   string    `json:"fileName" bson:"fileName"`
	UploadedAt time.Time `json:"uploadedAt" bson:"uploadedAt"`
}
