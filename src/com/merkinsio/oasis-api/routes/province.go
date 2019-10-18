package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendProvincesAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllProvinces)).Methods("GET")
}
