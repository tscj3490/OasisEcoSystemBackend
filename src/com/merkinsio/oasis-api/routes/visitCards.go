package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendVisitCardRoutes(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetVisitCards)).Methods("GET")
	router.Handle("/{visitCardId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.GetVisitCardByID)).Methods("GET")
	router.Handle("/", server.AuthenticateWithUser(server.CreateVisitCard)).Methods("POST")
	router.Handle("/{visitCardId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdateVisitCard)).Methods("POST")
	router.Handle("/{visitCardId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeleteVisitCard)).Methods("DELETE")
}
