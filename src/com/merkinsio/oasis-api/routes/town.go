package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendTownsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetTownsByProvince)).Queries("province", "{^[0-9]+$}").Methods("GET")
}
