package server

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type googleSignIn struct {
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
	FamilyName  string `json:"familyName"`
	GivenName   string `json:"givenName"`
	ID          string `json:"id"`
	IDToken     string `json:"idToken"`
	Name        string `json:"name"`
	Photo       string `json:"photo"`
}

// ManageGoogleAuthentication Checks the google user and creates/updates it on the DB
func ManageGoogleAuthentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body googleSignIn
	if err := parseBody(r, &body); err != nil {
		log.Errorf("Error in social.ManageGoogleAuthentication.parseBody -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the idToken token
	tokenInfoURL := fmt.Sprintf("%s?id_token=%s", constants.GoogleTokenInfoURL, body.IDToken)
	response, err := http.Get(tokenInfoURL)
	if err != nil || response.StatusCode != 200 {
		message := "Invalid token"
		if err != nil {
			message = err.Error()
		}
		respondWithError(w, http.StatusBadRequest, message)
		return
	}

	// Check if the user exists already in the database
	retrievedUser, err := domain.GetUserByEmail(body.Email)
	if err != nil {
		if err.Error() != "not found" {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else if retrievedUser == nil {
			user := domain.User{
				Email:     body.Email,
				FirstName: body.GivenName,
				LastName:  body.FamilyName,
				Role:      constants.RoleFarmer,
				JoinedBy:  constants.JoinedGoogle,
				SocialAccounts: domain.SocialAccounts{
					Google: &(domain.SocialAccount{
						ID:        body.ID,
						CreatedAt: time.Now().UTC(),
					}),
				},
			}
			if _, err := services.StoreUser(services.StoreUserPayload{User: &user}); err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
	} else {
		if retrievedUser.SocialAccounts.Google == nil || retrievedUser.SocialAccounts.Google.ID == "" {
			if err := domain.UpdateUserGoogleSocialDataByEmail(retrievedUser.Email, domain.SocialAccount{
				ID:        body.ID,
				CreatedAt: time.Now().UTC(),
			}); err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": body.IDToken,
	})

}

type facebookSignIn struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenSource string `json:"accessTokenSource"`
	ApplicationID     string `json:"applicationID"`
	ExpirationTime    int    `json:"expirationTime"`
	LastRefreshTime   int    `json:"lastRefreshTime"`
	UserID            string `json:"userID"`
}

type facebookMeResponseBody struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"`
}

// ManageFacebookAuthentication Checks the facebook user and creates/updates it on the DB
func ManageFacebookAuthentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var fbBody facebookSignIn
	if err := parseBody(r, &fbBody); err != nil {
		log.Errorf("Error in social.ManageFacebookAuthentication.parseBody -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the accessToken
	meURL := fmt.Sprintf("%s?fields=email,first_name,last_name&access_token=%s", constants.FacebookMeURL, fbBody.AccessToken)
	response, err := http.Get(meURL)
	if err != nil || response.StatusCode != 200 {
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else if response != nil {
			respondWithJSON(w, http.StatusBadRequest, response)
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid token")
		}
		return
	}

	// Retrieve faceboook user data from response
	body := facebookMeResponseBody{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		log.Errorf("Error in social.ManageFacebookAuthentication.Decode -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer response.Body.Close()

	// Check if the user exists already in the database
	retrievedUser, err := domain.GetUserByEmail(body.Email)
	userID := ""
	if err != nil {
		if err.Error() != "not found" {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		} else if retrievedUser == nil {
			user := domain.User{
				Email:     body.Email,
				FirstName: body.FirstName,
				LastName:  body.LastName,
				Role:      constants.RoleFarmer,
				JoinedBy:  constants.JoinedFacebook,
				SocialAccounts: domain.SocialAccounts{
					Facebook: &(domain.SocialAccount{
						ID:        body.ID,
						CreatedAt: time.Now().UTC(),
					}),
				},
			}
			storedUser, err := services.StoreUser(services.StoreUserPayload{User: &user})
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			userID = storedUser.ID.Hex()
		}
	} else {
		if retrievedUser.SocialAccounts.Google == nil || retrievedUser.SocialAccounts.Facebook.ID == "" {
			if err := domain.UpdateUserFacebookSocialDataByEmail(retrievedUser.Email, domain.SocialAccount{
				ID:        body.ID,
				CreatedAt: time.Now().UTC(),
			}); err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			userID = retrievedUser.ID.Hex()
		}
	}

	// Genereate the new JWTTOken
	_, generatedToken, err := generateJWTToken(userID, constants.RoleFarmer)
	if err != nil {
		log.Errorf("Error in social.ManageFacebookAuthentication.generateJWTToken -> error: %v", err.Error())
		panic(err.Error())
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": generatedToken,
	})
}
