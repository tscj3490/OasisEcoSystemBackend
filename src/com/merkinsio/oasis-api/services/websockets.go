package services

import (
	"com/merkinsio/oasis-api/config"
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

//TODO: REMOVE CLIENTS AND USE REDIS INSTEAD
var Clients map[*websocket.Conn]map[string]string     // connected clients
var NotificationsChannel chan domain.WebsocketMessage // broadcast channel

//TODO: Configure the upgrader
var WebsocketUpgrader websocket.Upgrader

func init() {
	Clients = make(map[*websocket.Conn]map[string]string)     // connected clients
	NotificationsChannel = make(chan domain.WebsocketMessage) // broadcast channel

	//TODO: Configure the upgrader
	WebsocketUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			logrus.Debugf("Origin request: %v", r.Header.Get("Origin"))
			if config.Config.GetString("environment") == "development" {
				return true
			}
			return r.Header.Get("Origin") == config.Config.GetString("http.origin")
		},
	}

}

func EmitMessage(w domain.WebsocketMessage) {
	// zone, err := domain.GetZoneByID(w.ZoneID)
	// if err == nil {
	// w.UserIDs = append(w.UserIDs, zone.AdvisorID)
	// w.UserIDs = append(w.UserIDs, zone.Sunboys...)

	NotificationsChannel <- w
	// } else {
	// logrus.Warnf("Message was not emitted; %v", err.Error())
	// }
}

//KeepAlive Keeps the socket connection alive
func KeepAlive(c *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now().UTC()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now().UTC()
		return nil
	})

	go func() {
		for {
			keepAliveMessage := domain.WebsocketMessage{
				Type: constants.WebsocketKeepAlive,
			}
			err := c.WriteJSON(keepAliveMessage)
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().UTC().Sub(lastResponse) > timeout {
				delete(Clients, c)
				c.Close()
				return
			}
		}
	}()
}
