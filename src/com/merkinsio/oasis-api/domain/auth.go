package domain

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

//AuthTokenValidTime The valid time of the JWT token
const AuthTokenValidTime = time.Hour * 24 * 365 // One year LUL nice security LUL

//TokenizedUser Holds the JWT tokenized data
type TokenizedUser struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Email    string        `json:"email" bson:"email"`
	Role     string        `json:"role" bson:"role"`
	Password string        `json:"-" bson:"password"`
}

//TokenClaims Holds the JWT token claims
type TokenClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

//GetTokenizedUserByEmail Retrierves the tokenized user by its email
func GetTokenizedUserByEmail(email string) (*TokenizedUser, error) {
	tokenized := TokenizedUser{}

	if err := DB.C("users").Find(bson.M{
		"email":  email,
		"active": true,
	}).One(&tokenized); err != nil {
		return nil, err
	}

	return &tokenized, nil
}
