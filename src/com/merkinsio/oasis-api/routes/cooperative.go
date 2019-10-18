package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

// AppendCooperativesAPI Holds all the /cooperatives routes
func AppendCooperativesAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllCooperatives)).Methods("GET")
	router.Handle("/{cooperativeId:[A-Fa-f\\d]{24}/visitCards}", server.AuthenticateWithUser(server.GetVisitCardsByCooperativeID)).Methods("GET")
	router.Handle("/", server.AuthenticateWithUser(server.CreateCooperative)).Methods("POST")
	router.Handle("/{cooperativeId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeleteCooperative)).Methods("DELETE")
	router.Handle("/{cooperativeId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.GetCooperativeByID)).Methods("GET")
	router.Handle("/{cooperativeId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdateCooperative)).Methods("POST")
	router.Handle("/{cooperativeId:[a-f\\d]{24}}/workers", server.AuthenticateWithUser(server.GetCooperativeWorkers)).Methods("GET")
}
