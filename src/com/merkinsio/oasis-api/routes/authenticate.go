package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendAuthenticateRoutes(router *mux.Router) {
	router.HandleFunc("/", server.Login).Methods("POST")
	router.Handle("/google", server.Wrap(server.ManageGoogleAuthentication)).Methods("POST")
	router.Handle("/facebook", server.Wrap(server.ManageFacebookAuthentication)).Methods("POST")
}
