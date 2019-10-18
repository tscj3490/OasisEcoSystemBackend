package server

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser := context.Get(r, "currentUser").(*domain.User)
	jwtToken := context.Get(r, "jwtToken").(*jwt.Token)

	// Upgrade initial GET request to a websocket
	ws, err := services.WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	services.Clients[ws] = map[string]string{
		"token":  jwtToken.Raw,
		"userId": currentUser.ID.Hex(),
	}

	for {
		var msg domain.WebsocketMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(services.Clients, ws)
			break
		}

		if msg.Type == constants.WebsocketGreetings {
			ws.WriteJSON(msg)
		}
		// Send the newly received message to the broadcast channel
		services.NotificationsChannel <- msg
	}
}

func HandleSocketMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-services.NotificationsChannel
		// Send it out to every client that is currently connected
		for client, values := range services.Clients {
			for _, id := range msg.UserIDs {
				if values["userId"] == id.Hex() {
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(services.Clients, client)
					}
				}
			}

		}
	}
}
