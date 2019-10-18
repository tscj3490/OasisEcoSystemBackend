package services

import (
	"com/merkinsio/oasis-api/domain"

	log "github.com/sirupsen/logrus"
)

//StoreUserPayload Holds the payload that the StoreUser method needs
type StoreUserPayload struct {
	User *domain.User
}

//StoreUser Creates and return the newly created user in the database
func StoreUser(payload StoreUserPayload) (*domain.User, error) {
	//TODO: Handle duplicate users and restrictions
	userID, err := domain.CreateUserByUser(payload.User)

	if err != nil {
		return nil, err
	}

	user, err := domain.GetUserByID(*userID)

	if err != nil {
		log.Errorf("Error in services.StoreUser -> error: %s", err.Error())
		return nil, err
	}

	return user, nil
}
