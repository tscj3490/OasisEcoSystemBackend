package server

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/mail"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// GetAllUsers Returns all the users
func GetAllUsers(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	users := []domain.User{}
	collection := domain.DB.C("users")

	collection.Find(bson.M{}).Select(bson.M{"password": 0}).All(&users)

	respondWithJSON(w, http.StatusOK, users)
}

//CreateUser Creates and returns the created user
func CreateUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// currentUser := context.Get(r, "currentUser").(*domain.User)
	userOwn := domain.UserAndOwner{}

	if err := parseBody(r, &userOwn); err != nil {
		log.Errorf("Error in server.CreateUser.parseBody -> error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := userOwn.User

	clearPassword := user.Password
	user.Email = strings.ToLower(user.Email)
	if user.Role == "" {
		user.Role = constants.RoleFarmer
	}

	// Check if the email is already in the database
	retrievedUser, _ := domain.GetUserByEmail(user.Email)
	if retrievedUser != nil {
		respondWithError(w, http.StatusConflict, "duplicateEmail")
		return
	}

	payload := services.StoreUserPayload{
		User: &user,
	}
	createdUser, err := services.StoreUser(payload)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set User as Cooperative Owner
	userOwn.User = *createdUser
	err = domain.SetCooperativeOwner(userOwn)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Send the email with the password
	emailBody, err := mail.ParseEmailFromTemplate("sendRegistrationPassword.html", mail.SendRegistrationPassword{
		Password: clearPassword,
		Name:     payload.User.FirstName,
	})
	if err != nil {
		log.Printf("Error in mail.ParseEmailFromTemplate: %v", err)
	}
	req := mail.CreateRequest([]string{payload.User.Email}, "Contraseña", emailBody)
	_, err = req.SendEmail()
	if err != nil {
		log.Printf("Error in req.SendMail: %v", err)
	}

	respondWithJSON(w, http.StatusOK, &createdUser)
}

//GetUserByID returns the user whose Id you put on the query
func GetUserByID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	userIDStr := vars["userId"]
	var userID bson.ObjectId

	if isValid := bson.IsObjectIdHex(userIDStr); isValid == false {
		respondWithError(w, http.StatusBadRequest, "invalid id")
	} else {
		userID = bson.ObjectIdHex(userIDStr)
	}

	var err error
	var user *domain.User
	user, err = domain.GetUserByID(userID)

	if err != nil {
		errorCode := http.StatusBadRequest
		if err.Error() == "not found" {
			errorCode = http.StatusNotFound
		}

		respondWithError(w, errorCode, fmt.Sprintf("User not found for id %v", userIDStr))
	} else {
		respondWithJSON(w, http.StatusOK, user)
	}
}

//RegisterUser Creates a client user
func RegisterUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := domain.User{}

	if err := parseBody(r, &user); err != nil {
		log.Errorf("Error in server.RegisterUser.parseBody -> error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Email = strings.ToLower(user.Email)
	user.Role = constants.RoleFarmer
	user.TermOfServices = false

	// Check if the email is already in the database
	retrievedUser, _ := domain.GetUserByEmail(user.Email)
	if retrievedUser != nil {
		respondWithError(w, http.StatusConflict, "duplicateEmail")
		return
	}

	payload := services.StoreUserPayload{
		User: &user,
	}
	createdUser, err := services.StoreUser(payload)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, &createdUser)
}

// TermOfServicesAcceptedByUser accepted by user
func TermOfServicesAcceptedByUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := domain.User{}

	if err := parseBody(r, &user); err != nil {
		log.Errorf("Error in server.TermOfServicesAcceptedByUser.parseBody -> error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdUser, err := domain.TermOfServicesAcceptedByUserID(user)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, &createdUser)
}

// GetAllTechnicians Returns all the users with technician role
func GetAllTechnicians(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	users := []domain.User{}
	collection := domain.DB.C("users")

	collection.Find(bson.M{"role": constants.RoleTechnician}).Select(bson.M{"password": 0}).All(&users)

	respondWithJSON(w, http.StatusOK, users)
}

// GetAllWorkers Returns all the users with technician role
func GetAllWorkers(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	users := []domain.User{}
	collection := domain.DB.C("users")

	collection.Find(bson.M{"role": constants.RoleWorker}).Select(bson.M{"password": 0}).All(&users)

	respondWithJSON(w, http.StatusOK, users)
}

// UpdateUser updates ands returns the user
func UpdateUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	userOwn := domain.UserAndOwner{}

	if err := parseBody(r, &userOwn); err != nil {
		log.Errorf("Error in server.CreateUser.parseBody -> error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	body := userOwn.User

	vars := mux.Vars(r)

	userIDStr := vars["userId"]

	if valid := bson.IsObjectIdHex(userIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	userID := bson.ObjectIdHex(userIDStr)
	body.ID = userID

	err := domain.UpdateUserByUser(body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		user, err := domain.GetUserByID(userID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, user)
		}
	}

	// Set User as Cooperative Owner
	err = domain.SetCooperativeOwner(userOwn)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
}

// DeleteUser changes the active property
func DeleteUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	var userID bson.ObjectId

	if isValid := bson.IsObjectIdHex(userIDStr); isValid == false {
		respondWithError(w, http.StatusBadRequest, "invalid id")
	} else {
		userID = bson.ObjectIdHex(userIDStr)
	}

	if err := domain.DeleteUserByID(userID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

func GetClientsByCooperativeID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	cooperativeIDStr := vars["cooperativeId"]

	if valid := bson.IsObjectIdHex(cooperativeIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	cooperativeID := bson.ObjectIdHex(cooperativeIDStr)

	users, err := domain.GetClientsByCooperativeID(&cooperativeID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, users)
	}

}
