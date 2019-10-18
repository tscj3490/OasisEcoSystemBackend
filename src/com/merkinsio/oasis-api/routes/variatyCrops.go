package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendVariatyCropsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllVariaty)).Methods("GET")
	router.Handle("/", server.AuthenticateWithUser(server.CreateVariatyCrop)).Methods("POST")
	router.Handle("/{variatyId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeleteVariatyCrop)).Methods("DELETE")
	router.Handle("/{variatyId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdateVariatyCrop)).Methods("POST")
	router.Handle("/cultivation/", server.AuthenticateWithUser(server.GetConditionalVariaties)).Queries("declaredCultivation", "{^[A-Z ]+$}").Methods("GET")
}
