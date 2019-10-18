package server

import (
	"com/merkinsio/oasis-api/domain"
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2/bson"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func parseBody(r *http.Request, payload interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(payload); err != nil {
		return err
	}
	return nil
}

func generateJWTToken(id, role string) (*jwt.Token, string, error) {
	authTokenExpiration := time.Now().UTC().Add(domain.AuthTokenValidTime).Unix()
	authClaims := domain.TokenClaims{
		jwt.StandardClaims{
			Subject:   id,
			ExpiresAt: authTokenExpiration,
		},
		role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	tokenString, err := token.SignedString([]byte("9cdd3cbd-b48f-4187-a48a-07e67e2e9ab3")) //TODO: CHANGE MY SECRET

	return token, tokenString, err
}

func parseStringToObjectID(str string) *bson.ObjectId {
	if bson.IsObjectIdHex(str) {
		result := bson.ObjectIdHex(str)
		return &result
	}
	return nil
}

//Wrap Wraps the handlers into a negroni handler to be used in the routes
func Wrap(handlers ...func(http.ResponseWriter, *http.Request, http.HandlerFunc)) *negroni.Negroni {
	negroniHandlers := []negroni.Handler{}

	// User passed handlers
	for _, handler := range handlers {
		negroniHandlers = append(negroniHandlers, negroni.HandlerFunc(handler))
	}
	return negroni.New(negroniHandlers...)
}
