package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendParamsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllParams)).Methods("GET")
	router.Handle("/{treatment:[a-zA-Z_]+}", server.AuthenticateWithUser(server.GetConditionalParams)).Methods("GET")
	router.Handle("/", server.AuthenticateWithUser(server.CreateParam)).Methods("POST")
	router.Handle("/{paramId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeleteParam)).Methods("DELETE")
	router.Handle("/{paramId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdateParam)).Methods("POST")
}
