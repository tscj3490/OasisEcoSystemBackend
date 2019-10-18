package server

import (
	"com/merkinsio/oasis-api/config"
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/mail"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/negroni"

	jwt "github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Login Generates the user jwt token
func Login(w http.ResponseWriter, r *http.Request) {
	var body login
	parseBody(r, &body)
	tokenized, err := domain.GetTokenizedUserByEmail(body.Email)

	if err != nil {
		respondWithError(w, 401, "invalidEmailOrPassword")
		return
	}

	noMatch := bcrypt.CompareHashAndPassword([]byte(tokenized.Password), []byte(body.Password))

	if noMatch != nil {
		respondWithError(w, 401, "invalidEmailOrPassword")
	} else {
		_, tokenString, err := generateJWTToken(tokenized.ID.Hex(), tokenized.Role)

		if err != nil {
			log.Fatal(err)
			respondWithError(w, 500, "Error while signing the token")
		} else {
			respondWithJSON(w, 200, map[string]interface{}{
				"token":       tokenString,
				"currentUser": tokenized,
			})
		}
	}
}

// AuthenticateWithUser Performs the token check and stores the current
// user in the gorilla context in order to be accessed by the handlers
func AuthenticateWithUser(handlers ...func(http.ResponseWriter, *http.Request, http.HandlerFunc)) *negroni.Negroni {
	// Default handlers
	negroniHandlers := []negroni.Handler{
		negroni.HandlerFunc(validateTokenHandler),
		negroni.HandlerFunc(getCurrentUser),
	}

	// User passed handlers
	for _, handler := range handlers {
		negroniHandlers = append(negroniHandlers, negroni.HandlerFunc(handler))
	}

	return negroni.New(negroniHandlers...)
}

// UserRecoveryMiddleware
func UserRecoveryMiddleware(handlers ...func(http.ResponseWriter, *http.Request, http.HandlerFunc)) *negroni.Negroni {
	// Default handlers
	negroniHandlers := []negroni.Handler{}

	// User passed handlers
	for _, handler := range handlers {
		negroniHandlers = append(negroniHandlers, negroni.HandlerFunc(handler))
	}

	return negroni.New(negroniHandlers...)
}

func validateTokenHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := jwtRequest.ParseFromRequest(r, jwtRequest.OAuth2Extractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte("9cdd3cbd-b48f-4187-a48a-07e67e2e9ab3"), nil //TODO: CHANGE MY SECRET
		},
	)

	context.Set(r, "jwtToken", token)

	if err == nil {
		if token.Valid {
			// var claims jwt.MapClaims
			context.Set(r, "claims", token.Claims.(jwt.MapClaims))
			next(w, r)
		} else {
			respondWithError(w, 401, "Token is not valid")
		}
	} else if token != nil { // This means we receive a jwt token but with other signature
		//Google token validation
		tokeninfoURL := fmt.Sprintf("%s?id_token=%s", constants.GoogleTokenInfoURL, token.Raw)
		response, err := http.Get(tokeninfoURL)
		if err != nil {
			log.Errorf("Error in validateTokenHanlder.googleapis -> error: %v", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else if response.StatusCode == 200 {
			context.Set(r, "claims", token.Claims.(jwt.MapClaims))
			next(w, r)
		} else {
			respondWithError(w, 401, "Token is not valid")
		}
	} else {
		// Send error by default
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access to this resource")
	}

}

func getCurrentUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	claims := context.Get(r, "claims").(jwt.MapClaims)
	userIDHex := claims["sub"].(string)
	iss := ""
	if claims["iss"] != nil {
		iss = claims["iss"].(string)
	}
	var userID bson.ObjectId

	var user *domain.User
	var authType string
	var err error

	if bson.IsObjectIdHex(userIDHex) {
		userID = bson.ObjectIdHex(userIDHex)
		user, err = domain.GetUserByID(userID)
		authType = "normal"
	} else if iss == "https://accounts.google.com" || iss == "accounts.google.com" {
		//it comes from google
		user, err = domain.GetUserByGoogleID(userIDHex)
		authType = "google"
	} else {
		respondWithError(w, 400, "Invalid token")
		return
	}

	// user, err := domain.GetUserByID(userID)

	if err != nil {
		if err.Error() == "not found" {
			respondWithError(w, 404, "User not found")
		} else {
			respondWithError(w, 400, "Invalid token")
		}
		return
	}

	context.Set(r, "currentUser", user)
	context.Set(r, "authType", authType)
	next(w, r)
}

// Me Sends to the client the current user by the
// given token in the header
func Me(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser := context.Get(r, "currentUser").(*domain.User)
	authType := context.Get(r, "authType").(string)
	jwtToken := context.Get(r, "jwtToken").(*jwt.Token)
	vars := mux.Vars(r)
	versionCode := vars["versionCode"]

	if jwtToken == nil {
		respondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	if versionCode != "" {
		vCode, e := strconv.Atoi(versionCode)
		if e != nil {
			respondWithError(w, http.StatusBadRequest, e.Error())
			return
		}
		if vCode < config.Config.GetInt("app.version") {
			respondWithError(w, http.StatusUnauthorized, "Invalid version app")
			return
		}
	}

	var tokenString string
	var err error

	switch authType {
	case "google":
		tokenString = jwtToken.Raw
	default:
		_, tokenString, err = generateJWTToken(currentUser.ID.Hex(), currentUser.Role)
	}

	if err != nil {
		respondWithError(w, 404, "Bad token")
	} else {
		respondWithJSON(w, 200, map[string]interface{}{
			"token":       tokenString,
			"currentUser": currentUser,
		})
	}
}

//UserMatchRoles checks the roles for the current user and returns a 401 response
//if no role is matched
func UserMatchRoles(roles ...string) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		currentUser := context.Get(r, "currentUser").(*domain.User)

		matchedRole := false
		for _, role := range roles {
			if currentUser.Role == role {
				matchedRole = true
				break
			}
		}

		if matchedRole {
			next(w, r)
		} else {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized access to this resource")
		}
	}
}

// PasswordRecovery sends to the client the current user by the
// given token in the header
func PasswordRecovery(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var recovery login
	parseBody(r, &recovery)

	retrievedUser, _ := domain.GetUserByEmail(recovery.Email)
	token := uuid.NewV4().String()
	domain.UpdateUserByIdWithRecoveryToken(*retrievedUser, token)

	emailBody, err := mail.ParseEmailFromTemplate("sendRecoveryPassword.html", mail.SendRecoveryPassword{
		Token: token,
	})
	if err != nil {
		log.Printf("Error in mail.ParseEmailFromTemplate: %v", err)
	}
	req := mail.CreateRequest([]string{retrievedUser.Email}, "Recuperar contraseña", emailBody)
	_, err = req.SendEmail()
	if err != nil {
		log.Printf("Error in req.SendMail: %v", err)
	}
	respondWithJSON(w, 200, map[string]interface{}{
		"token": token,
	})

}

// ChangePasswordRecovery
func ChangePasswordRecovery(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var recovery login

	parseBody(r, &recovery)

	vars := mux.Vars(r)
	tokenStr := vars["token"]

	hashPassword, err := domain.GenerateStringHashedPassword(recovery.Password)
	if err != nil {
		log.Printf("Error in GenerateStringHashedPassword: ", err)
	}

	err2 := domain.UpdatePasswordRecoveryByToken(tokenStr, hashPassword)
	if err2 != nil {
		log.Printf("Error UpdatePasswordRecoveryByToken: %v", err2)
	}

	respondWithJSON(w, 200, map[string]interface{}{})
}
